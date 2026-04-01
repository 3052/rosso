package crave

import (
   "encoding/json"
   "net/http"
   "net/url"
   "strings"
)

func (s *Subscription) String() string {
   var data strings.Builder
   data.WriteString("display name = ")
   data.WriteString(s.Experience.DisplayName)
   data.WriteString("\nexpiration date = ")
   data.WriteString(s.ExpirationDate)
   return data.String()
}

type Subscription struct {
   Experience struct {
      DisplayName string
   }
   ExpirationDate string
}

func (a *Account) FetchSubscriptions() ([]Subscription, error) {
   req := http.Request{
      URL: &url.URL{
         Scheme: "https",
         Host:   "account.bellmedia.ca",
         Path:   "/api/subscription/v5",
      },
      Header: http.Header{},
   }
   req.Header.Set("authorization", "Bearer "+a.AccessToken)
   resp, err := http.DefaultClient.Do(&req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   var result struct {
      Subscriptions []Subscription
   }
   err = json.NewDecoder(resp.Body).Decode(&result)
   if err != nil {
      return nil, err
   }
   return result.Subscriptions, nil
}
