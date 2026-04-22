package plex

import (
   "io"
   urlpkg "net/url"

   "41.neocities.org/maya"
)

func GetLicense(license string, authToken string, body []byte) ([]byte, error) {
   endpoint := &urlpkg.URL{
      Scheme: "https",
      Host:   "vod.provider.plex.tv",
      Path:   license,
   }

   query := urlpkg.Values{}
   query.Set("x-plex-drm", "widevine")
   query.Set("x-plex-token", authToken)
   endpoint.RawQuery = query.Encode()

   resp, err := maya.Post(endpoint, nil, body)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()

   return io.ReadAll(resp.Body)
}
