package crave

import (
   "encoding/json"
   "net/http"
   "net/url"
   "strings"
)

func two(magic_link_token string) (*account, error) {
   data := url.Values{
      "grant_type":[]string{"magic_link_token"},
      "magic_link_token":[]string{magic_link_token},
   }.Encode()
   req, err := http.NewRequest(
      "POST", "https://account.bellmedia.ca/api/login/v2.2",
      strings.NewReader(data),
   )
   if err != nil {
      return nil, err
   }
   req.Header.Add("Authorization", "Basic Y3JhdmUtd2ViOmRlZmF1bHQ=")
   req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
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
