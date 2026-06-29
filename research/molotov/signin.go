// signin.go
package molotov

import (
   "bytes"
   "encoding/json"
   "fmt"
   "net/http"
)

// Signin performs the authentication request and returns the access token.
func Signin(username, password string) (string, error) {
   url := "https://api-eu.fubo.tv/v2/signin"

   reqBody, err := json.Marshal(SigninRequest{
      Username: username,
      Password: password,
   })
   if err != nil {
      return "", err
   }

   req, err := http.NewRequest("PUT", url, bytes.NewBuffer(reqBody))
   if err != nil {
      return "", err
   }

   req.Header.Set("Content-Type", "application/json")
   req.Header.Set("x-device-id", DeviceID)
   req.Header.Set("x-device-group", "desktop")
   req.Header.Set("x-device-type", "desktop")
   req.Header.Set("x-client-version", "6.12.0")
   req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:140.0) Gecko/20100101 Firefox/140.0")

   client := &http.Client{}
   resp, err := client.Do(req)
   if err != nil {
      return "", err
   }
   defer resp.Body.Close()

   if resp.StatusCode != http.StatusOK {
      return "", fmt.Errorf("signin failed with status: %d", resp.StatusCode)
   }

   var signinResp SigninResponse
   if err := json.NewDecoder(resp.Body).Decode(&signinResp); err != nil {
      return "", err
   }

   return signinResp.Payload.AccessToken, nil
}

type SigninRequest struct {
   Username string `json:"username"`
   Password string `json:"password"`
}

type SigninResponse struct {
   Payload struct {
      AccessToken string `json:"access_token"`
   } `json:"payload"`
}
