package peacock

import (
   "bytes"
   "crypto/md5"
   "encoding/hex"
   "encoding/json"
   "fmt"
   "net/http"
   "time"
)

const playBase = "https://play.clients.peacocktv.com"

// TokenResponse is the response from POST /auth/throttled/tokens.
type TokenResponse struct {
   UserToken                     string    `json:"userToken"`
   TokenExpiryTime               time.Time `json:"tokenExpiryTime"`
   RecommendedTokenReacquireTime time.Time `json:"recommendedTokenReacquireTime"`
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

// ExchangeToken trades the OAuth2 activation token for a long-lived user token
// using the embedded mTLS certificate. The activation token is sent
// in the request body as authToken, not as a bearer header.
// The returned user token can be used as a bearer credential for playback.
func (c *Client) ExchangeToken(authToken *OAuthAuthorizeResponse) (*TokenResponse, error) {
   if authToken == nil {
      return nil, fmt.Errorf("exchange token: nil authToken")
   }
   if authToken.Properties.AccessToken == "" {
      return nil, fmt.Errorf("exchange token: empty access token")
   }

   client, err := mtlsClient(c.HTTP.Timeout)
   if err != nil {
      return nil, fmt.Errorf("exchange token: %w", err)
   }

   body := tokenRequest{
      Auth: tokenAuth{
         AuthScheme:        "OAUTH",
         AuthIssuer:        "NOWTV",
         Provider:          "NBCU",
         ProviderTerritory: "US",
         Proposition:       "NBCUOTT",
         AuthToken:         authToken.Properties.AccessToken,
      },
      Device: tokenDevice{
         Type:        "TV",
         Platform:    "ANDROIDTV",
         ID:          c.DeviceID,
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

   req.Header.Set("User-Agent", "Mozilla/5.0 (Linux; Android 12; sdk_gphone64_x86_64 Build/SE1A.220826.008; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/91.0.4472.114 Mobile Safari/537.36")
   req.Header.Set("x-skyott-platform", "ANDROIDTV")
   req.Header.Set("x-skyott-proposition", "NBCUOTT")
   req.Header.Set("x-skyott-provider", "NBCU")
   req.Header.Set("x-skyott-territory", "US")
   req.Header.Set("x-skyott-activeterritory", "US")
   req.Header.Set("x-skyott-language", "en-US")
   req.Header.Set("x-skyott-device", "TV")
   req.Header.Set("x-skyott-broadcastregions", "INPATTERN_US_CENTRAL")
   req.Header.Set("x-skyint-requestid", randomUUID())
   req.Header.Set("Content-Type", "application/vnd.tokens.v1+json")
   req.Header.Set("Accept", "application/vnd.tokens.v1+json")
   req.Header.Set("Content-MD5", contentMD5)
   req.Header.Set("Origin", "https://tv.clients.peacocktv.com")

   resp, err := doRequest(client, req)
   if err != nil {
      return nil, fmt.Errorf("exchange token: %w", err)
   }
   defer resp.Body.Close()

   if resp.StatusCode != http.StatusOK {
      return nil, fmt.Errorf("exchange token: bad status %d", resp.StatusCode)
   }

   var out TokenResponse
   if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
      return nil, fmt.Errorf("exchange token: decode: %w", err)
   }
   return &out, nil
}

// token.go
