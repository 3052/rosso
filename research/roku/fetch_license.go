package roku

import (
   "io"
   "net/url"

   "41.neocities.org/maya"
)

func FetchLicense(widevineConfig Widevine, challenge []byte) ([]byte, error) {
   endpoint, err := url.Parse(widevineConfig.LicenseServer)
   if err != nil {
      return nil, err
   }

   headers := map[string]string{
      "content-type": "application/x-protobuf",
   }

   resp, err := maya.Post(endpoint, headers, challenge)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()

   return io.ReadAll(resp.Body)
}
