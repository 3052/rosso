package peacock

import (
   "bytes"
   "crypto/md5"
   "crypto/tls"
   _ "embed"
   "encoding/hex"
   "encoding/json"
   "fmt"
   "io"
   "log"
   "net/http"
   "net/url"
   "strings"
   "time"
)

// oauth_authorize.go
const idBase = "https://rango.id.peacocktv.com"

// signin.go
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

// OAuthAuthorizeResponse is the Siren JSON response from GET /oauth/authorize/service/international.
type OAuthAuthorizeResponse struct {
   Class      []string `json:"class"`
   Properties struct {
      EventType   string `json:"eventType"`
      Href        string `json:"href"`
      Type        string `json:"type"`
      AccessToken string `json:"access_token"`
   } `json:"properties"`
}

// OAuthAuthorize follows the OAuth authorize redirect to obtain the access token.
// It requires the skyCEsidmesso01 cookie set by SignIn.
func OAuthAuthorize(cookie *http.Cookie) (*OAuthAuthorizeResponse, error) {
   if cookie == nil {
      return nil, fmt.Errorf("oauth authorize: nil cookie")
   }
   if cookie.Value == "" {
      return nil, fmt.Errorf("oauth authorize: empty cookie value")
   }
   params := url.Values{}
   params.Set("client_id", "nbcu_tvclient")
   params.Set("response_type", "token")
   oauthUrl := idBase + "/oauth/authorize/service/international?" + params.Encode()
   req, err := http.NewRequest(http.MethodGet, oauthUrl, nil)
   if err != nil {
      return nil, fmt.Errorf("oauth authorize: create request: %w", err)
   }
   req.AddCookie(cookie)
   req.Header.Set("x-skyott-territory", "US")
   req.Header.Set("x-skyott-provider", "NBCU")
   resp, err := doRequest(http.DefaultClient, req)
   if err != nil {
      return nil, fmt.Errorf("oauth authorize: %w", err)
   }
   defer resp.Body.Close()

   if resp.StatusCode != http.StatusOK {
      body, _ := io.ReadAll(resp.Body)
      return nil, fmt.Errorf("oauth authorize: bad status %d: %s", resp.StatusCode, body)
   }

   var out OAuthAuthorizeResponse
   if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
      return nil, fmt.Errorf("oauth authorize: decode: %w", err)
   }

   if out.Properties.AccessToken == "" {
      return nil, fmt.Errorf("oauth authorize: access_token not found in response")
   }

   return &out, nil
}

func (*OAuthAuthorizeResponse) CachePath() string {
   return "peacock/oauth_authorize"
}

// SignInParams holds the credentials for the sign-in request.
type SignInParams struct {
   UserIdentifier string
   Password       string
   RememberMe     bool
}

// SignInResponse is the Siren JSON response from POST /signin/service/international.
// Cookie is the skyCEsidmesso01 cookie from the Set-Cookie header; it is required
// as input to OAuthAuthorize.
type SignInResponse struct {
   Cookie     *http.Cookie `json:"-"`
   Class      []string     `json:"class"`
   Properties struct {
      TargetApi string `json:"targetApi"`
      Type      string `json:"type"`
      EventType string `json:"eventType"`
      Href      string `json:"href"`
      Data      struct {
         DeviceID string `json:"deviceid"`
      } `json:"data"`
   } `json:"properties"`
}

// SignIn authenticates with email and password via POST /signin/service/international.
// The skyCEsidmesso01 cookie from the response is returned in SignInResponse.Cookie
// for use by OAuthAuthorize. The device id is available in Properties.Data.DeviceID for
// use by ExchangeToken.
func SignIn(params *SignInParams) (*SignInResponse, error) {
   if params.UserIdentifier == "" {
      return nil, fmt.Errorf("sign in: empty userIdentifier")
   }
   if params.Password == "" {
      return nil, fmt.Errorf("sign in: empty password")
   }
   signInUrl := idBase + "/signin/service/international"
   form := url.Values{}
   form.Set("password", params.Password)
   form.Set("userIdentifier", params.UserIdentifier)
   req, err := http.NewRequest(http.MethodPost, signInUrl, strings.NewReader(form.Encode()))
   if err != nil {
      return nil, fmt.Errorf("sign in: create request: %w", err)
   }
   req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
   req.Header.Set("x-skyott-territory", "US")
   req.Header.Set("x-skyott-provider", "NBCU")
   req.Header.Set("user-agent", "user-agent (Android; Build/user-agent)")
   resp, err := doRequest(http.DefaultClient, req)
   if err != nil {
      return nil, fmt.Errorf("sign in: %w", err)
   }
   defer resp.Body.Close()

   if resp.StatusCode != http.StatusCreated {
      body, _ := io.ReadAll(resp.Body)
      return nil, fmt.Errorf("sign in: bad status %d: %s", resp.StatusCode, body)
   }

   var out SignInResponse
   for _, cookie := range resp.Cookies() {
      if cookie.Name == "skyCEsidmesso01" {
         out.Cookie = cookie
         break
      }
   }
   if out.Cookie == nil {
      return nil, fmt.Errorf("sign in: skyCEsidmesso01 cookie not found in response")
   }

   if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
      return nil, fmt.Errorf("sign in: decode: %w", err)
   }
   return &out, nil
}

func (*SignInResponse) CachePath() string {
   return "rosso/peacock/SignInResponse"
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
