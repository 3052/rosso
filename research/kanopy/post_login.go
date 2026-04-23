package kanopy

import (
   "encoding/json"
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

type Login struct {
   Jwt    string `json:"jwt"`
   UserId int    `json:"userId"`
}

func PostLogin(email string, password string) (*Login, error) {
   endpoint := &url.URL{
      Scheme: "https",
      Host:   "www.kanopy.com",
      Path:   "/kapi/login",
   }

   payload := LoginRequest{
      CredentialType: "email",
      EmailUser: EmailUser{
         Email:    email,
         Password: password,
      },
   }

   body, err := json.Marshal(payload)
   if err != nil {
      return nil, err
   }

   resp, err := maya.Post(endpoint, nil, body)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()

   var loginData Login
   if err := json.NewDecoder(resp.Body).Decode(&loginData); err != nil {
      return nil, err
   }

   return &loginData, nil
}
