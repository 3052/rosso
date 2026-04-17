package plex

import (
   "fmt"
   "io"
   "net/http"
)

// GetVODMetadata retrieves media part and stream details using the media's ratingKey
func GetVODMetadata(ratingKey, plexToken string) ([]byte, error) {
   reqURL := fmt.Sprintf("https://vod.provider.plex.tv/library/metadata/%s", ratingKey)

   req, err := http.NewRequest("GET", reqURL, nil)
   if err != nil {
      return nil, err
   }

   req.Header.Set("Accept", "application/json")
   req.Header.Set("X-Plex-Token", plexToken)

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
