package roku

import (
   "encoding/json"
   "net/url"

   "41.neocities.org/maya"
)

type ContentToken string

type AccountToken struct {
   AuthToken  ContentToken `json:"authToken"`
   IsLoggedIn bool         `json:"isLoggedIn"`
}

func FetchAccountToken(userToken ContentToken) (*AccountToken, error) {
   endpoint := &url.URL{
      Scheme: "https",
      Host:   "googletv.web.roku.com",
      Path:   "/api/v1/account/token",
   }

   headers := map[string]string{
      "user-agent": "trc-googletv; production; 0",
   }
   if userToken != "" {
      headers["x-roku-content-token"] = string(userToken)
   }

   resp, err := maya.Get(endpoint, headers)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()

   var account AccountToken
   if err := json.NewDecoder(resp.Body).Decode(&account); err != nil {
      return nil, err
   }

   return &account, nil
}
