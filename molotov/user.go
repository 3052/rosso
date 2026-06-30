// user.go
package molotov

import (
   "encoding/json"
   "fmt"
   "net/http"
)

type UserResponse struct {
   Data struct {
      ID       string `json:"id"`
      Profiles []struct {
         ID string `json:"id"`
      } `json:"profiles"`
   } `json:"data"`
}

// GetUser fetches the user profile using the token from the Signin response.
func GetUser(signinResp *SigninResponse) (*UserResponse, error) {
   url := "https://api-eu.fubo.tv/user"

   req, err := http.NewRequest("GET", url, nil)
   if err != nil {
      return nil, err
   }

   req.Header.Set("Authorization", "Bearer "+signinResp.Payload.AccessToken)
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

   var userResp UserResponse
   if err := json.NewDecoder(resp.Body).Decode(&userResp); err != nil {
      return nil, err
   }

   if len(userResp.Data.Profiles) == 0 {
      return nil, fmt.Errorf("no profiles found for user")
   }

   return &userResp, nil
}
