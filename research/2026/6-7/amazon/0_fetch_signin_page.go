package amazon

import (
   "fmt"
   "io"
   "net/http"
   "net/url"
   "regexp"
)

const AuthDeviceType = "A1MPSLFC7L5AFK"

func FetchSignInPage(client *http.Client, deviceID, frcCookie string) (url.Values, string, string, error) {
   codeVerifier, codeChallenge, err := GeneratePKCE()
   if err != nil {
      return nil, "", "", err
   }

   mapMDCookie, err := GenerateMapMD()
   if err != nil {
      return nil, "", "", err
   }

   // Inject the map-md and static frc cookies into the jar
   amazonURL, _ := url.Parse("https://www.amazon.com")
   cookiesToSet := []*http.Cookie{
      {Name: "map-md", Value: mapMDCookie, Domain: ".amazon.com", Path: "/"},
   }
   if frcCookie != "" {
      cookiesToSet = append(cookiesToSet, &http.Cookie{Name: "frc", Value: frcCookie, Domain: ".amazon.com", Path: "/"})
   }
   client.Jar.SetCookies(amazonURL, cookiesToSet)

   reqURL, err := url.Parse("https://www.amazon.com/ap/signin")
   if err != nil {
      return nil, "", "", err
   }

   clientID := GenerateClientID(deviceID, AuthDeviceType)

   q := reqURL.Query()
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
   q.Set("openid.oa2.client_id", clientID)
   q.Set("disableLoginPrepopulate", "0")
   q.Set("openid.ns", "http://specs.openid.net/auth/2.0")

   reqURL.RawQuery = q.Encode()
   signInUrlStr := reqURL.String()

   req, err := http.NewRequest(http.MethodGet, signInUrlStr, nil)
   if err != nil {
      return nil, "", "", err
   }

   req.Header.Set("upgrade-insecure-requests", "1")
   req.Header.Set("user-agent", "Mozilla/5.0 (Linux; Android 11; sdk_gphone_x86_64 Build/RSR1.240422.006; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/83.0.4103.106 Mobile Safari/537.36")
   req.Header.Set("accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
   req.Header.Set("x-requested-with", "com.amazon.avod.thirdpartyclient")
   req.Header.Set("sec-fetch-site", "none")
   req.Header.Set("sec-fetch-mode", "navigate")
   req.Header.Set("sec-fetch-user", "?1")
   req.Header.Set("sec-fetch-dest", "document")
   req.Header.Set("accept-language", "en-US,en;q=0.9")

   resp, err := client.Do(req)
   if err != nil {
      return nil, "", "", err
   }
   defer resp.Body.Close()

   if resp.StatusCode != http.StatusOK {
      return nil, "", "", fmt.Errorf("expected 200 OK, got status code: %d", resp.StatusCode)
   }

   bodyBytes, err := io.ReadAll(resp.Body)
   if err != nil {
      return nil, "", "", err
   }
   html := string(bodyBytes)

   formRegex := regexp.MustCompile(`(?s)<form[^>]*name="signIn"[^>]*method="post"[^>]*action="[^"]*signin[^"]*"[^>]*>(.*?)</form>`)
   formMatch := formRegex.FindStringSubmatch(html)
   if len(formMatch) < 2 {
      return nil, "", "", fmt.Errorf("signIn form not found in the HTML response")
   }
   formHtml := formMatch[1]

   inputRegex := regexp.MustCompile(`(?i)<input[^>]+>`)
   nameRegex := regexp.MustCompile(`(?i)name=['"]([^'"]+)['"]`)
   valueRegex := regexp.MustCompile(`(?i)value=['"]([^'"]*)['"]`)

   formValues := url.Values{}
   inputs := inputRegex.FindAllString(formHtml, -1)

   for _, inputStr := range inputs {
      nameMatch := nameRegex.FindStringSubmatch(inputStr)
      if len(nameMatch) >= 2 {
         name := nameMatch[1]
         value := ""
         valueMatch := valueRegex.FindStringSubmatch(inputStr)
         if len(valueMatch) >= 2 {
            value = valueMatch[1]
         }
         formValues.Set(name, value)
      }
   }

   if formValues.Get("anti-csrftoken-a2z") == "" {
      return nil, "", "", fmt.Errorf("missing 'anti-csrftoken-a2z' in form")
   }

   return formValues, codeVerifier, signInUrlStr, nil
}
