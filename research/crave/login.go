package crave

import (
   "encoding/json"
   "net/url"

   "41.neocities.org/maya"
)

type Login struct {
   AccessToken  string `json:"access_token"`
   RefreshToken string `json:"refresh_token"`
   AccountId    string `json:"account_id"`
}

func LoginAccount(username string, password string) (*Login, error) {
   endpoint := url.URL{
      Scheme: "https",
      Host:   "account.bellmedia.ca",
      Path:   "/api/login/v2.1",
   }

   values := url.Values{}
   values.Set("grant_type", "password")
   values.Set("password", password)
   values.Set("username", username)

   headers := map[string]string{
      "content-type":  "application/x-www-form-urlencoded",
      "authorization": "Basic Y3JhdmUtd2ViOmRlZmF1bHQ=",
   }

   resp, err := maya.Post(&endpoint, headers, []byte(values.Encode()))
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()

   var targetLogin Login
   if err := json.NewDecoder(resp.Body).Decode(&targetLogin); err != nil {
      return nil, err
   }

   return &targetLogin, nil
}
