// FILE: rakuten/fetch_season.go
package rakuten

import (
   "encoding/json"
   "net/url"
   "strconv"

   "41.neocities.org/maya"
)

type SeasonDetails struct {
   Id       string    `json:"id"`
   Title    string    `json:"title"`
   Episodes []Episode `json:"episodes"`
}

type Episode struct {
   Id          string      `json:"id"`
   Title       string      `json:"title"`
   ViewOptions ViewOptions `json:"view_options"`
}

func FetchSeason(id string, rating *Classification, region *Market) (*SeasonDetails, error) {
   target := &url.URL{
      Scheme: "https",
      Host:   "gizmo.rakuten.tv",
      Path:   "/v3/seasons/" + id,
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
      Data SeasonDetails `json:"data"`
   }

   if err := json.NewDecoder(resp.Body).Decode(&respWrapper); err != nil {
      return nil, err
   }

   return &respWrapper.Data, nil
}
