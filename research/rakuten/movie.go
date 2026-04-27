package rakuten

import (
   "encoding/json"
   "net/url"
   "strconv"

   "41.neocities.org/maya"
)

type MovieId string
type MovieContentId string

type Movie struct {
   Id MovieContentId `json:"id"`
}

func GetMovie(movieId MovieId, classId ClassificationId, market MarketCode) (*Movie, error) {
   endpoint := url.URL{
      Scheme: "https",
      Host:   "gizmo.rakuten.tv",
      Path:   "/v3/movies/" + string(movieId),
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
      Data Movie `json:"data"`
   }
   decoder := json.NewDecoder(resp.Body)
   if err := decoder.Decode(&wrapper); err != nil {
      return nil, err
   }

   return &wrapper.Data, nil
}
