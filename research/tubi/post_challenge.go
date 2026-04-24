package tubi

import (
   "io"
   "net/url"

   "41.neocities.org/maya"
)

func PostChallenge(server *LicenseServer, body []byte) ([]byte, error) {
   target, err := url.Parse(server.Url)
   if err != nil {
      return nil, err
   }

   headers := map[string]string{
      "content-type": "application/x-protobuf",
   }

   resp, err := maya.Post(target, headers, body)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()

   return io.ReadAll(resp.Body)
}
