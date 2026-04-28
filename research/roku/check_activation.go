package roku

import (
   "encoding/json"
   "net/url"

   "41.neocities.org/maya"
)

type ActivationDetails struct {
   Code      string       `json:"code"`
   Token     ContentToken `json:"token"`
   CreatedAt int64        `json:"createdAt"`
   Platform  string       `json:"platform"`
   Status    string       `json:"status"`
   Profiles  []Profile    `json:"profiles"`
}

type Profile struct {
   Id      string `json:"id"`
   IsKids  bool   `json:"isKids"`
   IsOwner bool   `json:"isOwner"`
}

func CheckActivation(userToken ContentToken, activationData Activation) (*ActivationDetails, error) {
   endpoint := &url.URL{
      Scheme: "https",
      Host:   "googletv.web.roku.com",
      Path:   "/api/v1/account/activation/" + activationData.Code,
   }

   headers := map[string]string{
      "user-agent":           "trc-googletv; production; 0",
      "x-roku-content-token": string(userToken),
   }

   resp, err := maya.Get(endpoint, headers)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()

   var details ActivationDetails
   if err := json.NewDecoder(resp.Body).Decode(&details); err != nil {
      return nil, err
   }

   return &details, nil
}
