// File: license.go
package tubi

import (
   "io"
   "net/url"

   "41.neocities.org/maya"
)

func PostLicense(license_server *LicenseServer, body []byte) ([]byte, error) {
   target, err := url.Parse(license_server.Url)
   if err != nil {
      return nil, err
   }

   headers := map[string]string{
      "content-type": "application/x-protobuf",
   }

   response, err := maya.Post(target, headers, body)
   if err != nil {
      return nil, err
   }
   defer response.Body.Close()

   payload, err := io.ReadAll(response.Body)
   if err != nil {
      return nil, err
   }

   return payload, nil
}
