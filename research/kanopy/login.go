package kanopy

import (
   "encoding/json"
   "io"
   "net/url"

   "41.neocities.org/maya"
)

type LoginResponse struct {
   JWT               string `json:"jwt"`
   VisitorID         string `json:"visitorId"`
   UserID            int64  `json:"userId"`
   KanopyKidsEnabled bool   `json:"kanopyKidsEnabled"`
   WebshopID         int    `json:"webshopId"`
   WebshopCode       string `json:"webshopCode"`
   UserRole          string `json:"userRole"`
}

type loginRequest struct {
   CredentialType string    `json:"credentialType"`
   EmailUser      emailUser `json:"emailUser"`
}

type emailUser struct {
   Email    string `json:"email"`
   Password string `json:"password"`
}

func Login(email string, password string) (*LoginResponse, error) {
   target := &url.URL{
      Scheme: "https",
      Host:   "www.kanopy.com",
      Path:   "/kapi/login",
   }

   headers := map[string]string{
      "content-type": "application/json",
      "user-agent":   "!",
   }

   reqBody := loginRequest{
      CredentialType: "email",
      EmailUser: emailUser{
         Email:    email,
         Password: password,
      },
   }

   bodyBytes, err := json.Marshal(reqBody)
   if err != nil {
      return nil, err
   }

   resp, err := maya.Post(target, headers, bodyBytes)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()

   respBytes, err := io.ReadAll(resp.Body)
   if err != nil {
      return nil, err
   }

   var loginResp LoginResponse
   if err := json.Unmarshal(respBytes, &loginResp); err != nil {
      return nil, err
   }

   return &loginResp, nil
}
