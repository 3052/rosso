package amazon

import (
   "bytes"
   "encoding/json"
   "fmt"
   "net/http"
)

func (*TokenPair) CachePath() string {
   return "rosso/amazon/TokenPair"
}

// TokenPair represents the access and refresh tokens returned upon successful registration.
type TokenPair struct {
   AccessToken  string `json:"access_token"`
   RefreshToken string `json:"refresh_token"`
}

// PollRegister attempts to register the device. This should typically be called in a loop
// until it returns success (after the user links the device on the web).
func PollRegister(publicCode, privateCode string) (*TokenPair, error) {
   url := "https://api.amazon.com/auth/register"

   payload := map[string]interface{}{
      "auth_data": map[string]interface{}{
         "code_pair": map[string]string{
            "public_code":  publicCode,
            "private_code": privateCode,
         },
      },
      "registration_data": map[string]string{
         "domain":           "Device",
         "device_name":      "%FIRST_NAME%'s%DUPE_STRATEGY_1ST% " + DeviceModel,
         "app_name":         "AIV",
         "app_version":      "3.12.0",
         "device_model":     DeviceModel,
         "os_version":       DeviceOS,
         "device_type":      DeviceTypeID,
         "device_serial":    DeviceID,
         "software_version": "999",
      },
      "requested_token_type": []string{"bearer"},
   }

   body, err := json.Marshal(payload)
   if err != nil {
      return nil, err
   }

   req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
   if err != nil {
      return nil, err
   }

   req.Header.Set("User-Agent", UserAgent)
   req.Header.Set("Content-Type", "application/json")
   req.Header.Set("Accept", "application/json")

   client := &http.Client{}
   resp, err := client.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()

   if resp.StatusCode == http.StatusUnauthorized {
      return nil, fmt.Errorf("authorization pending/unauthorized")
   } else if resp.StatusCode != http.StatusOK {
      return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
   }

   // We can embed the new TokenPair struct directly into our anonymous decoder struct
   var result struct {
      Response struct {
         Success struct {
            Tokens struct {
               Bearer TokenPair `json:"bearer"`
            } `json:"tokens"`
         } `json:"success"`
      } `json:"response"`
   }

   if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
      return nil, err
   }

   bearer := result.Response.Success.Tokens.Bearer
   return &bearer, nil
}
