package plex

import (
   "net/url"
)

func BuildManifestUrl(part *VodPart, authToken string) *url.URL {
   endpoint := &url.URL{
      Scheme: "https",
      Host:   "vod.provider.plex.tv",
      Path:   part.Key,
   }

   query := url.Values{}
   query.Set("x-plex-token", authToken)
   endpoint.RawQuery = query.Encode()

   return endpoint
}
