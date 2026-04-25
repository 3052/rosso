// FILE: widevine_license.go
package plex

import (
   "errors"
   "io"
   "net/url"

   "41.neocities.org/maya"
)

func AcquireWidevineLicense(media *VodMedia, anonymous *AnonymousUser, body []byte) ([]byte, error) {
   if len(media.Part) == 0 {
      return nil, errors.New("no media parts found")
   }
   if media.Part[0].License == "" {
      return nil, errors.New("no license path found")
   }

   endpoint := &url.URL{
      Scheme: "https",
      Host:   "vod.provider.plex.tv",
      Path:   media.Part[0].License,
   }

   query := url.Values{}
   query.Set("x-plex-drm", "widevine")
   query.Set("x-plex-token", anonymous.AuthToken)
   endpoint.RawQuery = query.Encode()

   resp, err := maya.Post(endpoint, nil, body)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()

   return io.ReadAll(resp.Body)
}
