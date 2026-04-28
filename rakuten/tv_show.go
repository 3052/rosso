package rakuten

import (
   "encoding/json"
   "net/url"
   "strconv"
   "strings"

   "41.neocities.org/maya"
)

type TvShow struct {
   Id      string   `json:"id"`
   Title   string   `json:"title"`
   Seasons []Season `json:"seasons"`
}

func FetchTvShow(tvShowId string, userClassification Classification, targetMarket Market) (*TvShow, error) {
   target := &url.URL{
      Scheme: "https",
      Host:   "gizmo.rakuten.tv",
      Path:   "/v3/tv_shows/" + tvShowId,
   }

   query := url.Values{}
   query.Set("classification_id", strconv.Itoa(userClassification.NumericalId))
   query.Set("device_identifier", "atvui40")
   query.Set("market_code", targetMarket.Code)
   target.RawQuery = query.Encode()

   resp, err := maya.Get(target, nil)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()

   var apiResp struct {
      Data TvShow `json:"data"`
   }

   if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
      return nil, err
   }

   return &apiResp.Data, nil
}

func (show *TvShow) String() string {
   var builder strings.Builder
   for index, currentSeason := range show.Seasons {
      if index >= 1 {
         builder.WriteByte('\n')
      }
      builder.WriteString("season id = ")
      builder.WriteString(currentSeason.Id)
   }
   return builder.String()
}
