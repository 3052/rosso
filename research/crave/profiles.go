package crave

import (
   "encoding/json"
   "net/url"

   "41.neocities.org/maya"
)

type Profile struct {
   Id string `json:"id"`
}

func GetProfiles(activeLogin *Login) ([]Profile, error) {
   endpoint := url.URL{
      Scheme: "https",
      Host:   "account.bellmedia.ca",
      Path:   "/api/profile/v2/account/" + activeLogin.AccountId,
   }

   headers := map[string]string{
      "authorization": "Bearer " + activeLogin.AccessToken,
   }

   resp, err := maya.Get(&endpoint, headers)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()

   var targetProfiles []Profile
   if err := json.NewDecoder(resp.Body).Decode(&targetProfiles); err != nil {
      return nil, err
   }

   return targetProfiles, nil
}
