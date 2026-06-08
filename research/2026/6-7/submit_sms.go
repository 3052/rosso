// submit_sms.go
package amazon

import (
   "fmt"
   "net/http"
   "net/url"
   "strings"
)

func SubmitSMS(client *http.Client, actionUrl string, hiddenParams map[string]string, phoneNumber string) (string, error) {
   data := url.Values{}
   for k, v := range hiddenParams {
      data.Set(k, v)
   }
   // Per the trace: submit the phone number, leave password blank, set signInWithOTP to true
   data.Set("email", phoneNumber)
   data.Set("password", "")
   data.Set("signInWithOTP", "true")

   req, err := http.NewRequest("POST", actionUrl, strings.NewReader(data.Encode()))
   if err != nil {
      return "", fmt.Errorf("failed to create request: %w", err)
   }
   req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
   req.Header.Set("User-Agent", "Mozilla/5.0 (Linux; Android 11; sdk_gphone_x86_64 Build/RSR1.240422.006; wv) AppleWebKit/537.36")
   req.Header.Set("x-requested-with", "com.amazon.avod.thirdpartyclient")

   resp, err := client.Do(req)
   if err != nil {
      return "", fmt.Errorf("request failed: %w", err)
   }
   defer resp.Body.Close()

   if resp.StatusCode == http.StatusFound || resp.StatusCode == http.StatusMovedPermanently {
      loc := resp.Header.Get("Location")
      if !strings.HasPrefix(loc, "http") {
         loc = "https://www.amazon.com" + loc
      }
      return loc, nil
   }

   return "", fmt.Errorf("expected 302 redirect to OTP trigger, got %d", resp.StatusCode)
}
