package amazon

import (
   "bytes"
   "encoding/json"
   "fmt"
   "net/http"
)

// TokenPair represents the access and refresh tokens returned upon successful
// registration
type TokenPair struct {
   AccessToken  string `json:"access_token"`
   RefreshToken string `json:"refresh_token"`
}

// PollRegister attempts to register the device. This should typically be called in a loop
// until it returns success (after the user links the device on the web).
func PollRegister(publicCode, privateCode string) (*TokenPair, error) {
   payload := map[string]any{
      "auth_data": map[string]any{
         "code_pair": map[string]string{
            "public_code":  publicCode,
            "private_code": privateCode,
         },
      },
      "registration_data": map[string]string{
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
      "requested_token_type": []string{"bearer"},
   }
   body, err := json.Marshal(payload)
   if err != nil {
      return nil, err
   }
   req, err := http.NewRequest(
      "POST", HostAmazonAPI+"/auth/register", bytes.NewBuffer(body),
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
   var result struct {
      Response struct {
         Success struct {
            Tokens struct {
               Bearer TokenPair `json:"bearer"`
            } `json:"tokens"`
         } `json:"success"`
         Error struct {
            Code    string `json:"code"`
            Message string `json:"message"`
         } `json:"error"`
      } `json:"response"`
   }
   if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
      return nil, err
   }
   if result.Response.Error.Code != "" {
      return nil, fmt.Errorf("amazon API error: %s - %s", result.Response.Error.Code, result.Response.Error.Message)
   }
   if resp.StatusCode != http.StatusOK {
      return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
   }
   bearer := result.Response.Success.Tokens.Bearer
   return &bearer, nil
}

func (*TokenPair) CachePath() string {
   return "rosso/amazon/TokenPair"
}
