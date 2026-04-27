package rakuten

import (
   "io"
   "net/url"

   "41.neocities.org/maya"
)

type StreamingUuid string

func AcquireLicense(uuid StreamingUuid, challenge []byte) ([]byte, error) {
   endpoint := &url.URL{
      Scheme: "https",
      Host:   "prod-playready.rakuten.tv",
      Path:   "/v1/licensing/pr",
   }
   values := url.Values{}
   values.Set("uuid", string(uuid))
   endpoint.RawQuery = values.Encode()

   resp, err := maya.Post(endpoint, nil, challenge)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()

   return io.ReadAll(resp.Body)
}
