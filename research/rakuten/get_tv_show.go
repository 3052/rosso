package rakuten

import (
   "encoding/json"
   "net/url"

   "41.neocities.org/maya"
)

type TvShow struct {
   Id      string   `json:"id"`
   Title   string   `json:"title"`
   Seasons []Season `json:"seasons"`
}

type Season struct {
   Id    string `json:"id"`
   Title string `json:"title"`
}

func GetTvShow(showId string) (*TvShow, error) {
   query := url.Values{}
   query.Set("classification_id", "23")
   query.Set("device_identifier", "atvui40")
   query.Set("market_code", "fr")

   target := &url.URL{
      Scheme:   "https",
      Host:     "gizmo.rakuten.tv",
      Path:     "/v3/tv_shows/" + showId,
      RawQuery: query.Encode(),
   }

   resp, err := maya.Get(target, nil)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()

   var wrapper struct {
      Data *TvShow `json:"data"`
   }
   if err := json.NewDecoder(resp.Body).Decode(&wrapper); err != nil {
      return nil, err
   }
   return wrapper.Data, nil
}
