package amazon

import (
   "fmt"
   "io"
   "net/http"
)

func TriggerOTP(client *http.Client, otpUrl string) (string, map[string]string, error) {
   req, err := http.NewRequest("GET", otpUrl, nil)
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

   actionUrl, hiddenParams := extractFormActionAndHiddenInputs(string(body), otpUrl)
   return actionUrl, hiddenParams, nil
}
