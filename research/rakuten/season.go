package rakuten

import (
   "encoding/json"
   "net/url"

   "41.neocities.org/maya"
)

type SeasonDetails struct {
   Id          string        `json:"id"`
   NumericalId int           `json:"numerical_id"`
   Title       string        `json:"title"`
   Number      int           `json:"number"`
   Year        int           `json:"year"`
   ShortPlot   string        `json:"short_plot"`
   Plot        string        `json:"plot"`
   Episodes    []EpisodeItem `json:"episodes"`
}

type EpisodeItem struct {
   Id                string `json:"id"`
   NumericalId       int    `json:"numerical_id"`
   Title             string `json:"title"`
   ShortPlot         string `json:"short_plot"`
   DisplayName       string `json:"display_name"`
   Year              int    `json:"year"`
   Duration          int    `json:"duration"`
   Label             string `json:"label"`
   Number            int    `json:"number"`
   DurationInSeconds int    `json:"duration_in_seconds"`
   Plot              string `json:"plot"`
   SeasonId          string `json:"season_id"`
   SeasonNumber      int    `json:"season_number"`
   TvShowId          string `json:"tv_show_id"`
   TvShowTitle       string `json:"tv_show_title"`
}

func GetSeasonDetails(season *SeasonSummary) (*SeasonDetails, error) {
   location := &url.URL{
      Scheme: "https",
      Host:   "gizmo.rakuten.tv",
      Path:   "/v3/seasons/" + url.PathEscape(season.Id),
   }
   query := url.Values{}
   query.Set("classification_id", "18")
   query.Set("device_identifier", "atvui40")
   query.Set("market_code", "uk")
   location.RawQuery = query.Encode()

   resp, err := maya.Get(location, nil)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()

   var response struct {
      Data *SeasonDetails `json:"data"`
   }
   if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
      return nil, err
   }
   return response.Data, nil
}
