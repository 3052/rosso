package itv

import (
   _ "embed"
   "errors"
   "fmt"
   "net/url"
   "path"
   "strings"
)

func (m *MediaFile) GetManifest() (*url.URL, error) {
   return url.Parse(strings.Replace(m.Href, "itvpnpctv", "itvpnpdotcom", 1))
}

// FetchPlayReady fetches a playlist with PlayReady DRM requirements
func FetchPlayReady(urlData string) (*Playlist, error) {
   return fetchPlaylist(urlData, "playready", "SL3000")
}

// FetchWidevine fetches a playlist with Widevine DRM requirements
func FetchWidevine(urlData string) (*Playlist, error) {
   return fetchPlaylist(urlData, "widevine", "L3")
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
