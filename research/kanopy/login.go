package kanopy

import (
   "encoding/json"
   "io"
   "net/url"

   "41.neocities.org/maya"
)

type LoginRequestEmailUser struct {
   Email    string `json:"email"`
   Password string `json:"password"`
}

type LoginRequest struct {
   CredentialType string                `json:"credentialType"`
   EmailUser      LoginRequestEmailUser `json:"emailUser"`
}

type LoginResponse struct {
   JWT       string `json:"jwt"`
   VisitorId string `json:"visitorId"`
   UserId    int    `json:"userId"`
}

func Login(email string, password string) (*LoginResponse, error) {
   targetUrl, err := url.Parse("https://www.kanopy.com/kapi/login")
   if err != nil {
      return nil, err
   }

   requestPayload := LoginRequest{
      CredentialType: "email",
      EmailUser: LoginRequestEmailUser{
         Email:    email,
         Password: password,
      },
   }

   bodyBytes, err := json.Marshal(requestPayload)
   if err != nil {
      return nil, err
   }

   requestHeaders := map[string]string{
      "content-type": "application/json",
   }

   response, err := maya.Post(targetUrl, requestHeaders, bodyBytes)
   if err != nil {
      return nil, err
   }
   defer response.Body.Close()

   responseBytes, err := io.ReadAll(response.Body)
   if err != nil {
      return nil, err
   }

   var loginResponse LoginResponse
   err = json.Unmarshal(responseBytes, &loginResponse)
   if err != nil {
      return nil, err
   }

   return &loginResponse, nil
}
