package amazon

import (
   "fmt"
   "io"
   "net/http"
   "net/url"
)

func InitSignin(client *http.Client, clientId, codeChallenge string) (string, map[string]string, error) {
   baseURL := "https://www.amazon.com/ap/signin"
   u, _ := url.Parse(baseURL)
   q := u.Query()
   q.Set("openid.pape.max_auth_age", "0")
   q.Set("openid.identity", "http://specs.openid.net/auth/2.0/identifier_select")
   q.Set("accountStatusPolicy", "P1")
   q.Set("language", "en_US")
   q.Set("pageId", "amzn_dv_ios_blue")
   q.Set("openid.return_to", "https://www.amazon.com/ap/maplanding")
   q.Set("openid.assoc_handle", "amzn_piv_android_v2_us")
   q.Set("openid.oa2.response_type", "code")
   q.Set("countryCode", "US")
   q.Set("openid.mode", "checkid_setup")
   q.Set("openid.ns.pape", "http://specs.openid.net/extensions/pape/1.0")
   q.Set("openid.oa2.code_challenge_method", "S256")
   q.Set("openid.ns.oa2", "http://www.amazon.com/ap/ext/oauth/2")
   q.Set("openid.oa2.code_challenge", codeChallenge)
   q.Set("openid.oa2.scope", "device_auth_access")
   q.Set("openid.claimed_id", "http://specs.openid.net/auth/2.0/identifier_select")
   q.Set("openid.oa2.client_id", clientId)
   q.Set("disableLoginPrepopulate", "0")
   q.Set("openid.ns", "http://specs.openid.net/auth/2.0")
   u.RawQuery = q.Encode()

   req, err := http.NewRequest("GET", u.String(), nil)
   if err != nil {
      return "", nil, fmt.Errorf("failed to create request: %w", err)
   }
   req.Header.Set("User-Agent", "Mozilla/5.0 (Linux; Android 11; sdk_gphone_x86_64 Build/RSR1.240422.006; wv) AppleWebKit/537.36")
   req.Header.Set("x-requested-with", "com.amazon.avod.thirdpartyclient")

   resp, err := client.Do(req)
   if err != nil {
      return "", nil, fmt.Errorf("request failed: %w", err)
   }
   defer resp.Body.Close()

   body, err := io.ReadAll(resp.Body)
   if err != nil {
      return "", nil, fmt.Errorf("failed to read body: %w", err)
   }

   actionUrl, hiddenParams := extractFormActionAndHiddenInputs(string(body), baseURL)
   return actionUrl, hiddenParams, nil
}
