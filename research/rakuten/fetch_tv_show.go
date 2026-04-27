package rakuten

import (
   "encoding/json"
   "net/url"
   "strconv"

   "41.neocities.org/maya"
)

type TvShow struct {
   Id      string   `json:"id"`
   Seasons []Season `json:"seasons"`
}

type Season struct {
   Id string `json:"id"`
}

func FetchTvShow(tvShowId string, rating *Classification, region *Market) (*TvShow, error) {
   target := &url.URL{
      Scheme: "https",
      Host:   "gizmo.rakuten.tv",
      Path:   "/v3/tv_shows/" + tvShowId,
   }

   query := url.Values{}
   query.Set("classification_id", strconv.Itoa(rating.NumericalId))
   query.Set("device_identifier", "atvui40")
   query.Set("market_code", region.Code)
   target.RawQuery = query.Encode()

   resp, err := maya.Get(target, nil)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()

   var respWrapper struct {
      Data TvShow `json:"data"`
   }

   if err := json.NewDecoder(resp.Body).Decode(&respWrapper); err != nil {
      return nil, err
   }

   return &respWrapper.Data, nil
}
