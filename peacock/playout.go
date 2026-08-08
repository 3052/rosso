package peacock

import (
   "bytes"
   "crypto/md5"
   "encoding/hex"
   "encoding/json"
   "fmt"
   "net/http"
   "net/url"
   "time"
)

// PlayoutVod requests a VOD playout URL from POST /video/playouts/vod using the
// embedded mTLS certificate. token supplies the userToken; providerVariantID
// identifies the asset to play. vcodec selects the requested video codec
// (e.g. "H264" or "H265"). protection selects the DRM system ("WIDEVINE" or
// "PLAYREADY").
func PlayoutVod(token *TokenResponse, providerVariantID, vcodec, protection string) (*PlayoutVodResponse, error) {
   if token == nil {
      return nil, fmt.Errorf("playout vod: nil token")
   }
   if token.UserToken == "" {
      return nil, fmt.Errorf("playout vod: empty userToken")
   }
   if providerVariantID == "" {
      return nil, fmt.Errorf("playout vod: empty providerVariantID")
   }
   if vcodec == "" {
      return nil, fmt.Errorf("playout vod: empty vcodec")
   }
   if protection == "" {
      return nil, fmt.Errorf("playout vod: empty protection")
   }
   body := playoutRequest{
      Device: playoutDevice{
         Capabilities: []playoutCapability{
            {
               Acodec:     "AAC",
               Container:  "ISOBMFF",
               Protection: protection,
               Transport:  "DASH",
               Vcodec:     vcodec,
            },
         },
         MaxVideoFormat: "UHD",
      },
      ProviderVariantID:            providerVariantID,
      PersonaParentalControlRating: "9",
   }
   raw, err := json.Marshal(body)
   if err != nil {
      return nil, fmt.Errorf("playout vod: marshal: %w", err)
   }
   hash := md5.Sum(raw)
   contentMD5 := hex.EncodeToString(hash[:])
   req, err := http.NewRequest(http.MethodPost, playBase+"/video/playouts/vod", bytes.NewReader(raw))
   if err != nil {
      return nil, fmt.Errorf("playout vod: create request: %w", err)
   }
   req.Header.Set("x-skyott-usertoken", token.UserToken)
   req.Header.Set("Content-Type", "application/vnd.playvod.v1+json")
   req.Header.Set("Content-MD5", contentMD5)
   client, err := mtlsClient()
   if err != nil {
      return nil, fmt.Errorf("playout vod: %w", err)
   }
   resp, err := doRequest(client, req)
   if err != nil {
      return nil, fmt.Errorf("playout vod: %w", err)
   }
   defer resp.Body.Close()

   var out PlayoutVodResponse
   if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
      return nil, fmt.Errorf("playout vod: decode: %w", err)
   }
   if out.ErrorCode != "" {
      return nil, &out
   }
   return &out, nil
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

type playoutCapability struct {
   Acodec     string `json:"acodec"`
   Container  string `json:"container"`
   Protection string `json:"protection"`
   Transport  string `json:"transport"`
   Vcodec     string `json:"vcodec"`
}

type playoutClient struct {
   ThirdParties   []string `json:"thirdParties"`
   VariantCapable bool     `json:"variantCapable"`
}

type playoutDevice struct {
   Capabilities          []playoutCapability `json:"capabilities"`
   MaxVideoFormat        string              `json:"maxVideoFormat"`
   SupportedColourSpaces []string            `json:"supportedColourSpaces"`
   Model                 string              `json:"model"`
   HdcpEnabled           bool                `json:"hdcpEnabled"`
}

type playoutRequest struct {
   Device                       playoutDevice `json:"device"`
   Client                       playoutClient `json:"client"`
   ProviderVariantID            string        `json:"providerVariantId"`
   PersonaParentalControlRating string        `json:"personaParentalControlRating"`
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

func (*PlayoutVodResponse) CachePath() string {
   return "rosso/peacock/PlayoutVodResponse"
}

// Error implements the error interface.
func (r *PlayoutVodResponse) Error() string {
   return r.ErrorCode + ": " + r.Description
}

// Fastly returns the parsed URL of the FASTLY CDN endpoint from the playout response.
func (r *PlayoutVodResponse) Fastly() (*url.URL, error) {
   for _, endpoint := range r.Asset.Endpoints {
      if endpoint.Cdn == "FASTLY" {
         parsed, err := url.Parse(endpoint.Url)
         if err != nil {
            return nil, fmt.Errorf("fastly: parse url: %w", err)
         }
         return parsed, nil
      }
   }
   return nil, fmt.Errorf("fastly cdn endpoint not found")
}

// playout.go
