package plex

import (
   "encoding/json"
   "fmt"
   "net/http"
)

type VODMetadataResponse struct {
   MediaContainer struct {
      Metadata []struct {
         Media []struct {
            ID       string `json:"id"`
            Protocol string `json:"protocol"`
         } `json:"Media"`
      } `json:"Metadata"`
   } `json:"MediaContainer"`
}

// GetVODMetadata retrieves media streams. You can iterate through response.MediaContainer.Metadata[0].Media
// to find the protocol="dash" ID to pass to the license server.
func GetVODMetadata(ratingKey, plexToken string) (*VODMetadataResponse, error) {
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

   var result VODMetadataResponse
   if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
      return nil, err
   }

   return &result, nil
}
