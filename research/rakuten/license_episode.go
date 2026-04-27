package rakuten

import (
   "io"
   "net/url"

   "41.neocities.org/maya"
)

func AcquireEpisodeLicense(uuid EpisodeLicenseUuid, challenge []byte) ([]byte, error) {
   endpoint := url.URL{
      Scheme: "https",
      Host:   "prod-playready.rakuten.tv",
      Path:   "/v1/licensing/pr",
   }

   query := url.Values{}
   query.Set("uuid", string(uuid))
   endpoint.RawQuery = query.Encode()

   resp, err := maya.Post(&endpoint, nil, challenge)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()

   return io.ReadAll(resp.Body)
}
