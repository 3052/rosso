package crave

import (
   "encoding/json"
   "net/url"

   "41.neocities.org/maya"
)

type ProfileToken struct {
   AccessToken  string `json:"access_token"`
   RefreshToken string `json:"refresh_token"`
   Scope        string `json:"scope"`
   TokenType    string `json:"token_type"`
   ExpiresIn    int    `json:"expires_in"`
}

func SwitchProfile(account *AccountToken, activeProfile *Profile) (*ProfileToken, error) {
   endpoint := &url.URL{
      Scheme: "https",
      Host:   "account.bellmedia.ca",
      Path:   "/api/login/v2.2",
   }

   headers := map[string]string{
      "authorization": "Basic Y3JhdmUtd2ViOmRlZmF1bHQ=",
      "content-type":  "application/x-www-form-urlencoded",
   }

   values := url.Values{}
   values.Set("grant_type", "refresh_token")
   values.Set("profile_id", activeProfile.Id)
   values.Set("refresh_token", account.RefreshToken)

   body := []byte(values.Encode())

   resp, err := maya.Post(endpoint, headers, body)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()

   token := &ProfileToken{}
   if err := json.NewDecoder(resp.Body).Decode(token); err != nil {
      return nil, err
   }

   return token, nil
}
