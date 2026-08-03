package peacock

import (
   "bytes"
   "crypto/md5"
   "crypto/tls"
   _ "embed"
   "encoding/hex"
   "encoding/json"
   "fmt"
   "log"
   "net/http"
   "time"
)

// client.go
const playBase = "https://play.clients.peacocktv.com"

//go:embed cert.pem
var certPEM []byte

//go:embed key.pem
var keyPEM []byte

// doRequest logs the request method and URL, then sends the request
// using the provided http.Client.
func doRequest(client *http.Client, req *http.Request) (*http.Response, error) {
   log.Println(req.Method, req.URL)
   return client.Do(req)
}

// mtlsClient returns an *http.Client configured with the embedded mTLS
// certificate, ProxyFromEnvironment, and the given timeout.
func mtlsClient() (*http.Client, error) {
   cert, err := tls.X509KeyPair(certPEM, keyPEM)
   if err != nil {
      return nil, err
   }

   return &http.Client{
      Transport: &http.Transport{
         Proxy: http.ProxyFromEnvironment,
         TLSClientConfig: &tls.Config{
            Certificates: []tls.Certificate{cert},
         },
      },
   }, nil
}

// TokenResponse is the response from POST /auth/throttled/tokens.
// On a non-2xx response the server populates ErrorCode/Description
// instead of the token fields, in which case TokenResponse doubles
// as the returned error.
type TokenResponse struct {
   UserToken                     string    `json:"userToken"`
   TokenExpiryTime               time.Time `json:"tokenExpiryTime"`
   RecommendedTokenReacquireTime time.Time `json:"recommendedTokenReacquireTime"`
   ErrorCode                     string    `json:"errorCode"`
   Description                   string    `json:"description"`
}

// ExchangeToken trades the OAuth2 activation token for a long-lived user token
// using the embedded mTLS certificate. The activation token is sent
// in the request body as authToken, not as a bearer header.
// The returned user token can be used as a bearer credential for playback.
func ExchangeToken(authToken *OAuthAuthorizeResponse, signIn *SignInResponse) (*TokenResponse, error) {
   if authToken == nil {
      return nil, fmt.Errorf("exchange token: nil authToken")
   }
   if authToken.Properties.AccessToken == "" {
      return nil, fmt.Errorf("exchange token: empty access token")
   }
   if signIn == nil {
      return nil, fmt.Errorf("exchange token: nil signIn")
   }
   if signIn.Properties.Data.DeviceID == "" {
      return nil, fmt.Errorf("exchange token: empty deviceID")
   }
   client, err := mtlsClient()
   if err != nil {
      return nil, fmt.Errorf("exchange token: %w", err)
   }
   body := tokenRequest{
      Auth: tokenAuth{
         AuthScheme:        "OAUTH",
         Provider:          "NBCU",
         ProviderTerritory: "US",
         Proposition:       "NBCUOTT",
         AuthToken:         authToken.Properties.AccessToken,
      },
      Device: tokenDevice{
         Type:        "TV",
         Platform:    "ANDROIDTV",
         ID:          signIn.Properties.Data.DeviceID,
         DrmDeviceID: "UNKNOWN",
      },
   }
   raw, err := json.Marshal(body)
   if err != nil {
      return nil, fmt.Errorf("exchange token: marshal: %w", err)
   }
   hash := md5.Sum(raw)
   contentMD5 := hex.EncodeToString(hash[:])
   req, err := http.NewRequest(http.MethodPost, playBase+"/auth/throttled/tokens", bytes.NewReader(raw))
   if err != nil {
      return nil, fmt.Errorf("exchange token: create request: %w", err)
   }
   req.Header.Set("Content-Type", "application/vnd.tokens.v1+json")
   req.Header.Set("Content-MD5", contentMD5)
   resp, err := doRequest(client, req)
   if err != nil {
      return nil, fmt.Errorf("exchange token: %w", err)
   }
   defer resp.Body.Close()

   if resp.StatusCode == http.StatusUnsupportedMediaType {
      return nil, fmt.Errorf("exchange token: %s", resp.Status)
   }

   var out TokenResponse
   if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
      return nil, fmt.Errorf("exchange token: decode: %w", err)
   }

   if out.ErrorCode != "" {
      return nil, &out
   }

   return &out, nil
}

func (*TokenResponse) CachePath() string {
   return "rosso/peacock/TokenResponse"
}

// Error implements the error interface. It returns a non-empty
// string only when the server reported an errorCode, allowing
// *TokenResponse to be returned directly as the error value.
func (t *TokenResponse) Error() string {
   if t == nil || t.ErrorCode == "" {
      return ""
   }
   if t.Description == "" {
      return fmt.Sprintf("peacock: %s", t.ErrorCode)
   }
   return fmt.Sprintf("peacock: %s: %s", t.ErrorCode, t.Description)
}

type tokenAuth struct {
   AuthScheme        string `json:"authScheme"`
   AuthIssuer        string `json:"authIssuer"`
   Provider          string `json:"provider"`
   ProviderTerritory string `json:"providerTerritory"`
   Proposition       string `json:"proposition"`
   AuthToken         string `json:"authToken"`
}

type tokenDevice struct {
   Type        string `json:"type"`
   Platform    string `json:"platform"`
   ID          string `json:"id"`
   DrmDeviceID string `json:"drmDeviceId"`
}

// tokenRequest is the body sent to POST /auth/throttled/tokens.
type tokenRequest struct {
   Auth   tokenAuth   `json:"auth"`
   Device tokenDevice `json:"device"`
}

// token.go
