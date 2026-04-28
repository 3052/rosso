package roku

import (
   "encoding/json"
   "net/url"

   "41.neocities.org/maya"
)

type ActivationStatus struct {
   Code      ActivationCode `json:"code"`
   Token     ContentToken   `json:"token"`
   CreatedAt int64          `json:"createdAt"`
   Profiles  []Profile      `json:"profiles"`
   Platform  string         `json:"platform"`
   Status    string         `json:"status"`
}

type Profile struct {
   Id      string `json:"id"`
   IsKids  bool   `json:"isKids"`
   IsOwner bool   `json:"isOwner"`
}

func FetchActivationStatus(token ContentToken, code ActivationCode) (*ActivationStatus, error) {
   endpoint := &url.URL{
      Scheme: "https",
      Host:   "googletv.web.roku.com",
      Path:   "/api/v1/account/activation/" + string(code),
   }

   headers := map[string]string{
      "user-agent":           "trc-googletv; production; 0",
      "x-roku-content-token": string(token),
   }

   resp, err := maya.Get(endpoint, headers)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()

   var status ActivationStatus
   if err := json.NewDecoder(resp.Body).Decode(&status); err != nil {
      return nil, err
   }

   return &status, nil
}
