// playback.go
package main

import (
   "encoding/json"
   "fmt"
   "io"
   "net/http"
)

// GetPlaybackInfo fetches content metadata to grab the Package and Destination IDs
func GetPlaybackInfo(contentID, token string) (*PlaybackResponse, error) {
   url := fmt.Sprintf("https://playback.rte-api.bellmedia.ca/contents/%s", contentID)
   req, err := http.NewRequest("GET", url, nil)
   if err != nil {
      return nil, err
   }
   // Essential headers from HAR
   req.Header.Set("x-playback-language", "EN")
   req.Header.Set("x-client-platform", "platform_jasper_html")
   req.Header.Set("authorization", "Bearer "+token)
   client := &http.Client{}
   resp, err := client.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()

   if resp.StatusCode != http.StatusOK {
      return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, resp.Status)
   }

   body, err := io.ReadAll(resp.Body)
   if err != nil {
      return nil, err
   }

   var playResp PlaybackResponse
   if err := json.Unmarshal(body, &playResp); err != nil {
      return nil, err
   }

   return &playResp, nil
}
// PlaybackResponse maps only the necessary fields from the playback response JSON
type PlaybackResponse struct {
   AvailableContentPackages []struct {
      ID            int `json:"id"`
      DestinationID int `json:"destinationId"`
   } `json:"availableContentPackages"`
}
