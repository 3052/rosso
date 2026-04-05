package oldflix

import (
   "encoding/json"
   "fmt"
   "net/http"
   "net/url"
   "strings"
)

func (b *Browse) Watch(trackId, token string) (string, error) {
   data := url.Values{}
   data.Set("id", b.Id)
   data.Set("m", b.Movie.Id)
   data.Set("tk", trackId) // tk is the audio/language track id
   req, err := http.NewRequest("POST", BaseUrl+"/api/watch/play", strings.NewReader(data.Encode()))
   if err != nil {
      return "", err
   }
   req.Header.Set("Authorization", "Bearer "+token)
   req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
   req.Header.Set("User-Agent", "okhttp/4.12.0")

   resp, err := http.DefaultClient.Do(req)
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

type PlaylistItem struct {
   File string `json:"file"`
}

type WatchPlayResponse struct {
   Playlist []PlaylistItem `json:"playlist"`
   Status   int            `json:"status"`
}
