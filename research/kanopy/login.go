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

type LoginPayload struct {
   CredentialType string    `json:"credentialType"`
   EmailUser      EmailUser `json:"emailUser"`
}

type LoginResponse struct {
   Jwt       string `json:"jwt"`
   VisitorId string `json:"visitorId"`
   UserId    int    `json:"userId"`
}

func Login(emailUser *EmailUser) (*LoginResponse, error) {
   payload := LoginPayload{
      CredentialType: "email",
      EmailUser:      *emailUser,
   }
   bodyBytes, marshalError := json.Marshal(payload)
   if marshalError != nil {
      return nil, marshalError
   }

   targetUrl, parseError := url.Parse("https://www.kanopy.com/kapi/login")
   if parseError != nil {
      return nil, parseError
   }

   headers := map[string]string{
      "content-type": "application/json",
      "user-agent":   "!",
   }

   resp, requestError := maya.Post(targetUrl, headers, bodyBytes)
   if requestError != nil {
      return nil, requestError
   }
   defer resp.Body.Close()

   var loginResp LoginResponse
   decodeError := json.NewDecoder(resp.Body).Decode(&loginResp)
   if decodeError != nil {
      return nil, decodeError
   }
   return &loginResp, nil
}
