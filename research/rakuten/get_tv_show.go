package rakuten

import (
   "encoding/json"
   "net/url"
   "strconv"

   "41.neocities.org/maya"
)

type SeasonId string

type TvShow struct {
   Id      string     `json:"id"`
   Title   string     `json:"title"`
   Seasons []TvSeason `json:"seasons"`
}

type TvSeason struct {
   Id    SeasonId `json:"id"`
   Title string   `json:"title"`
}

func GetTvShow(showId string, marketCode string, classificationId int) (*TvShow, error) {
   link := &url.URL{
      Scheme: "https",
      Host:   "gizmo.rakuten.tv",
      Path:   "/v3/tv_shows/" + showId,
   }
   values := url.Values{}
   values.Set("classification_id", strconv.Itoa(classificationId))
   values.Set("device_identifier", "atvui40")
   values.Set("market_code", marketCode)
   link.RawQuery = values.Encode()

   resp, err := maya.Get(link, nil)
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
