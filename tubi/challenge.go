package tubi

import (
   "io"
   "net/url"

   "41.neocities.org/maya"
)

func PostChallenge(server LicenseServer, challenge []byte) ([]byte, error) {
   target, err := url.Parse(server.Url)
   if err != nil {
      return nil, err
   }

   headers := map[string]string{
      "content-type": "application/x-protobuf",
   }

   response, err := maya.Post(target, headers, challenge)
   if err != nil {
      return nil, err
   }
   defer response.Body.Close()

   return io.ReadAll(response.Body)
}
