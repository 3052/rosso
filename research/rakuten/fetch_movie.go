package rakuten

import (
   "encoding/json"
   "net/url"
   "strconv"

   "41.neocities.org/maya"
)

type Movie struct {
   Id          string      `json:"id"`
   ViewOptions ViewOptions `json:"view_options"`
}

type ViewOptions struct {
   Private Private `json:"private"`
}

type Private struct {
   Streams []Stream `json:"streams"`
}

type Stream struct {
   AudioLanguages []Language `json:"audio_languages"`
}

func FetchMovie(movieId string, rating *Classification, region *Market) (*Movie, error) {
   target := &url.URL{
      Scheme: "https",
      Host:   "gizmo.rakuten.tv",
      Path:   "/v3/movies/" + movieId,
   }

   query := url.Values{}
   query.Set("classification_id", strconv.Itoa(rating.NumericalId))
   query.Set("device_identifier", "atvui40")
   query.Set("market_code", region.Code)
   target.RawQuery = query.Encode()

   resp, err := maya.Get(target, nil)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()

   var respWrapper struct {
      Data Movie `json:"data"`
   }

   if err := json.NewDecoder(resp.Body).Decode(&respWrapper); err != nil {
      return nil, err
   }

   return &respWrapper.Data, nil
}
