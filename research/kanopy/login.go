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
   Jwt               string `json:"jwt"`
   VisitorId         string `json:"visitorId"`
   UserId            int    `json:"userId"`
   KanopyKidsEnabled bool   `json:"kanopyKidsEnabled"`
   WebshopId         int    `json:"webshopId"`
   WebshopCode       string `json:"webshopCode"`
   UserRole          string `json:"userRole"`
}

func LoginUser(email string, password string) (*LoginResponse, error) {
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

   respBody, err := io.ReadAll(resp.Body)
   if err != nil {
      return nil, err
   }

   var login LoginResponse
   if err := json.Unmarshal(respBody, &login); err != nil {
      return nil, err
   }

   return &login, nil
}
