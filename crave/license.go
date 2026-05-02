// FILE: crave/license.go
package crave

import (
   "encoding/base64"
   "encoding/json"
   "io"
   "net/url"

   "41.neocities.org/maya"
)

func AcquireLicense(challenge []byte, token *ProfileToken, activePlayback *Playback) ([]byte, error) {
   endpoint := &url.URL{
      Scheme: "https",
      Host:   "license.9c9media.com",
      Path:   "/playready",
   }

   bodyMap := map[string]interface{}{
      "payload": base64.StdEncoding.EncodeToString(challenge),
      "playbackContext": map[string]interface{}{
         "contentId":        activePlayback.ContentId,
         "contentpackageId": activePlayback.ContentPackage.Id,
         "destinationId":    activePlayback.DestinationId,
         "jwt":              token.AccessToken,
         "platformId":       48,
      },
   }

   body, err := json.Marshal(bodyMap)
   if err != nil {
      return nil, err
   }

   resp, err := maya.Post(endpoint, nil, body)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()

   return io.ReadAll(resp.Body)
}
