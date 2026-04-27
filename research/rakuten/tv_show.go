package rakuten

import (
   "encoding/json"
   "net/url"

   "41.neocities.org/maya"
)

type TvShowResponse struct {
   Data TvShow `json:"data"`
}

type TvShow struct {
   Id      string         `json:"id"`
   Title   string         `json:"title"`
   Seasons []TvShowSeason `json:"seasons"`
}

type TvShowSeason struct {
   Id    string `json:"id"`
   Title string `json:"title"`
}

func GetTvShow(tvShowId string, sessionData *SessionData) (*TvShow, error) {
   endpoint := &url.URL{
      Scheme: "https",
      Host:   "gizmo.rakuten.tv",
      Path:   "/v3/tv_shows/" + tvShowId,
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

   var tvShowResponse TvShowResponse
   if err := json.NewDecoder(resp.Body).Decode(&tvShowResponse); err != nil {
      return nil, err
   }

   return &tvShowResponse.Data, nil
}
