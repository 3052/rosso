package roku

import (
   "encoding/json"
   "net/url"

   "41.neocities.org/maya"
)

type AccountActivation struct {
   Code string `json:"code"`
}

func CreateAccountActivation(token *AccountToken) (*AccountActivation, error) {
   target := &url.URL{
      Scheme: "https",
      Host:   "googletv.web.roku.com",
      Path:   "/api/v1/account/activation",
   }
   headers := map[string]string{
      "content-type":         "application/json",
      "user-agent":           "trc-googletv; production; 0",
      "x-roku-content-token": token.AuthToken,
   }

   reqBody, err := json.Marshal(map[string]string{
      "platform": "googletv",
   })
   if err != nil {
      return nil, err
   }

   resp, err := maya.Post(target, headers, reqBody)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()

   var activation AccountActivation
   if err := json.NewDecoder(resp.Body).Decode(&activation); err != nil {
      return nil, err
   }
   return &activation, nil
}
