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

//   "scope": [
//     "account:write",
//     "default",
//     "offline_download:10",
//     "password_token",
//     "subscription:crave_total,cravep,cravetv,free,se"
//   ]
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
