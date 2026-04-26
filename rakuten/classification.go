package rakuten

import (
   "41.neocities.org/maya"
   _ "embed"
   "encoding/json"
   "net/url"
)

//go:embed classification.json
var classification_json []byte

func (c *Content) FetchClassification() (*Classification, error) {
   target := url.URL{
      Scheme:   "https",
      Host:     "gizmo.rakuten.tv",
      Path:     "/v3/me/start",
      RawQuery: url.Values{"market_code": {c.MarketCode}}.Encode(),
   }

   header := map[string]string{"content-type": "application/json"}

   resp, err := maya.Post(&target, header, classification_json)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   var result struct {
      Data struct {
         Profile struct {
            Classification Classification
         }
      }
   }
   err = json.NewDecoder(resp.Body).Decode(&result)
   if err != nil {
      return nil, err
   }
   return &result.Data.Profile.Classification, nil
}

type Classification struct {
   NumericalId int `json:"numerical_id"`
}
