package rakuten

import (
   "encoding/json"
   "net/url"
   "strconv"

   "41.neocities.org/maya"
)

type SeasonEpisode struct {
   Id string `json:"id"`
}

type TvSeason struct {
   Id       string          `json:"id"`
   Episodes []SeasonEpisode `json:"episodes"`
}

func FetchTvSeason(session *UserSession, showSeason *TvShowSeason) (*TvSeason, error) {
   endpoint := &url.URL{
      Scheme: "https",
      Host:   "gizmo.rakuten.tv",
      Path:   "/v3/seasons/" + showSeason.Id,
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
      Data TvSeason `json:"data"`
   }
   if err := json.NewDecoder(resp.Body).Decode(&wrapper); err != nil {
      return nil, err
   }

   return &wrapper.Data, nil
}
