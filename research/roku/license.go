package roku

import (
   "io"
   "net/url"

   "41.neocities.org/maya"
)

func FetchLicense(targetWidevine Widevine, body []byte) ([]byte, error) {
   endpoint, err := url.Parse(targetWidevine.LicenseServer)
   if err != nil {
      return nil, err
   }

   headers := map[string]string{
      "content-type": "application/x-protobuf",
   }

   resp, err := maya.Post(endpoint, headers, body)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()

   return io.ReadAll(resp.Body)
}
