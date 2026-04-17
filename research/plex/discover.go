package plex

import (
   "encoding/json"
   "fmt"
   "net/http"
   "net/url"
)

type DiscoverMatchesResponse struct {
   MediaContainer struct {
      Metadata []struct {
         RatingKey string `json:"ratingKey"`
         Title     string `json:"title"`
         Type      string `json:"type"`
         Guid      string `json:"guid"`
      } `json:"Metadata"`
   } `json:"MediaContainer"`
}

// GetDiscoverMatches returns the parsed metadata including the critical ratingKey
// needed to fetch the VOD playback details.
func GetDiscoverMatches(movieURLPath, plexToken string) (*DiscoverMatchesResponse, error) {
   baseURL, _ := url.Parse("https://discover.provider.plex.tv/library/metadata/matches")

   q := baseURL.Query()
   q.Set("url", movieURLPath)
   q.Set("x-plex-token", plexToken)
   baseURL.RawQuery = q.Encode()

   req, err := http.NewRequest("GET", baseURL.String(), nil)
   if err != nil {
      return nil, err
   }

   req.Header.Set("Accept", "application/json")

   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()

   if resp.StatusCode < 200 || resp.StatusCode >= 300 {
      return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
   }

   var result DiscoverMatchesResponse
   if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
      return nil, err
   }

   return &result, nil
}
