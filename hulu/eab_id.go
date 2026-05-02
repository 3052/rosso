package hulu

import (
   "41.neocities.org/maya"
   "encoding/json"
   "errors"
   "net/url"
)

type DeepLink struct {
   EabId   string `json:"eab_id"`
   Message string
}

func (d *Device) DeepLink(id string) (*DeepLink, error) {
   resp, err := maya.Get(
      &url.URL{
         Scheme: "https",
         Host:   "discover.hulu.com",
         Path:   "/content/v5/deeplink/playback",
         RawQuery: url.Values{
            "id":        {id},
            "namespace": {"entity"},
         }.Encode(),
      },
      map[string]string{"authorization": "Bearer " + d.UserToken},
   )
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   var result DeepLink
   err = json.NewDecoder(resp.Body).Decode(&result)
   if err != nil {
      return nil, err
   }
   if result.Message != "" {
      return nil, errors.New(result.Message)
   }
   return &result, nil
}
