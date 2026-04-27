package rakuten

import (
   "io"
   "net/url"

   "41.neocities.org/maya"
)

func AcquireEpisodeLicense(info *EpisodeStreamInfo, challenge []byte) ([]byte, error) {
   link := &url.URL{
      Scheme: "https",
      Host:   "prod-playready.rakuten.tv",
      Path:   "/v1/licensing/pr",
   }
   values := url.Values{}
   values.Set("uuid", string(info.Wrid))
   link.RawQuery = values.Encode()

   resp, err := maya.Post(link, nil, challenge)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()

   return io.ReadAll(resp.Body)
}
