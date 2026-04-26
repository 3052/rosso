package rakuten

import (
   "encoding/json"
   "net/url"

   "41.neocities.org/maya"
)

type SeasonDetails struct {
   Id       string    `json:"id"`
   Title    string    `json:"title"`
   Episodes []Episode `json:"episodes"`
}

type Episode struct {
   Id    string `json:"id"`
   Title string `json:"title"`
}

func GetSeasonDetails(seasonItem *Season) (*SeasonDetails, error) {
   query := url.Values{}
   query.Set("classification_id", "23")
   query.Set("device_identifier", "atvui40")
   query.Set("market_code", "fr")

   target := &url.URL{
      Scheme:   "https",
      Host:     "gizmo.rakuten.tv",
      Path:     "/v3/seasons/" + seasonItem.Id,
      RawQuery: query.Encode(),
   }

   resp, err := maya.Get(target, nil)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()

   var wrapper struct {
      Data *SeasonDetails `json:"data"`
   }
   if err := json.NewDecoder(resp.Body).Decode(&wrapper); err != nil {
      return nil, err
   }
   return wrapper.Data, nil
}
