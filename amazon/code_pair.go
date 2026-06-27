package amazon

import (
   "bytes"
   "encoding/json"
   "fmt"
   "net/http"
   "strings"
)

// CodePair represents the public and private codes used for device linking.
type CodePair struct {
   PublicCode  string `json:"public_code"`
   PrivateCode string `json:"private_code"`
}

// CreateCodePair requests a public and private code pair for device linking.
func CreateCodePair() (*CodePair, error) {
   payload := map[string]any{
      "code_data": map[string]string{
         "domain":           "Device",
         "device_name":      DeviceName,
         "app_name":         "AIV",
         "app_version":      "3.12.0",
         "device_model":     "sdk_gphone_x86",
         "os_version":       "Android",
         "device_type":      DeviceTypeID, // from HAR: A2SNKIF736WF4T
         "device_serial":    DeviceID,     // from HAR: uuidb43bee409bd448cfb5ba3337bd241645
         "software_version": "999",
      },
   }
   body, err := json.Marshal(payload)
   if err != nil {
      return nil, err
   }
   req, err := http.NewRequest(
      "POST", HostAmazonAPI+"/auth/create/codepair",
      bytes.NewBuffer(body),
   )
   if err != nil {
      return nil, err
   }

   // Headers matching the HAR file
   req.Header.Set("User-Agent", "Android/google/sdk_gphone_x86/generic_x86_arm:11/RSR1.240422.006/12134477:userdebug/dev-keys, Ignition X/15.5.2026042820-android, Google")
   req.Header.Set("Accept-Encoding", "identity")
   req.Header.Set("content-type", "application/json")
   req.Header.Set("accept-language", "en_US")
   req.Header.Set("accept", "application/json")

   client := &http.Client{}
   resp, err := client.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()

   if resp.StatusCode != http.StatusOK {
      return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
   }

   // Decode directly into our new struct type
   var result CodePair
   if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
      return nil, err
   }

   return &result, nil
}

func (*CodePair) CachePath() string {
   return "rosso/amazon/CodePair"
}

// amazon.com/mytv
func (c *CodePair) String() string {
   var data strings.Builder
   data.WriteString("Please navigate to https://primevideo.com/ontv\n")
   data.WriteString("or https://amazon.com/code\n")
   data.WriteString("Enter the following code: ")
   data.WriteString(c.PublicCode)
   return data.String()
}
