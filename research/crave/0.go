package crave

import (
   "encoding/json"
   "net/http"
   "net/url"
   "strings"
)

func (z *zero) String() string {
   var data strings.Builder
   data.WriteString("refresh token = ")
   data.WriteString(z.RefreshToken)
   data.WriteString("\naccess token = ")
   data.WriteString(z.AccessToken)
   return data.String()
}

// {
//   "localization": "en-CA",
//   "user_name": "",
//   "brand_policies": [
//     "casting:airplay",
//     "casting:chromecast",
//     "device:5",
//     "offline_download",
//     "platform:android",
//     "platform:android_tv",
//     "platform:fire_tv",
//     "platform:hisense",
//     "platform:ios",
//     "platform:lg_tv",
//     "platform:ps5",
//     "platform:roku",
//     "platform:samsung_tv",
//     "platform:sony_ps4",
//     "platform:stb",
//     "platform:tvos",
//     "platform:web",
//     "platform:x1",
//     "platform:xbox_one",
//     "playback_quality:4k",
//     "playback_quality:hd",
//     "playback_quality:sd",
//     "stream_concurrency:4",
//     "subscription:crave_total",
//     "subscription:cravep",
//     "subscription:cravetv",
//     "subscription:free",
//     "subscription:se"
//   ],
//   "creation_date": 1774203235086,
//   "ais_id": null,
//   "authorities": [
//     "REGULAR_USER"
//   ],
//   "client_id": "crave-web",
//   "brand_id": "1d72d990cb765de7e4211111",
//   "account_id": "69967da9c93eef5df20f8712",
//   "profile_id": null,
//   "scope": [
//     "account:write",
//     "default",
//     "offline_download:10",
//     "password_token",
//     "subscription:crave_total,cravep,cravetv,free,se"
//   ],
//   "exp": 1774217635,
//   "iat": 1774203235,
//   "jti": "84246206-53ea-4d00-86a8-88ae6dba7451",
//   "account": {
//     "status": "ACTIVE"
//   }
// }
type zero struct {
   RefreshToken string `json:"refresh_token"`
   AccessToken string `json:"access_token"`
}

func fetch_zero(username, password string) (*zero, error) {
   data := url.Values{
      "grant_type": {"password"},
      "password": {password},
      "username": {username},
   }.Encode()
   req, err := http.NewRequest(
      "POST", "https://account.bellmedia.ca/api/login/v2.1",
      strings.NewReader(data),
   )
   if err != nil {
      return nil, err
   }
   req.Header.Set("content-type", "application/x-www-form-urlencoded")
   req.SetBasicAuth("crave-web", "default")
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   result := &zero{}
   err = json.NewDecoder(resp.Body).Decode(result)
   if err != nil {
      return nil, err
   }
   return result, nil
}
