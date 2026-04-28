package roku

import (
   "encoding/json"
   "net/url"

   "41.neocities.org/maya"
)

type ActivationStatus struct {
   Code      string    `json:"code"`
   Token     string    `json:"token"`
   CreatedAt int64     `json:"createdAt"`
   Profiles  []Profile `json:"profiles"`
   Platform  string    `json:"platform"`
   Status    string    `json:"status"`
}

type Profile struct {
   Id      string `json:"id"`
   IsKids  bool   `json:"isKids"`
   IsOwner bool   `json:"isOwner"`
}

func GetActivationStatus(token *AccountToken, activation *AccountActivation) (*ActivationStatus, error) {
   target := &url.URL{
      Scheme: "https",
      Host:   "googletv.web.roku.com",
      Path:   "/api/v1/account/activation/" + activation.Code,
   }
   headers := map[string]string{
      "user-agent":           "trc-googletv; production; 0",
      "x-roku-content-token": token.AuthToken,
   }

   resp, err := maya.Get(target, headers)
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
