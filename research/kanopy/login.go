// file: login.go
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
   targetUrl := &url.URL{
      Scheme: "https",
      Host:   "www.kanopy.com",
      Path:   "/kapi/login",
   }

   reqBody := LoginRequest{
      CredentialType: "email",
      EmailUser:      emailUser,
   }

   jsonData, err := json.Marshal(reqBody)
   if err != nil {
      return nil, err
   }

   headers := map[string]string{
      "content-type": "application/json",
      "user-agent":   "!",
   }

   resp, err := maya.Post(targetUrl, headers, jsonData)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()

   bodyBytes, err := io.ReadAll(resp.Body)
   if err != nil {
      return nil, err
   }

   var loginResp LoginResponse
   if err := json.Unmarshal(bodyBytes, &loginResp); err != nil {
      return nil, err
   }

   return &loginResp, nil
}
