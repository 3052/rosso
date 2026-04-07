package itv

import (
   "bytes"
   _ "embed"
   "encoding/json"
   "errors"
   "io"
   "net/http"
   "net/http/cookiejar"
   "net/url"
   "path"
   "strconv"
   "strings"
)

func (p *Playlist) PlayReady(id string) error {
   data, err := json.Marshal(map[string]any{
      "client": map[string]string{
         "id": "browser",
      },
      "variantAvailability": map[string]any{
         "drm": map[string]string{
            "maxSupported": "SL3000",
            "system":       "playready",
         },
         "featureset": []string{
            "hd",
            "mpeg-dash",
            "single-track",
            "playready",
         },
         "platformTag": "ctv", // 1080p
      },
   })
   if err != nil {
      return err
   }
   req, err := http.NewRequest(
      "POST", "https://magni.itv.com/playlist/itvonline/ITV/"+id,
      bytes.NewReader(data),
   )
   if err != nil {
      return err
   }
   req.Header.Set("accept", "application/vnd.itv.vod.playlist.v4+json")
   req.Header.Set("user-agent", "!")
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return err
   }
   defer resp.Body.Close()
   return json.NewDecoder(resp.Body).Decode(p)
}

func Titles(legacyId string) ([]Title, error) {
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
            "query":     {graphql_compact(programme_page)},
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

///

func (m *MediaFile) Widevine(data []byte) ([]byte, error) {
   resp, err := http.Post(
      m.KeyServiceUrl, "application/x-protobuf", bytes.NewReader(data),
   )
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   return io.ReadAll(resp.Body)
}

func (t *Title) String() string {
   var data strings.Builder
   if t.Series != nil {
      data.WriteString("series = ")
      data.WriteString(strconv.Itoa(t.Series.SeriesNumber))
      data.WriteString("\nepisode = ")
      data.WriteString(strconv.Itoa(t.EpisodeNumber))
   }
   if t.Title != "" {
      if data.Len() >= 1 {
         data.WriteByte('\n')
      }
      data.WriteString("title = ")
      data.WriteString(t.Title)
   }
   if data.Len() >= 1 {
      data.WriteByte('\n')
   }
   data.WriteString("playlist = ")
   data.WriteString(t.LatestAvailableVersion.PlaylistUrl)
   return data.String()
}

func graphql_compact(data string) string {
   return strings.Join(strings.Fields(data), " ")
}

func FetchPlaylist(urlData string) (*Playlist, error) {
   data, err := json.Marshal(map[string]any{
      "client": map[string]string{
         "id": "browser",
      },
      "variantAvailability": map[string]any{
         "drm": map[string]string{
            "maxSupported": "L3",
            "system":       "widevine",
         },
         "featureset": []string{ // need all these to get 720p
            "hd",
            "mpeg-dash",
            "single-track",
            "widevine",
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
   err = json.NewDecoder(resp.Body).Decode(&result)
   if err != nil {
      return nil, err
   }
   if result.Error != "" {
      return nil, errors.New(result.Error)
   }
   return &result, nil
}

func (m *MediaFile) Dash() (*Dash, error) {
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

func (p *Playlist) FullHd() (*MediaFile, error) {
   for _, file := range p.Playlist.Video.MediaFiles {
      if file.Resolution == "1080" {
         return &file, nil
      }
   }
   return nil, errors.New("full hd (1080p) media file not found")
}

type Playlist struct {
   Error    string
   Playlist struct {
      Video struct {
         MediaFiles []MediaFile
      }
   }
}

type MediaFile struct {
   Href          string
   KeyServiceUrl string
   Resolution    string
}

type Dash struct {
   Body []byte
   Url  *url.URL
}
