package plex

import (
   "io"
   "net/url"

   "41.neocities.org/maya"
)

func PostLicense(authToken string, part *MediaPart, challenge []byte) ([]byte, error) {
   query := url.Values{}
   query.Set("x-plex-drm", "widevine")
   query.Set("x-plex-token", authToken)

   targetUrl := &url.URL{
      Scheme:   "https",
      Host:     "vod.provider.plex.tv",
      Path:     "/library/parts/" + part.Id + "/license",
      RawQuery: query.Encode(),
   }

   resp, err := maya.Post(targetUrl, nil, challenge)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()

   return io.ReadAll(resp.Body)
}
