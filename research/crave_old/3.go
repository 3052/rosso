package crave

import (
   "encoding/json"
   "net/http"
   "net/url"
   "strings"
)

func (a *account) three() error {
   data := url.Values{
      "grant_type": {"refresh_token"},
      "refresh_token": {a.RefreshToken},
   }.Encode()
   req, err := http.NewRequest(
      "POST", "https://account.bellmedia.ca/api/login/v2.2",
      strings.NewReader(data),
   )
   if err != nil {
      return err
   }
   req.Header.Set("content-type", "application/x-www-form-urlencoded")
   req.Header.Set("authorization", "Basic Y3JhdmUtd2ViOmRlZmF1bHQ=")
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return err
   }
   defer resp.Body.Close()
   return json.NewDecoder(resp.Body).Decode(a)
}
