package itv

import (
   "bytes"
   "encoding/json"
   "errors"
   "fmt"
   "io"
   "net/http"
   "net/http/cookiejar"
   "net/url"
   "path"
   "strings"
   _ "embed"
)

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

func (m *MediaFile) FetchDash() (*Dash, error) {
   var err error
   http.DefaultClient.Jar, err = cookiejar.New(nil)
   if err != nil {
      return nil, err
   }
   resp, err := http.Get(strings.Replace(m.Href, "itvpnpctv", "itvpnpdotcom", 1))
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   body, err := io.ReadAll(resp.Body)
   if err != nil {
      return nil, err
   }
   return &Dash{Body: body, Url: resp.Request.URL}, nil
}

func (p *Playlist) Get1080() (*MediaFile, error) {
   for _, file := range p.Playlist.Video.MediaFiles {
      if file.Resolution == "1080" {
         return &file, nil
      }
   }
   return nil, errors.New("1080p media file not found")
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

type Playlist struct {
   Error    string
   Playlist struct {
      Video struct {
         MediaFiles []MediaFile
      }
   }
}

type Dash struct {
   Body []byte
   Url  *url.URL
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
