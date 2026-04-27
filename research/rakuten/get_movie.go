package rakuten

import (
   "encoding/json"
   "net/url"
   "strconv"

   "41.neocities.org/maya"
)

type CinemaMovie struct {
   Id    string `json:"id"`
   Title string `json:"title"`
}

func GetMovie(movieId string, marketCode string, classificationId int) (*CinemaMovie, error) {
   link := &url.URL{
      Scheme: "https",
      Host:   "gizmo.rakuten.tv",
      Path:   "/v3/movies/" + movieId,
   }
   values := url.Values{}
   values.Set("classification_id", strconv.Itoa(classificationId))
   values.Set("device_identifier", "atvui40")
   values.Set("market_code", marketCode)
   link.RawQuery = values.Encode()

   resp, err := maya.Get(link, nil)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()

   var respWrapper struct {
      Data CinemaMovie `json:"data"`
   }
   if err := json.NewDecoder(resp.Body).Decode(&respWrapper); err != nil {
      return nil, err
   }
   return &respWrapper.Data, nil
}
