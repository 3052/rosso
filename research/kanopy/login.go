// File: login.go
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

type loginRequest struct {
   CredentialType string     `json:"credentialType"`
   EmailUser      *EmailUser `json:"emailUser"`
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

func Login(emailUser *EmailUser) (*LoginResponse, error) {
   reqURL, err := url.Parse("https://www.kanopy.com/kapi/login")
   if err != nil {
      return nil, err
   }

   reqBody := loginRequest{
      CredentialType: "email",
      EmailUser:      emailUser,
   }

   payload, err := json.Marshal(reqBody)
   if err != nil {
      return nil, err
   }

   headers := map[string]string{
      "content-type": "application/json",
      "user-agent":   "!",
   }

   resp, err := maya.Post(reqURL, headers, payload)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()

   var result LoginResponse
   if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
      return nil, err
   }

   return &result, nil
}
