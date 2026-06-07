package amazon

import (
   "fmt"
   "io"
   "net/http"
   "net/url"
   "regexp"
)

// FetchCVFPage requests the CVF OTP page (using the URL returned by SubmitCredentials).
func FetchCVFPage(client *http.Client, cvfUrl string, referer string) (url.Values, error) {
   req, err := http.NewRequest(http.MethodGet, cvfUrl, nil)
   if err != nil {
      return nil, err
   }

   req.Header.Set("upgrade-insecure-requests", "1")
   req.Header.Set("user-agent", "Mozilla/5.0 (Linux; Android 11; sdk_gphone_x86_64 Build/RSR1.240422.006; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/83.0.4103.106 Mobile Safari/537.36")
   req.Header.Set("accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
   req.Header.Set("x-requested-with", "com.amazon.avod.thirdpartyclient")
   req.Header.Set("referer", referer) // Required by WAF
   req.Header.Set("sec-fetch-site", "same-origin")
   req.Header.Set("sec-fetch-mode", "navigate")
   req.Header.Set("sec-fetch-user", "?1")
   req.Header.Set("sec-fetch-dest", "document")
   req.Header.Set("accept-language", "en-US,en;q=0.9")

   // Ensure we do not follow any further redirects silently
   client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
      return http.ErrUseLastResponse
   }

   resp, err := client.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()

   if resp.StatusCode != http.StatusOK {
      return nil, fmt.Errorf("expected 200 OK, got status code: %d", resp.StatusCode)
   }

   bodyBytes, err := io.ReadAll(resp.Body)
   if err != nil {
      return nil, err
   }
   html := string(bodyBytes)

   captchaRegex := regexp.MustCompile(`(?s)<form[^>]*id="cvf-aamation-challenge-form"[^>]*>(.*?)</form>`)
   if captchaMatch := captchaRegex.FindStringSubmatch(html); len(captchaMatch) >= 2 {
      return nil, fmt.Errorf("CAPTCHA_REQUIRED")
   }

   formRegex := regexp.MustCompile(`(?s)<form[^>]*id="verification-code-form"[^>]*>(.*?)</form>`)
   formMatch := formRegex.FindStringSubmatch(html)
   if len(formMatch) < 2 {
      return nil, fmt.Errorf("verification-code-form not found in the HTML response")
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
      return nil, fmt.Errorf("missing 'anti-csrftoken-a2z' in CVF form")
   }

   return formValues, nil
}
