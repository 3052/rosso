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

// oauth_authorize.go
