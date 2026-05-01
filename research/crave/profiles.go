package crave

import (
   "encoding/json"
   "net/url"

   "41.neocities.org/maya"
)

type Profile struct {
   Id                string   `json:"id"`
   AccountId         string   `json:"accountId"`
   Nickname          string   `json:"nickname"`
   HasPin            bool     `json:"hasPin"`
   Master            bool     `json:"master"`
   Maturity          string   `json:"maturity"`
   Onboarded         bool     `json:"onboarded"`
   UiLanguage        string   `json:"uiLanguage"`
   PlaybackLanguages []string `json:"playbackLanguages"`
   LastModifiedDate  string   `json:"lastModifiedDate"`
   AvatarUrl         string   `json:"avatarUrl"`
}

func GetProfiles(account *AccountToken) ([]Profile, error) {
   endpoint := &url.URL{
      Scheme: "https",
      Host:   "account.bellmedia.ca",
      Path:   "/api/profile/v2/account/" + account.AccountId,
   }

   headers := map[string]string{
      "authorization": "Bearer " + account.AccessToken,
   }

   resp, err := maya.Get(endpoint, headers)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()

   var activeProfiles []Profile
   if err := json.NewDecoder(resp.Body).Decode(&activeProfiles); err != nil {
      return nil, err
   }

   return activeProfiles, nil
}
