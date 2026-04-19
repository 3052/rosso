// login.go
package kanopy

import (
   "encoding/json"
   "io"
   "net/url"

   "41.neocities.org/maya"
)

type EmailUser struct {
   Email    string `json:"email"`
   Password string `json:"password"`
}

type LoginRequest struct {
   CredentialType string    `json:"credentialType"`
   EmailUser      EmailUser `json:"emailUser"`
}

type LoginResponse struct {
   JWT    string `json:"jwt"`
   UserId int    `json:"userId"`
}

func Login(req *LoginRequest) (*LoginResponse, error) {
   reqBody, err := json.Marshal(req)
   if err != nil {
      return nil, err
   }

   targetUrl, err := url.Parse("https://www.kanopy.com/kapi/login")
   if err != nil {
      return nil, err
   }

   headers := map[string]string{
      "content-type": "application/json",
   }

   resp, err := maya.Post(targetUrl, headers, reqBody)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()

   respBody, err := io.ReadAll(resp.Body)
   if err != nil {
      return nil, err
   }

   var loginResp LoginResponse
   if err := json.Unmarshal(respBody, &loginResp); err != nil {
      return nil, err
   }

   return &loginResp, nil
}
