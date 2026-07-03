package itv

import (
   "bytes"
   _ "embed"
   "encoding/json"
   "errors"
   "fmt"
   "io"
   "log"
   "net/http"
   "net/url"
   "path"
   "strings"
)

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

// doRequest is the unified HTTP client function that logs the request and executes it
func doRequest(req *http.Request) (*http.Response, error) {
   log.Println(req.Method, req.URL)
   return http.DefaultClient.Do(req)
}

func graphql_compact(data string) string {
   return strings.Join(strings.Fields(data), " ")
}

type MediaFile struct {
   Href          string // MPD
   KeyServiceUrl string // DRM
   Resolution    string
}

func (*MediaFile) CachePath() string {
   return "rosso/itv/MediaFile"
}

func (m *MediaFile) FetchKeyService(body []byte) ([]byte, error) {
   req, err := http.NewRequest(http.MethodPost, m.KeyServiceUrl, bytes.NewReader(body))
   if err != nil {
      return nil, err
   }

   resp, err := doRequest(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   return io.ReadAll(resp.Body)
}

type Playlist struct {
   Error    string
   Playlist struct {
      Video struct {
         MediaFiles []MediaFile
      }
   }
}

// FetchPlayReady fetches a playlist with PlayReady DRM requirements
func FetchPlayReady(address string) (*Playlist, error) {
   return fetchPlaylist(address, "playready", "SL3000")
}

// FetchWidevine fetches a playlist with Widevine DRM requirements
func FetchWidevine(address string) (*Playlist, error) {
   return fetchPlaylist(address, "widevine", "L3")
}

// fetchPlaylist is the common underlying function doing the heavy lifting
func fetchPlaylist(address, drmSystem, maxSupported string) (*Playlist, error) {
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

   req, err := http.NewRequest(http.MethodPost, address, bytes.NewReader(body))
   if err != nil {
      return nil, err
   }
   req.Header.Set("Accept", "application/vnd.itv.vod.playlist.v4+json")
   req.Header.Set("User-Agent", "!")

   resp, err := doRequest(req)
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

func FetchTitles(legacyId string) ([]Title, error) {
   variables, err := json.Marshal(map[string]string{
      "brandLegacyId": legacyId,
   })
   if err != nil {
      return nil, err
   }

   endpoint := &url.URL{
      Scheme: "https",
      Host:   "content-inventory.prd.oasvc.itv.com",
      Path:   "/discovery",
      RawQuery: url.Values{
         "query":     {graphql_compact(programme_page)},
         "variables": {string(variables)},
      }.Encode(),
   }

   req, err := http.NewRequest(http.MethodGet, endpoint.String(), nil)
   if err != nil {
      return nil, err
   }

   resp, err := doRequest(req)
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

   if len(result.Data.Titles) == 0 {
      return nil, errors.New("no titles found")
   }

   return result.Data.Titles, nil
}

func (t *Title) String() string {
   data := &strings.Builder{}
   if t.Series != nil {
      fmt.Fprintln(data, "series:", t.Series.SeriesNumber)
      fmt.Fprintln(data, "episode:", t.EpisodeNumber)
   }
   if t.Title != "" {
      fmt.Fprintln(data, "title:", t.Title)
   }
   fmt.Fprint(data, "playlist: ", t.LatestAvailableVersion.PlaylistUrl)
   return data.String()
}
