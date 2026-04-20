// login.go
package kanopy

import (
   "encoding/json"
   "fmt"
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
   JWT               string `json:"jwt"`
   VisitorID         string `json:"visitorId"`
   UserID            int    `json:"userId"`
   KanopyKidsEnabled bool   `json:"kanopyKidsEnabled"`
   WebshopID         int    `json:"webshopId"`
   WebshopCode       string `json:"webshopCode"`
   UserRole          string `json:"userRole"`
}

func Login(req *LoginRequest) (*LoginResponse, error) {
   targetURL, err := url.Parse("https://www.kanopy.com/kapi/login")
   if err != nil {
      return nil, err
   }

   bodyBytes, err := json.Marshal(req)
   if err != nil {
      return nil, err
   }

   headers := map[string]string{
      "content-type": "application/json",
      "user-agent":   "!",
   }

   resp, err := maya.Post(targetURL, headers, bodyBytes)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()

   if resp.StatusCode != 200 {
      return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
   }

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
