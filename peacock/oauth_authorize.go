package peacock

import (
   "encoding/json"
   "fmt"
   "io"
   "net/http"
   "net/url"
)

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
// It requires the session cookies set by SignIn, so the Client's HTTP transport
// must have a cookie jar configured.
func (c *Client) OAuthAuthorize() (string, error) {
   params := url.Values{}
   params.Set("client_id", "nbcu_tvclient")
   params.Set("redirect_uri", "nbcu://auth")
   params.Set("response_type", "token")

   oauthUrl := idBase + "/oauth/authorize/service/international?" + params.Encode()

   req, err := http.NewRequest(http.MethodGet, oauthUrl, nil)
   if err != nil {
      return "", fmt.Errorf("oauth authorize: create request: %w", err)
   }
   req.Header.Set("Accept", "*/*")
   req.Header.Set("Accept-Language", "en-US,en;q=0.9")
   req.Header.Set("Origin", "https://tv.clients.peacocktv.com")
   req.Header.Set("Referer", "https://tv.clients.peacocktv.com/")
   req.Header.Set("X-Requested-With", "com.peacocktv.peacockandroid")
   req.Header.Set("Sec-Fetch-Site", "same-site")
   req.Header.Set("Sec-Fetch-Mode", "cors")
   req.Header.Set("Sec-Fetch-Dest", "empty")
   req.Header.Set("User-Agent", "Mozilla/5.0 (Linux; Android 12; sdk_gphone64_x86_64 Build/SE1A.220826.008; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/91.0.4472.114 Mobile Safari/537.36")
   req.Header.Set("x-skyott-platform", "ANDROIDTV")
   req.Header.Set("x-skyott-proposition", "NBCUOTT")
   req.Header.Set("x-skyott-provider", "NBCU")
   req.Header.Set("x-skyott-territory", "US")
   req.Header.Set("x-skyott-activeterritory", "US")
   req.Header.Set("x-skyott-language", "en-US")
   req.Header.Set("x-skyott-device", "TV")
   req.Header.Set("x-skyott-broadcastregions", "INPATTERN_US_CENTRAL")
   req.Header.Set("x-deviceid", c.DeviceID)
   req.Header.Set("x-skyint-requestid", randomUUID())

   resp, err := c.HTTP.Do(req)
   if err != nil {
      return "", fmt.Errorf("oauth authorize: %w", err)
   }
   defer resp.Body.Close()

   if resp.StatusCode != http.StatusOK {
      body, _ := io.ReadAll(resp.Body)
      return "", fmt.Errorf("oauth authorize: bad status %d: %s", resp.StatusCode, body)
   }

   var out OAuthAuthorizeResponse
   if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
      return "", fmt.Errorf("oauth authorize: decode: %w", err)
   }

   if out.Properties.AccessToken == "" {
      return "", fmt.Errorf("oauth authorize: access_token not found in response")
   }

   return out.Properties.AccessToken, nil
}

// oauth_authorize.go
