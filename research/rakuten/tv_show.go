package rakuten

import (
   "encoding/json"
   "net/url"
   "strconv"

   "41.neocities.org/maya"
)

type TvShowSeason struct {
   Id          string `json:"id"`
   NumericalId int    `json:"numerical_id"`
   Number      int    `json:"number"`
}

type TvShow struct {
   Id      string         `json:"id"`
   Seasons []TvShowSeason `json:"seasons"`
}

func FetchTvShow(session *UserSession, tvShowId string) (*TvShow, error) {
   endpoint := &url.URL{
      Scheme: "https",
      Host:   "gizmo.rakuten.tv",
      Path:   "/v3/tv_shows/" + tvShowId,
   }
   values := url.Values{}
   values.Set("classification_id", strconv.Itoa(session.Profile.Classification.NumericalId))
   values.Set("device_identifier", "atvui40")
   values.Set("market_code", session.Market.Code)
   endpoint.RawQuery = values.Encode()

   resp, err := maya.Get(endpoint, nil)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()

   var wrapper struct {
      Data TvShow `json:"data"`
   }
   if err := json.NewDecoder(resp.Body).Decode(&wrapper); err != nil {
      return nil, err
   }

   return &wrapper.Data, nil
}
