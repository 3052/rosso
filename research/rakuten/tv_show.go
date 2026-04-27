package rakuten

import (
   "encoding/json"
   "net/url"
   "strconv"

   "41.neocities.org/maya"
)

type ContentId string
type ContentType string

type Season struct {
   Id   ContentId   `json:"id"`
   Type ContentType `json:"type"`
}

type TvShowResponse struct {
   Id      ContentId   `json:"id"`
   Type    ContentType `json:"type"`
   Seasons []*Season   `json:"seasons"`
}

func GetTvShow(showId string, sessionResp *SessionResponse) (*TvShowResponse, error) {
   query := make(url.Values)
   query.Set("classification_id", strconv.Itoa(sessionResp.Profile.Classification.NumericalId))
   query.Set("device_identifier", "atvui40")
   query.Set("market_code", sessionResp.Market.Code)

   endpoint := &url.URL{
      Scheme:   "https",
      Host:     "gizmo.rakuten.tv",
      Path:     "/v3/tv_shows/" + showId,
      RawQuery: query.Encode(),
   }

   resp, err := maya.Get(endpoint, nil)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()

   var wrapper struct {
      Data *TvShowResponse `json:"data"`
   }
   if err := json.NewDecoder(resp.Body).Decode(&wrapper); err != nil {
      return nil, err
   }

   return wrapper.Data, nil
}
