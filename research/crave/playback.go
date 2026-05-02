// FILE: crave/playback.go
package crave

import (
   "encoding/json"
   "net/url"
   "strconv"

   "41.neocities.org/maya"
)

type Playback struct {
   ContentId      int            `json:"contentId,string"`
   DestinationId  int            `json:"destinationId"`
   ContentPackage ContentPackage `json:"contentPackage"`
}

type ContentPackage struct {
   Id                int    `json:"id"`
   DurationInSeconds int    `json:"durationInSeconds"`
   Language          string `json:"language"`
   IsDescribedVideo  bool   `json:"isDescribedVideo"`
}

func GetPlayback(token *ProfileToken, activeMedia *Media) (*Playback, error) {
   endpoint := &url.URL{
      Scheme: "https",
      Host:   "playback.rte-api.bellmedia.ca",
      Path:   "/contents/" + strconv.Itoa(activeMedia.FirstContent.Id),
   }

   headers := map[string]string{
      "x-client-platform":   "platform_jasper_web",
      "authorization":       "Bearer " + token.AccessToken,
      "x-playback-language": "EN",
   }

   resp, err := maya.Get(endpoint, headers)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()

   activePlayback := &Playback{}
   if err := json.NewDecoder(resp.Body).Decode(activePlayback); err != nil {
      return nil, err
   }

   return activePlayback, nil
}
