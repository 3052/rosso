package crave

import (
   "encoding/json"
   "net/url"

   "41.neocities.org/maya"
)

type Session struct {
   AccessToken string `json:"access_token"`
}

func CreateSession(activeLogin *Login, activeProfile *Profile) (*Session, error) {
   endpoint := url.URL{
      Scheme: "https",
      Host:   "account.bellmedia.ca",
      Path:   "/api/login/v2.2",
   }

   values := url.Values{}
   values.Set("grant_type", "refresh_token")
   values.Set("profile_id", activeProfile.Id)
   values.Set("refresh_token", activeLogin.RefreshToken)

   headers := map[string]string{
      "content-type":  "application/x-www-form-urlencoded",
      "authorization": "Basic Y3JhdmUtd2ViOmRlZmF1bHQ=",
   }

   resp, err := maya.Post(&endpoint, headers, []byte(values.Encode()))
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()

   var targetSession Session
   if err := json.NewDecoder(resp.Body).Decode(&targetSession); err != nil {
      return nil, err
   }

   return &targetSession, nil
}
