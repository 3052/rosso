package itv

import (
   "41.neocities.org/maya"
   "encoding/json"
   "errors"
   "io"
   "net/http"
   "net/url"
   "strings"
)

func FetchTitles(legacyId string) ([]Title, error) {
   var data strings.Builder
   err := json.NewEncoder(&data).Encode(map[string]string{
      "brandLegacyId": legacyId,
   })
   if err != nil {
      return nil, err
   }
   req := http.Request{
      URL: &url.URL{
         Scheme: "https",
         Host:   "content-inventory.prd.oasvc.itv.com",
         Path:   "/discovery",
         RawQuery: url.Values{
            "query":     {programme_page},
            "variables": {data.String()},
         }.Encode(),
      },
      Header: http.Header{},
   }
   resp, err := http.DefaultClient.Do(&req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   var result struct {
      Data struct {
         Titles []Title
      }
   }
   err = json.NewDecoder(resp.Body).Decode(&result)
   if err != nil {
      return nil, err
   }
   return result.Data.Titles, nil
}

func (m *MediaFile) FetchKeyService(body []byte) ([]byte, error) {
   target, err := url.Parse(m.KeyServiceUrl)
   if err != nil {
      return nil, err
   }
   resp, err := maya.Post(target, nil, body)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   return io.ReadAll(resp.Body)
}

// fetchPlaylist is the common underlying function doing the heavy lifting
func fetchPlaylist(urlData, drmSystem, maxSupported string) (*Playlist, error) {
   body, err := json.Marshal(map[string]any{
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
   target, err := url.Parse(urlData)
   if err != nil {
      return nil, err
   }
   resp, err := maya.Post(
      target,
      map[string]string{
         "accept":     "application/vnd.itv.vod.playlist.v4+json",
         "user-agent": "!",
      },
      body,
   )
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
