package rakuten

import (
   "41.neocities.org/maya"
   _ "embed"
   "encoding/json"
   "net/url"
)

//go:embed start.json
var start_json []byte

func FetchProfile(marketCode string) (*Profile, error) {
   target := url.URL{
      Scheme:   "https",
      Host:     "gizmo.rakuten.tv",
      Path:     "/v3/me/start",
      RawQuery: url.Values{"market_code": {marketCode}}.Encode(),
   }

   header := map[string]string{"content-type": "application/json"}

   resp, err := maya.Post(&target, header, start_json)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   var result struct {
      Data struct {
         Profile Profile
      }
   }
   err = json.NewDecoder(resp.Body).Decode(&result)
   if err != nil {
      return nil, err
   }
   return &result.Data.Profile, nil
}

type Profile struct {
   Classification struct {
      NumericalId int `json:"numerical_id"`
   }
}
