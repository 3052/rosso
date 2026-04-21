// File: challenge.go
package tubi

import (
   "io"
   neturl "net/url"

   "41.neocities.org/maya"
)

func PostChallenge(url string, challengeData []byte) ([]byte, error) {
   targetUrl, err := neturl.Parse(url)
   if err != nil {
      return nil, err
   }

   headers := map[string]string{
      "content-type": "application/x-protobuf",
   }

   resp, err := maya.Post(targetUrl, headers, challengeData)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()

   return io.ReadAll(resp.Body)
}
