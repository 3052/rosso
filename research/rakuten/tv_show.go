package rakuten

import (
   "encoding/json"
   "net/url"

   "41.neocities.org/maya"
)

type TvShow struct {
   Id          string          `json:"id"`
   NumericalId int             `json:"numerical_id"`
   Title       string          `json:"title"`
   Plot        string          `json:"plot"`
   Countries   []Country       `json:"countries"`
   Genres      []Genre         `json:"genres"`
   Seasons     []SeasonSummary `json:"seasons"`
}

type Country struct {
   Id          string `json:"id"`
   NumericalId int    `json:"numerical_id"`
   Name        string `json:"name"`
}

type Genre struct {
   Id          string `json:"id"`
   NumericalId int    `json:"numerical_id"`
   Name        string `json:"name"`
   Adult       bool   `json:"adult"`
   Erotic      bool   `json:"erotic"`
   Kids        bool   `json:"kids"`
}

type SeasonSummary struct {
   Id               string `json:"id"`
   NumericalId      int    `json:"numerical_id"`
   Title            string `json:"title"`
   ShortPlot        string `json:"short_plot"`
   DisplayName      string `json:"display_name"`
   Year             int    `json:"year"`
   Label            string `json:"label"`
   Number           int    `json:"number"`
   TvShowId         string `json:"tv_show_id"`
   TvShowTitle      string `json:"tv_show_title"`
   NumberOfEpisodes int    `json:"number_of_episodes"`
}

func GetTvShow(showId string) (*TvShow, error) {
   location := &url.URL{
      Scheme: "https",
      Host:   "gizmo.rakuten.tv",
      Path:   "/v3/tv_shows/" + url.PathEscape(showId),
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
      Data *TvShow `json:"data"`
   }
   if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
      return nil, err
   }
   return response.Data, nil
}
