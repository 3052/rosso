package crave

import (
   "encoding/json"
   "net/http"
   "net/url"
   "strings"
)

func fetch_account(username, password string) (*account, error) {
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
   result := &account{}
   err = json.NewDecoder(resp.Body).Decode(result)
   if err != nil {
      return nil, err
   }
   return result, nil
}

type account struct {
   RefreshToken string `json:"refresh_token"`
   AccessToken string `json:"access_token"`
}

func (a *account) String() string {
   var data strings.Builder
   data.WriteString("refresh token = ")
   data.WriteString(a.RefreshToken)
   data.WriteString("\naccess token = ")
   data.WriteString(a.AccessToken)
   return data.String()
}
