package roku

import (
   "encoding/json"
   "net/url"

   "41.neocities.org/maya"
)

type ActivationCode string

type AccountActivation struct {
   Code ActivationCode `json:"code"`
}

type ActivationRequest struct {
   Platform string `json:"platform"`
}

func CreateAccountActivation(token ContentToken) (*AccountActivation, error) {
   endpoint := &url.URL{
      Scheme: "https",
      Host:   "googletv.web.roku.com",
      Path:   "/api/v1/account/activation",
   }

   headers := map[string]string{
      "content-type":         "application/json",
      "user-agent":           "trc-googletv; production; 0",
      "x-roku-content-token": string(token),
   }

   payload := ActivationRequest{
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

   var activation AccountActivation
   if err := json.NewDecoder(resp.Body).Decode(&activation); err != nil {
      return nil, err
   }

   return &activation, nil
}
