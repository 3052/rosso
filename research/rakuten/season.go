package rakuten

import (
   "encoding/json"
   "net/url"

   "41.neocities.org/maya"
)

type SeasonResponse struct {
   Data SeasonDetail `json:"data"`
}

type SeasonDetail struct {
   Id       string          `json:"id"`
   Episodes []SeasonEpisode `json:"episodes"`
}

type SeasonEpisode struct {
   Id    string `json:"id"`
   Title string `json:"title"`
}

func GetSeasonDetail(tvShowSeason *TvShowSeason, sessionData *SessionData) (*SeasonDetail, error) {
   endpoint := &url.URL{
      Scheme: "https",
      Host:   "gizmo.rakuten.tv",
      Path:   "/v3/seasons/" + tvShowSeason.Id,
   }

   query := url.Values{}
   query.Set("classification_id", sessionData.Profile.Classification.Id)
   query.Set("device_identifier", "atvui40")
   query.Set("market_code", sessionData.Market.Code)
   endpoint.RawQuery = query.Encode()

   resp, err := maya.Get(endpoint, nil)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()

   var seasonResponse SeasonResponse
   if err := json.NewDecoder(resp.Body).Decode(&seasonResponse); err != nil {
      return nil, err
   }

   return &seasonResponse.Data, nil
}
