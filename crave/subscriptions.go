// FILE: crave/subscriptions.go
package crave

import (
   "encoding/json"
   "net/url"

   "41.neocities.org/maya"
)

type LocalizedUrl struct {
   Fr string `json:"fr"`
   En string `json:"en"`
}

type ImageSet struct {
   Small LocalizedUrl `json:"small"`
}

type Circle struct {
   Svg ImageSet `json:"svg"`
   Png ImageSet `json:"png"`
}

type Logos struct {
   Circle Circle `json:"circle"`
}

type ContentPolicy struct {
   Sku string `json:"sku"`
}

type Experience struct {
   Id              string          `json:"id"`
   Sku             string          `json:"sku"`
   BrandId         string          `json:"brandId"`
   DisplayName     string          `json:"displayName"`
   Logos           Logos           `json:"logos"`
   ContentPolicies []ContentPolicy `json:"contentPolicies"`
}

type Subscription struct {
   Type              string     `json:"type"`
   Experience        Experience `json:"experience"`
   SubscriptionState string     `json:"subscriptionState"`
   StoreName         string     `json:"storeName"`
   ExpirationDate    string     `json:"expirationDate"`
   AutoRenewStatus   bool       `json:"autoRenewStatus"`
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
