package kanopy

import (
   "encoding/json"
   "net/url"

   "41.neocities.org/maya"
)

type LoginResponse struct {
   JWT       string `json:"jwt"`
   VisitorID string `json:"visitorId"`
   UserID    int    `json:"userId"`
   WebshopID int    `json:"webshopId"`
   UserRole  string `json:"userRole"`
}

func Login(email, password string) (*LoginResponse, error) {
   loginURL := &url.URL{
      Scheme: "https",
      Host:   "www.kanopy.com",
      Path:   "/kapi/login",
   }

   payload := map[string]any{
      "credentialType": "email",
      "emailUser": map[string]string{
         "email":    email,
         "password": password,
      },
   }

   body, err := json.Marshal(payload)
   if err != nil {
      return nil, err
   }

   headers := map[string]string{
      "content-type": "application/json",
      "user-agent":   "!",
   }

   resp, err := maya.Post(loginURL, headers, body)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()

   var loginResp LoginResponse
   if err := json.NewDecoder(resp.Body).Decode(&loginResp); err != nil {
      return nil, err
   }

   return &loginResp, nil
}
