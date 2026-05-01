package crave

import (
   "encoding/json"
   "net/url"

   "41.neocities.org/maya"
)

type Subscription struct {
   Type string `json:"type"`
}

type SubscriptionResponse struct {
   Subscriptions []Subscription `json:"subscriptions"`
}

func GetSubscriptions(activeSession *Session) (*SubscriptionResponse, error) {
   endpoint := url.URL{
      Scheme: "https",
      Host:   "account.bellmedia.ca",
      Path:   "/api/subscription/v5",
   }

   headers := map[string]string{
      "authorization": "Bearer " + activeSession.AccessToken,
   }

   resp, err := maya.Get(&endpoint, headers)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()

   var response SubscriptionResponse
   if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
      return nil, err
   }

   return &response, nil
}
