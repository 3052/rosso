// user.go
package molotov

import (
   "encoding/json"
   "fmt"
   "net/http"
)

func (*UserResponse) CachePath() string {
   return "rosso/molotov/UserResponse"
}

type UserResponse struct {
   ID       string    `json:"id"`
   Profiles []Profile `json:"profiles"`
}

type Profile struct {
   ID string `json:"id"`
}

// GetUser fetches the user profile using the token from the Signin response.
func GetUser(signinResp *SigninResponse) (*UserResponse, error) {
   url := "https://api-eu.fubo.tv/user"

   req, err := http.NewRequest("GET", url, nil)
   if err != nil {
      return nil, err
   }

   // Accessing the unwrapped field directly
   req.Header.Set("Authorization", "Bearer "+signinResp.AccessToken)
   req.Header.Set("x-device-id", DeviceID)
   req.Header.Set("x-device-group", "desktop")
   req.Header.Set("x-client-version", "6.12.0")
   req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:140.0) Gecko/20100101 Firefox/140.0")

   client := &http.Client{}
   resp, err := client.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()

   if resp.StatusCode != http.StatusOK {
      return nil, fmt.Errorf("get user failed with status: %d", resp.StatusCode)
   }

   // Unwrap the "data" envelope layer
   var envelope struct {
      Data UserResponse `json:"data"`
   }
   if err := json.NewDecoder(resp.Body).Decode(&envelope); err != nil {
      return nil, err
   }

   if len(envelope.Data.Profiles) == 0 {
      return nil, fmt.Errorf("no profiles found for user")
   }

   return &envelope.Data, nil
}
