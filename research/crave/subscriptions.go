package crave

import (
   "encoding/json"
   "net/url"

   "41.neocities.org/maya"
)

type Subscription struct {
   Type              string `json:"type"`
   SubscriptionState string `json:"subscriptionState"`
   StoreName         string `json:"storeName"`
   ExpirationDate    string `json:"expirationDate"`
   AutoRenewStatus   bool   `json:"autoRenewStatus"`
}

func GetSubscriptions(token *ProfileToken) ([]Subscription, error) {
   endpoint := &url.URL{
      Scheme: "https",
      Host:   "account.bellmedia.ca",
      Path:   "/api/subscription/v5",
   }

   headers := map[string]string{
      "authorization": "Bearer " + token.AccessToken,
   }

   resp, err := maya.Get(endpoint, headers)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()

   var wrapper struct {
      Subscriptions []Subscription `json:"subscriptions"`
   }
   if err := json.NewDecoder(resp.Body).Decode(&wrapper); err != nil {
      return nil, err
   }

   return wrapper.Subscriptions, nil
}
