package amazon

import (
   "bytes"
   "encoding/json"
   "fmt"
   "io"
   "net/http"
)

// InitiateMDSO informs the authentication proxy about the public code to initiate linking.
func InitiateMDSO(publicCode string) error {
   url := "https://s0s7.api.amazonvideo.com/cdp/authproxy/mdso/initiate"

   req, err := http.NewRequest("POST", url, nil)
   if err != nil {
      return err
   }

   q := req.URL.Query()
   q.Add("deviceTypeID", "A2SNKIF736WF4T")
   q.Add("deviceID", "uuidcbb2f9705f13437e9e515622dce02106")
   q.Add("firmware", "google/sdk_gphone_x86/generic_x86_arm:11/RSR1.240422.006/12134477:userdebug/dev-keys")
   q.Add("manufacturer", "Google")
   q.Add("chipset", "goldfish_x86")
   q.Add("model", "sdk_gphone_x86")
   q.Add("operatingSystem", "Android")
   q.Add("uxLocale", "en_US")
   req.URL.RawQuery = q.Encode()

   payload := map[string]interface{}{
      "publicCode":         publicCode,
      "generalDeviceGroup": true,
   }

   body, err := json.Marshal(payload)
   if err != nil {
      return err
   }

   req.Body = io.NopCloser(bytes.NewBuffer(body))
   req.ContentLength = int64(len(body))

   req.Header.Set("User-Agent", "Android/google/sdk_gphone_x86/generic_x86_arm:11/RSR1.240422.006/12134477:userdebug/dev-keys, Ignition X/15.5.2026042820-android, Google")
   req.Header.Set("Content-Type", "application/json")
   req.Header.Set("Accept", "application/json")

   client := &http.Client{}
   resp, err := client.Do(req)
   if err != nil {
      return err
   }
   defer resp.Body.Close()

   if resp.StatusCode != http.StatusOK {
      return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
   }

   return nil
}
