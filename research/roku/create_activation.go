package roku

import (
   "encoding/json"
   "net/url"

   "41.neocities.org/maya"
)

type Activation struct {
   Code string `json:"code"`
}

type ActivationPayload struct {
   Platform string `json:"platform"`
}

func CreateActivation(userToken ContentToken) (*Activation, error) {
   endpoint := &url.URL{
      Scheme: "https",
      Host:   "googletv.web.roku.com",
      Path:   "/api/v1/account/activation",
   }

   headers := map[string]string{
      "user-agent":           "trc-googletv; production; 0",
      "x-roku-content-token": string(userToken),
      "content-type":         "application/json",
   }

   payload := ActivationPayload{
      Platform: "googletv",
   }

   body, err := json.Marshal(payload)
   if err != nil {
      return nil, err
   }

   resp, err := maya.Post(endpoint, headers, body)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()

   var activation Activation
   if err := json.NewDecoder(resp.Body).Decode(&activation); err != nil {
      return nil, err
   }

   return &activation, nil
}
