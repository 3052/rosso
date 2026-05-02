package hulu

import (
   "41.neocities.org/maya"
   "encoding/json"
   "errors"
   "net/url"
)

type Device struct {
   DeviceToken string `json:"device_token"`
   Message     string // 2026-05-02
   UserToken   string `json:"user_token"`
}

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

   // Check if the API returned an explicit error message
   if result.Message != "" {
      return nil, errors.New(result.Message)
   }

   // NEW: Check if eab_id is missing (which means it's not playable)
   if result.EabId == "" {
      return nil, errors.New("content is not playable: missing eab_id in response")
   }

   return &result, nil
}
