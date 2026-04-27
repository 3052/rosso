package rakuten

import (
   "encoding/json"
   "net/url"

   "41.neocities.org/maya"
)

type MovieResponse struct {
   Data Movie `json:"data"`
}

type Movie struct {
   Id    string `json:"id"`
   Title string `json:"title"`
}

func GetMovie(movieId string, sessionData *SessionData) (*Movie, error) {
   endpoint := &url.URL{
      Scheme: "https",
      Host:   "gizmo.rakuten.tv",
      Path:   "/v3/movies/" + movieId,
   }

   query := url.Values{}
   query.Set("classification_id", sessionData.Profile.Classification.Id)
   query.Set("device_identifier", "atvui40")
   query.Set("market_code", sessionData.Market.Code)
   endpoint.RawQuery = query.Encode()

   resp, err := maya.Get(endpoint, nil)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()

   var movieResponse MovieResponse
   if err := json.NewDecoder(resp.Body).Decode(&movieResponse); err != nil {
      return nil, err
   }

   return &movieResponse.Data, nil
}
