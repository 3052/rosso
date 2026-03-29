package pluto

import (
   "bytes"
   "encoding/json"
   "errors"
   "fmt"
   "io"
   "net/http"
   "net/url"
   "strings"
)

func (v *Vod) String() string {
   data := &strings.Builder{}
   var lines bool
   for _, season := range v.Seasons {
      for _, episode := range season.Episodes {
         if lines {
            data.WriteString("\n\n")
         } else {
            lines = true
         }
         data.WriteString("season = ")
         fmt.Fprint(data, season.Number)
         data.WriteString("\nepisode = ")
         fmt.Fprint(data, episode.Number)
         data.WriteString("\nname = ")
         data.WriteString(episode.Name)
         data.WriteString("\nid = ")
         data.WriteString(episode.Id)
      }
   }
   return data.String()
}

// pluto.tv/on-demand/movies/64946365c5ae350013623630
// pluto.tv/on-demand/movies/disobedience-ca-2018-1-1
func FetchSeries(movieShow string) (*Series, error) {
   data := url.Values{}
   data.Set("appName", app_name)
   data.Set("appVersion", "9")
   data.Set("clientID", "9")
   data.Set("clientModelNumber", "9")
   data.Set("deviceMake", "9")
   data.Set("deviceModel", "9")
   data.Set("deviceVersion", "9")
   data.Set("drmCapabilities", drm_capabilities)
   if strings.Contains(movieShow, "-") {
      data.Set("episodeSlugs", movieShow)
   } else {
      data.Set("seriesIDs", movieShow)
   }
   req := http.Request{
      URL: &url.URL{
         Scheme:   "https",
         Host:     "boot.pluto.tv",
         Path:     "/v4/start",
         RawQuery: data.Encode(),
      },
      Header: http.Header{},
   }
   resp, err := http.DefaultClient.Do(&req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   var result Series
   err = json.NewDecoder(resp.Body).Decode(&result)
   if err != nil {
      return nil, err
   }
   if strings.Contains(movieShow, "-") {
      if result.Vod[0].Slug != movieShow {
         return nil, errors.New("slug mismatch")
      }
   } else if result.Vod[0].Id != movieShow {
      return nil, errors.New("id mismatch")
   }
   return &result, nil
}

type Vod struct {
   Id      string
   Seasons []struct {
      Episodes []struct {
         Id       string `json:"_id"`
         Name     string
         Number   int
         Stitched Stitched
      }
      Number int
   }
   Slug     string
   Stitched *Stitched
}

type Dash struct {
   Body []byte
   Url  *url.URL
}

// Define constants for the hardcoded URL parts
const (
   stitcherScheme = "https"
   stitcherHost   = "cfd-v4-service-stitcher-dash-use1-1.prd.pluto.tv"
)

var (
   app_name         = "androidtv"
   drm_capabilities = "widevine:L1"
)

func Widevine(data []byte) ([]byte, error) {
   resp, err := http.Post(
      "https://service-concierge.clusters.pluto.tv/v1/wv/alt",
      "application/x-protobuf", bytes.NewReader(data),
   )
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   return io.ReadAll(resp.Body)
}

func (s *Series) GetEpisodeUrl(episodeId string) (*url.URL, error) {
   // Iterate through all seasons and episodes to find the matching ID
   for _, season := range s.Vod[0].Seasons {
      for _, episode := range season.Episodes {
         if episode.Id == episodeId {
            // Directly access the path based on the data guarantees
            path := episode.Stitched.Paths[0].Path
            return s.buildStitcherUrl(path), nil
         }
      }
   }
   return nil, errors.New("episode not found")
}

func (s *Series) buildStitcherUrl(path string) *url.URL {
   stitcher := &url.URL{
      Host:   stitcherHost,
      Path:   "/v2" + path,
      Scheme: stitcherScheme,
   }
   values := url.Values{}
   values.Set("jwt", s.SessionToken)
   stitcher.RawQuery = values.Encode()
   return stitcher
}

type Stitched struct {
   Paths []struct {
      Path string
   }
}

type Series struct {
   SessionToken string
   Vod          []Vod
}

// It assumes Vod and Stitched.Paths always have at least one entry
func (s *Series) GetMovieUrl() *url.URL {
   // Directly access the required path based on the data guarantees
   path := s.Vod[0].Stitched.Paths[0].Path
   return s.buildStitcherUrl(path)
}
func FetchDash(urlData *url.URL) (*Dash, error) {
   req := http.Request{
      URL:    urlData,
      Header: http.Header{},
   }
   resp, err := http.DefaultClient.Do(&req)
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
