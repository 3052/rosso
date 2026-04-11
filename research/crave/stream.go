// stream.go
package main

import (
   "encoding/json"
   "fmt"
   "io"
   "net/http"
)

// StreamResponse captures the final playback URL (MPD)
type StreamResponse struct {
   Playback  string `json:"playback"`
   Trickplay string `json:"trickplay"`
}

// GetStreamMeta retrieves the final JSON object containing the actual MPD URL
func GetStreamMeta(contentID, packageID, destinationID, token string) (*StreamResponse, error) {
   // Platform "48" and query params represent Web Player / Xbox HD configs
   url := fmt.Sprintf(
      "https://stream.video.9c9media.com/meta/content/%s/contentpackage/%s/destination/%s/platform/48?format=mpd&filter=ff&uhd=false&hd=true&mcv=false&mca=false&mta=true&stt=true",
      contentID, packageID, destinationID,
   )

   req, err := http.NewRequest("GET", url, nil)
   if err != nil {
      return nil, err
   }

   // Essential headers from HAR
   req.Header.Set("accept", "*/*")
   req.Header.Set("accept-language", "en-US,en;q=0.9")
   req.Header.Set("authorization", "Bearer "+token)
   req.Header.Set("origin", "https://www.crave.ca")
   req.Header.Set("referer", "https://www.crave.ca/")
   req.Header.Set("user-agent", "Xbox One")

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

   var streamResp StreamResponse
   if err := json.Unmarshal(body, &streamResp); err != nil {
      return nil, err
   }

   return &streamResp, nil
}
