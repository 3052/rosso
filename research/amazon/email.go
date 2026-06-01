// post_email.go
package amazon

import (
   "fmt"
   "io"
   "net/http"
   "net/url"
   "os"
   "path/filepath"
   "strings"
)

func PostEmail(s *Session, action string, inputs map[string]string) (string, map[string]string, error) {
   data := url.Values{}
   //////////////////////////////////////////////////////////////////////////////
   for k, v := range inputs {
      data.Set(k, v)
   }
   data.Set("email", s.Email)
   req, err := http.NewRequest("POST", action, strings.NewReader(data.Encode()))
   if err != nil {
      return "", nil, err
   }
   req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
   req.Header.Set("Referer", "https://www.amazon.com/ap/signin")
   req.Header.Set("Sec-Fetch-Dest", "document")
   req.Header.Set("Sec-Fetch-Mode", "navigate")
   req.Header.Set("Sec-Fetch-Site", "same-origin")
   req.Header.Set("Upgrade-Insecure-Requests", "1")
   req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:140.0) Gecko/20100101 Firefox/140.0")
   //////////////////////////////////////////////////////////////////////////////
   resp, err := s.Client.Do(req)
   if err != nil {
      return "", nil, err
   }
   defer resp.Body.Close()

   body, _ := io.ReadAll(resp.Body)
   bodyStr := string(body)

   // If the password input field is not present, Amazon rejected the email (e.g. CAPTCHA, invalid email, rate limit)
   if !strings.Contains(bodyStr, `id="ap_password"`) {
      errDetails := CheckAmazonErrors(body)
      if errDetails == nil {
         errDetails = fmt.Errorf("did not reach password page")
      }
      errFile := filepath.Join(os.TempDir(), "error_post_email.html")
      os.WriteFile(errFile, body, 0644)
      return "", nil, fmt.Errorf("email submission failed: %v. Response saved to %s", errDetails, errFile)
   }

   nextAction, nextInputs := ExtractForm(bodyStr, "signIn")
   if nextAction == "" {
      errFile := filepath.Join(os.TempDir(), "error_post_email.html")
      os.WriteFile(errFile, body, 0644)
      return "", nil, fmt.Errorf("password form not found in response. Response saved to %s", errFile)
   }

   if strings.HasPrefix(nextAction, "/") {
      nextAction = "https://www.amazon.com" + nextAction
   }

   return nextAction, nextInputs, nil
}
