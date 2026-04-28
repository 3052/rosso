package rakuten

import (
   "encoding/json"
   "net/url"
   "strconv"

   "41.neocities.org/maya"
)

type Season struct {
   Id       string    `json:"id"`
   Title    string    `json:"title"`
   Episodes []Episode `json:"episodes"`
}

type Episode struct {
   Id          string      `json:"id"`
   Title       string      `json:"title"`
   ViewOptions ViewOptions `json:"view_options"`
}

func FetchSeason(seasonId string, userClassification Classification, targetMarket Market) (*Season, error) {
   target := &url.URL{
      Scheme: "https",
      Host:   "gizmo.rakuten.tv",
      Path:   "/v3/seasons/" + seasonId,
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
      Data Season `json:"data"`
   }

   if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
      return nil, err
   }

   return &apiResp.Data, nil
}

func (targetEpisode *Episode) String() string {
   return formatPlayableDetails(targetEpisode.Id, targetEpisode.Title, targetEpisode.ViewOptions.Private.Streams)
}
