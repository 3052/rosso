package crave

import (
   "encoding/json"
   "net/url"

   "41.neocities.org/maya"
)

type AvailableContentPackage struct {
   Id            int `json:"id"`
   DestinationId int `json:"destinationId"`
}

type Content struct {
   AvailableContentPackages []AvailableContentPackage `json:"availableContentPackages"`
}

func GetContent(activeSession *Session, activeMedia *Media) (*Content, error) {
   endpoint := url.URL{
      Scheme: "https",
      Host:   "playback.rte-api.bellmedia.ca",
      Path:   "/contents/" + activeMedia.FirstContent.Id,
   }

   headers := map[string]string{
      "x-client-platform":   "platform_jasper_web",
      "authorization":       "Bearer " + activeSession.AccessToken,
      "x-playback-language": "EN",
   }

   resp, err := maya.Get(&endpoint, headers)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()

   var targetContent Content
   if err := json.NewDecoder(resp.Body).Decode(&targetContent); err != nil {
      return nil, err
   }

   return &targetContent, nil
}
