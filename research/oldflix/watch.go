package oldflix

import (
   "encoding/json"
   "fmt"
   "net/http"
   "net/url"
   "strings"
)

type PlaylistItem struct {
   File string `json:"file"`
}

type WatchPlayResponse struct {
   Playlist []PlaylistItem `json:"playlist"`
   Status   int            `json:"status"`
}

// WatchPlay requests the final CDN-signed M3U8 stream URL
func (c *Client) WatchPlay(contentID, movieID, trackID string) (string, error) {
   data := url.Values{}
   data.Set("id", contentID)
   data.Set("m", movieID)
   data.Set("tk", trackID) // tk is the audio/language track id

   req, err := http.NewRequest("POST", BaseURL+"/api/watch/play", strings.NewReader(data.Encode()))
   if err != nil {
      return "", err
   }

   req.Header.Set("Authorization", "Bearer "+c.Token)
   req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
   req.Header.Set("User-Agent", "okhttp/4.12.0")

   resp, err := c.HTTPClient.Do(req)
   if err != nil {
      return "", err
   }
   defer resp.Body.Close()

   var watchResp WatchPlayResponse
   if err := json.NewDecoder(resp.Body).Decode(&watchResp); err != nil {
      return "", fmt.Errorf("failed to decode watch play response: %w", err)
   }

   if len(watchResp.Playlist) > 0 {
      return watchResp.Playlist[0].File, nil
   }

   return "", fmt.Errorf("no playlist found in response")
}
