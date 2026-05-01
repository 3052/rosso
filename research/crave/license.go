package crave

import (
   "encoding/json"
   "io"
   "net/url"
   "strconv"

   "41.neocities.org/maya"
)

type PlaybackContext struct {
   ContentId        int    `json:"contentId"`
   ContentpackageId int    `json:"contentpackageId"`
   DestinationId    int    `json:"destinationId"`
   Jwt              string `json:"jwt"`
   PlatformId       int    `json:"platformId"`
}

type LicenseRequest struct {
   Payload         string          `json:"payload"`
   PlaybackContext PlaybackContext `json:"playbackContext"`
}

func AcquireLicense(activeSession *Session, activeMedia *Media, available *AvailableContentPackage, challenge string) ([]byte, error) {
   endpoint := url.URL{
      Scheme: "https",
      Host:   "license.9c9media.com",
      Path:   "/playready",
   }

   contentIdInt, err := strconv.Atoi(activeMedia.FirstContent.Id)
   if err != nil {
      return nil, err
   }

   payload := LicenseRequest{
      Payload: challenge,
      PlaybackContext: PlaybackContext{
         ContentId:        contentIdInt,
         ContentpackageId: available.Id,
         DestinationId:    available.DestinationId,
         Jwt:              activeSession.AccessToken,
         PlatformId:       48,
      },
   }

   body, err := json.Marshal(payload)
   if err != nil {
      return nil, err
   }

   resp, err := maya.Post(&endpoint, nil, body)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()

   return io.ReadAll(resp.Body)
}
