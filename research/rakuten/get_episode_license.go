package rakuten

import (
   "io"
   "net/url"

   "41.neocities.org/maya"
)

func GetEpisodeLicense(streamItem *EpisodeStreaming, body []byte) ([]byte, error) {
   query := url.Values{}
   query.Set("uuid", streamItem.Id)

   target := &url.URL{
      Scheme:   "https",
      Host:     "prod-playready.rakuten.tv",
      Path:     "/v1/licensing/pr",
      RawQuery: query.Encode(),
   }

   resp, err := maya.Post(target, nil, body)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()

   return io.ReadAll(resp.Body)
}
