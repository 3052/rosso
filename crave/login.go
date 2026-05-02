package crave

import (
   "encoding/json"
   "net/url"

   "41.neocities.org/maya"
)

type AccountToken struct {
   AccessToken  string `json:"access_token"`
   RefreshToken string `json:"refresh_token"`
   AccountId    string `json:"account_id"`
   Jti          string `json:"jti"`
}

func PerformLogin(username string, password string) (*AccountToken, error) {
   endpoint := &url.URL{
      Scheme: "https",
      Host:   "account.bellmedia.ca",
      Path:   "/api/login/v2.1",
   }

   headers := map[string]string{
      "authorization": "Basic Y3JhdmUtd2ViOmRlZmF1bHQ=",
      "content-type":  "application/x-www-form-urlencoded",
   }

   values := url.Values{}
   values.Set("grant_type", "password")
   values.Set("password", password)
   values.Set("username", username)

   body := []byte(values.Encode())

   resp, err := maya.Post(endpoint, headers, body)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()

   account := &AccountToken{}
   if err := json.NewDecoder(resp.Body).Decode(account); err != nil {
      return nil, err
   }

   return account, nil
}
