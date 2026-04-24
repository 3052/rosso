package tubi

import (
   "io"
   "net/url"

   "41.neocities.org/maya"
)

func AcquireLicense(server *LicenseServer, body []byte) ([]byte, error) {
   target, err := url.Parse(server.Url)
   if err != nil {
      return nil, err
   }

   resp, err := maya.Post(target, nil, body)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()

   return io.ReadAll(resp.Body)
}
