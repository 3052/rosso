package rakuten

import (
   "io"
   "net/url"

   "41.neocities.org/maya"
)

func AcquireLicense(targetInfo *StreamInfo, challenge []byte) ([]byte, error) {
   query := make(url.Values)
   query.Set("uuid", targetInfo.Wrid)

   endpoint := &url.URL{
      Scheme:   "https",
      Host:     "prod-playready.rakuten.tv",
      Path:     "/v1/licensing/pr",
      RawQuery: query.Encode(),
   }

   resp, err := maya.Post(endpoint, nil, challenge)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()

   licenseBytes, err := io.ReadAll(resp.Body)
   if err != nil {
      return nil, err
   }

   return licenseBytes, nil
}
