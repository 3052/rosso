package rakuten

import (
   "encoding/json"
   "net/url"
   "strconv"

   "41.neocities.org/maya"
)

type EpisodeContentId string

type SeasonEpisode struct {
   Id EpisodeContentId `json:"id"`
}

type Season struct {
   Episodes []SeasonEpisode `json:"episodes"`
}

func GetSeason(seasonId SeasonId, classId ClassificationId, market MarketCode) (*Season, error) {
   endpoint := url.URL{
      Scheme: "https",
      Host:   "gizmo.rakuten.tv",
      Path:   "/v3/seasons/" + string(seasonId),
   }

   query := url.Values{}
   query.Set("classification_id", strconv.Itoa(int(classId)))
   query.Set("device_identifier", "atvui40")
   query.Set("market_code", string(market))
   endpoint.RawQuery = query.Encode()

   resp, err := maya.Get(&endpoint, nil)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()

   var wrapper struct {
      Data Season `json:"data"`
   }
   decoder := json.NewDecoder(resp.Body)
   if err := decoder.Decode(&wrapper); err != nil {
      return nil, err
   }

   return &wrapper.Data, nil
}
