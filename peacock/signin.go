package peacock

import (
   "encoding/json"
   "fmt"
   "io"
   "net/http"
   "net/url"
   "strings"
)

const idBase = "https://rango.id.peacocktv.com"

// SignInParams holds the credentials for the sign-in request.
type SignInParams struct {
   UserIdentifier string
   Password       string
   RememberMe     bool
}

// SignInResponse is the Siren JSON response from POST /signin/service/international.
type SignInResponse struct {
   Class      []string `json:"class"`
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
// The skyCEsidmesso01 cookie from the response is stored on the Client for use by OAuthAuthorize.
func (c *Client) SignIn(params *SignInParams) (*SignInResponse, error) {
   if params.UserIdentifier == "" {
      return nil, fmt.Errorf("sign in: empty userIdentifier")
   }
   if params.Password == "" {
      return nil, fmt.Errorf("sign in: empty password")
   }

   continuationParams := url.Values{}
   continuationParams.Set("response_type", "token")
   continuationParams.Set("client_id", "nbcu_tvclient")
   continuationParams.Set("redirect_uri", "nbcu://auth")
   continuationParams.Set("api_id", "oauth")
   continuationUrl := idBase + "/oauth/authorize/service/international?" + continuationParams.Encode()

   signInParams := url.Values{}
   signInParams.Set("continuationUrl", continuationUrl)
   signInUrl := idBase + "/signin/service/international?" + signInParams.Encode()

   form := url.Values{}
   form.Set("password", params.Password)
   form.Set("rememberMe", fmt.Sprintf("%t", params.RememberMe))
   form.Set("userIdentifier", params.UserIdentifier)

   req, err := http.NewRequest(http.MethodPost, signInUrl, strings.NewReader(form.Encode()))
   if err != nil {
      return nil, fmt.Errorf("sign in: create request: %w", err)
   }
   req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
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

   resp, err := doRequest(c.HTTP, req)
   if err != nil {
      return nil, fmt.Errorf("sign in: %w", err)
   }
   defer resp.Body.Close()

   if resp.StatusCode != http.StatusCreated {
      body, _ := io.ReadAll(resp.Body)
      return nil, fmt.Errorf("sign in: bad status %d: %s", resp.StatusCode, body)
   }

   for _, cookie := range resp.Cookies() {
      if cookie.Name == "skyCEsidmesso01" {
         c.skyCEsidmesso01 = cookie.Value
         break
      }
   }
   if c.skyCEsidmesso01 == "" {
      return nil, fmt.Errorf("sign in: skyCEsidmesso01 cookie not found in response")
   }

   var out SignInResponse
   if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
      return nil, fmt.Errorf("sign in: decode: %w", err)
   }
   return &out, nil
}

// signin.go
