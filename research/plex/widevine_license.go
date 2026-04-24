package plex

import (
   "io"
   "net/url"

   "41.neocities.org/maya"
)

func AcquireWidevineLicense(vod *VodMedia, anonymous *AnonymousUser, body []byte) ([]byte, error) {
   endpoint := &url.URL{
      Scheme: "https",
      Host:   "vod.provider.plex.tv",
      Path:   "/library/parts/" + vod.Id + "/license",
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
