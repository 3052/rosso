package amazon

import (
   "fmt"
   "io"
   "net/http"
   "net/url"
   "regexp"
)

// FetchCVFPage requests the CVF OTP page and parses the HTML to extract
// all the hidden input fields required for the verification POST request.
func FetchCVFPage(cvfUrl string, cookies []*http.Cookie) (url.Values, []*http.Cookie, error) {
   req, err := http.NewRequest("GET", cvfUrl, nil)
   if err != nil {
      return nil, nil, err
   }

   req.Header.Set("User-Agent", "Mozilla/5.0 (Linux; Android 11; sdk_gphone_x86_64 Build/RSR1.240422.006; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/83.0.4103.106 Mobile Safari/537.36")
   req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
   req.Header.Set("X-Requested-With", "com.amazon.avod.thirdpartyclient")

   for _, cookie := range cookies {
      req.AddCookie(cookie)
   }

   client := &http.Client{}
   resp, err := client.Do(req)
   if err != nil {
      return nil, nil, err
   }
   defer resp.Body.Close()

   if resp.StatusCode != http.StatusOK {
      return nil, nil, fmt.Errorf("expected 200 OK, got status code: %d", resp.StatusCode)
   }

   bodyBytes, err := io.ReadAll(resp.Body)
   if err != nil {
      return nil, nil, err
   }
   html := string(bodyBytes)

   // Isolate the specific form block to avoid grabbing inputs from the "Resend" or "WhatsApp" forms
   formRegex := regexp.MustCompile(`(?s)<form[^>]*id="verification-code-form"[^>]*>(.*?)</form>`)
   formMatch := formRegex.FindStringSubmatch(html)
   if len(formMatch) < 2 {
      return nil, nil, fmt.Errorf("verification-code-form not found in the HTML response")
   }
   formHtml := formMatch[1]

   // Extract all inputs within that form
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

   // NOTE: Before passing formValues to VerifyOTP(), you will need to add the actual OTP code:
   // formValues.Set("code", "123456")

   return formValues, resp.Cookies(), nil
}
