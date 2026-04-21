package plex

import (
   "encoding/json"
   "net/url"

   "41.neocities.org/maya"
)

type Match struct {
   MediaContainer MatchContainer `json:"MediaContainer"`
}

type MatchContainer struct {
   Metadata []MatchMetadata `json:"Metadata"`
}

type MatchMetadata struct {
   RatingKey string `json:"ratingKey"`
}

func GetMatch(matchUrl string, authToken string) (*Match, error) {
   query := url.Values{}
   query.Set("url", matchUrl)
   query.Set("x-plex-token", authToken)

   targetUrl := &url.URL{
      Scheme:   "https",
      Host:     "discover.provider.plex.tv",
      Path:     "/library/metadata/matches",
      RawQuery: query.Encode(),
   }

   resp, err := maya.Get(targetUrl, nil)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()

   var match Match
   if err := json.NewDecoder(resp.Body).Decode(&match); err != nil {
      return nil, err
   }

   return &match, nil
}
