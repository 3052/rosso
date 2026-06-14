package amazon

import (
   "bytes"
   "encoding/json"
   "fmt"
   "net/http"
)

// CreateCodePair requests a public and private code pair for device linking.
func CreateCodePair() (string, string, error) {
   url := "https://api.amazon.com/auth/create/codepair"

   payload := map[string]interface{}{
      "code_data": map[string]string{
         "domain":           "Device",
         "device_name":      "%FIRST_NAME%'s%DUPE_STRATEGY_1ST% sdk_gphone_x86",
         "app_name":         "AIV",
         "app_version":      "3.12.0",
         "device_model":     "sdk_gphone_x86",
         "os_version":       "Android",
         "device_type":      "A2SNKIF736WF4T",
         "device_serial":    "uuidcbb2f9705f13437e9e515622dce02106",
         "software_version": "999",
      },
   }

   body, err := json.Marshal(payload)
   if err != nil {
      return "", "", err
   }

   req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
   if err != nil {
      return "", "", err
   }

   req.Header.Set("User-Agent", "Android/google/sdk_gphone_x86/generic_x86_arm:11/RSR1.240422.006/12134477:userdebug/dev-keys, Ignition X/15.5.2026042820-android, Google")
   req.Header.Set("Content-Type", "application/json")
   req.Header.Set("Accept", "application/json")

   client := &http.Client{}
   resp, err := client.Do(req)
   if err != nil {
      return "", "", err
   }
   defer resp.Body.Close()

   if resp.StatusCode != http.StatusOK {
      return "", "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
   }

   var result struct {
      PublicCode  string `json:"public_code"`
      PrivateCode string `json:"private_code"`
   }

   if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
      return "", "", err
   }

   return result.PublicCode, result.PrivateCode, nil
}
