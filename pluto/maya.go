package pluto

import (
   "41.neocities.org/maya"
   "io"
   "net/url"
)

func FetchWidevine(body []byte) ([]byte, error) {
   resp, err := maya.Post(
      &url.URL{
         Scheme: "https",
         Host:   "service-concierge.clusters.pluto.tv",
         Path:   "/v1/wv/alt",
      },
      map[string]string{"content-type": "application/x-protobuf"},
      body,
   )
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   return io.ReadAll(resp.Body)
}
