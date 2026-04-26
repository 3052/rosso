package rakuten

import (
   "encoding/json"
   "net/url"

   "41.neocities.org/maya"
)

type Movie struct {
   Id    string `json:"id"`
   Title string `json:"title"`
}

func GetMovie(movieId string) (*Movie, error) {
   query := url.Values{}
   query.Set("classification_id", "41")
   query.Set("device_identifier", "atvui40")
   query.Set("market_code", "ie")

   target := &url.URL{
      Scheme:   "https",
      Host:     "gizmo.rakuten.tv",
      Path:     "/v3/movies/" + movieId,
      RawQuery: query.Encode(),
   }

   resp, err := maya.Get(target, nil)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()

   var wrapper struct {
      Data *Movie `json:"data"`
   }
   if err := json.NewDecoder(resp.Body).Decode(&wrapper); err != nil {
      return nil, err
   }
   return wrapper.Data, nil
}
