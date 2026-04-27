package rakuten

import (
   "encoding/json"
   "net/url"
   "strconv"

   "41.neocities.org/maya"
)

type Episode struct {
   Id   ContentId   `json:"id"`
   Type ContentType `json:"type"`
}

type SeasonResponse struct {
   Id       ContentId   `json:"id"`
   Type     ContentType `json:"type"`
   Episodes []*Episode  `json:"episodes"`
}

func GetSeason(targetSeason *Season, sessionResp *SessionResponse) (*SeasonResponse, error) {
   query := make(url.Values)
   query.Set("classification_id", strconv.Itoa(sessionResp.Profile.Classification.NumericalId))
   query.Set("device_identifier", "atvui40")
   query.Set("market_code", sessionResp.Market.Code)

   endpoint := &url.URL{
      Scheme:   "https",
      Host:     "gizmo.rakuten.tv",
      Path:     "/v3/seasons/" + string(targetSeason.Id),
      RawQuery: query.Encode(),
   }

   resp, err := maya.Get(endpoint, nil)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()

   var wrapper struct {
      Data *SeasonResponse `json:"data"`
   }
   if err := json.NewDecoder(resp.Body).Decode(&wrapper); err != nil {
      return nil, err
   }

   return wrapper.Data, nil
}
