package itv

import (
   "bytes"
   _ "embed"
   "encoding/json"
   "errors"
   "fmt"
   "io"
   "net/http"
   "net/url"
   "path"
   "strings"
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

//go:embed ProgrammePage.gql
var programme_page string

func ParseLegacyId(urlData string) string {
   // 1. Get the last part of the URL (e.g., "10a5356a0001B")
   base := path.Base(urlData)
   // 2. Split the string by the character 'a'
   parts := strings.Split(base, "a")
   // 3. Join them back together with '/'
   return strings.Join(parts, "/")
}

func (m *MediaFile) FetchKeyService(data []byte) ([]byte, error) {
   req, err := http.NewRequest("POST", m.KeyServiceUrl, bytes.NewReader(data))
   if err != nil {
      return nil, err
   }
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   return io.ReadAll(resp.Body)
}

type MediaFile struct {
   Href          string // MPD
   KeyServiceUrl string // DRM
   Resolution    string
}

type Playlist struct {
   Error    string
   Playlist struct {
      Video struct {
         MediaFiles []MediaFile
      }
   }
}

func (p *Playlist) Get1080() (*MediaFile, error) {
   for _, file := range p.Playlist.Video.MediaFiles {
      if file.Resolution == "1080" {
         return &file, nil
      }
   }
   return nil, errors.New("1080p media file not found")
}

type Title struct {
   LatestAvailableVersion struct {
      PlaylistUrl string
   }
   Series *struct {
      SeriesNumber int
   }
   EpisodeNumber int
   Title         string
}

func (t *Title) String() string {
   data := &strings.Builder{}
   if t.Series != nil {
      fmt.Fprintln(data, "series =", t.Series.SeriesNumber)
      fmt.Fprintln(data, "episode =", t.EpisodeNumber)
   }
   if t.Title != "" {
      fmt.Fprintln(data, "title =", t.Title)
   }
   fmt.Fprint(data, "playlist = ", t.LatestAvailableVersion.PlaylistUrl)
   return data.String()
}

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
