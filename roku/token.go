package roku

import (
   "encoding/json"
   "net/url"

   "41.neocities.org/maya"
)

type AccountToken struct {
   AuthToken  string `json:"authToken"`
   IsLoggedIn bool   `json:"isLoggedIn"`
   Ip         string `json:"ip"`
   Rida       string `json:"rida"`
}

func GetAccountToken(status *ActivationStatus) (*AccountToken, error) {
   target := &url.URL{
      Scheme: "https",
      Host:   "googletv.web.roku.com",
      Path:   "/api/v1/account/token",
   }
   headers := map[string]string{
      "user-agent": "trc-googletv; production; 0",
   }
   if status != nil {
      headers["x-roku-content-token"] = status.Token
   }

   resp, err := maya.Get(target, headers)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()

   var token AccountToken
   if err := json.NewDecoder(resp.Body).Decode(&token); err != nil {
      return nil, err
   }
   return &token, nil
}
