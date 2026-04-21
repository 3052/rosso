package kanopy

import (
   "encoding/json"
   "net/url"

   "41.neocities.org/maya"
)

type LoginResponse struct {
   Jwt       string `json:"jwt"`
   VisitorId string `json:"visitorId"`
   UserId    int    `json:"userId"`
   WebshopId int    `json:"webshopId"`
   UserRole  string `json:"userRole"`
}

func Login(email, password string) (*LoginResponse, error) {
   loginUrl := &url.URL{
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

   resp, err := maya.Post(loginUrl, headers, body)
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
