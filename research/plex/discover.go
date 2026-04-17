package plex

import (
   "fmt"
   "io"
   "net/http"
   "net/url"
)

// GetDiscoverMatches retrieves library matches for a given URL path (e.g., "/movie/vicky-cristina-barcelona")
func GetDiscoverMatches(movieURLPath, plexToken string) ([]byte, error) {
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

   return io.ReadAll(resp.Body)
}
