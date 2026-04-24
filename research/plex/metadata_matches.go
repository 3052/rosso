package plex

import (
   "encoding/json"
   "net/url"

   "41.neocities.org/maya"
)

type MatchesResponse struct {
   MediaContainer MatchContainer `json:"MediaContainer"`
}

type MatchContainer struct {
   Metadata []MatchItem `json:"Metadata"`
}

type MatchItem struct {
   Guid      string `json:"guid"`
   RatingKey string `json:"ratingKey"`
   Title     string `json:"title"`
   Type      string `json:"type"`
}

func GetMetadataMatches(urlPath string, anonymous *AnonymousUser) (*MatchesResponse, error) {
   endpoint := &url.URL{
      Scheme: "https",
      Host:   "discover.provider.plex.tv",
      Path:   "/library/metadata/matches",
   }

   query := url.Values{}
   query.Set("url", urlPath)
   query.Set("x-plex-token", anonymous.AuthToken)
   endpoint.RawQuery = query.Encode()

   resp, err := maya.Get(endpoint, nil)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()

   var matches MatchesResponse
   if err := json.NewDecoder(resp.Body).Decode(&matches); err != nil {
      return nil, err
   }

   return &matches, nil
}
