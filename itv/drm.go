package itv

import (
   "bytes"
   "encoding/json"
   "errors"
   "net/http"
)

// FetchPlayReady fetches a playlist with PlayReady DRM requirements
func FetchPlayReady(urlData string) (*Playlist, error) {
   return fetchPlaylist(urlData, "playready", "SL3000")
}

// FetchWidevine fetches a playlist with Widevine DRM requirements
func FetchWidevine(urlData string) (*Playlist, error) {
   return fetchPlaylist(urlData, "widevine", "L3")
}

// fetchPlaylist is the common underlying function doing the heavy lifting
func fetchPlaylist(urlData, drmSystem, maxSupported string) (*Playlist, error) {
   data, err := json.Marshal(map[string]any{
      "client": map[string]string{
         "id": "browser",
      },
      "variantAvailability": map[string]any{
         "drm": map[string]string{
            "maxSupported": maxSupported,
            "system":       drmSystem,
         },
         "featureset": []string{ // need all these to get 720p
            "hd",
            "mpeg-dash",
            "single-track",
            drmSystem, // Injects "playready" or "widevine"
         },
         "platformTag": "ctv", // 1080p
      },
   })
   if err != nil {
      return nil, err
   }
   req, err := http.NewRequest("POST", urlData, bytes.NewReader(data))
   if err != nil {
      return nil, err
   }
   req.Header.Set("accept", "application/vnd.itv.vod.playlist.v4+json")
   req.Header.Set("user-agent", "!")
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   var result Playlist
   if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
      return nil, err
   }
   if result.Error != "" {
      return nil, errors.New(result.Error)
   }
   return &result, nil
}
