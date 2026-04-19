// login.go
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

type LoginResponse struct {
   JWT               string `json:"jwt"`
   VisitorID         string `json:"visitorId"`
   UserID            int    `json:"userId"`
   KanopyKidsEnabled bool   `json:"kanopyKidsEnabled"`
   WebshopID         int    `json:"webshopId"`
   WebshopCode       string `json:"webshopCode"`
   UserRole          string `json:"userRole"`
}

func Login(email, password string) (*LoginResponse, error) {
   reqData := LoginRequest{
      CredentialType: "email",
      EmailUser: EmailUser{
         Email:    email,
         Password: password,
      },
   }
   body, err := json.Marshal(reqData)
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

   resp, err := maya.Post(targetUrl, headers, body)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()

   var loginResp LoginResponse
   err = json.NewDecoder(resp.Body).Decode(&loginResp)
   if err != nil {
      return nil, err
   }

   return &loginResp, nil
}
