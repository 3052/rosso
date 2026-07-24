package pluto

import (
   "bytes"
   "encoding/json"
   "errors"
   "fmt"
   "io"
   "log"
   "net/http"
   "net/url"
   "strings"
)

// Define constants for the hardcoded URL parts
const (
   stitcherScheme = "https"
   stitcherHost   = "cfd-v4-service-stitcher-dash-use1-1.prd.pluto.tv"
)

var (
   app_name         = "androidtv"
   drm_capabilities = "widevine:L1"
)

func FetchWidevine(body []byte) ([]byte, error) {
   target := &url.URL{
      Scheme: "https",
      Host:   "service-concierge.clusters.pluto.tv",
      Path:   "/v1/wv/alt",
   }
   req, err := http.NewRequest("POST", target.String(), bytes.NewReader(body))
   if err != nil {
      return nil, err
   }
   req.Header.Set("content-type", "application/x-protobuf")

   resp, err := do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   return io.ReadAll(resp.Body)
}

func build_stitcher(session_token, path string) *url.URL {
   stitcher := &url.URL{
      Host:   stitcherHost,
      Path:   "/v2" + path,
      Scheme: stitcherScheme,
   }
   values := url.Values{}
   values.Set("jwt", session_token)
   stitcher.RawQuery = values.Encode()
   return stitcher
}

func do(req *http.Request) (*http.Response, error) {
   log.Println(req.Method, req.URL)
   return http.DefaultClient.Do(req)
}

type Series struct {
   SessionToken string
   Vod          []Vod
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
   target := &url.URL{
      Scheme:   "https",
      Host:     "boot.pluto.tv",
      Path:     "/v4/start",
      RawQuery: data.Encode(),
   }
   req, err := http.NewRequest("GET", target.String(), nil)
   if err != nil {
      return nil, err
   }

   resp, err := do(req)
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

func (*Series) CachePath() string {
   return "rosso/pluto/Series"
}

func (s *Series) GetEpisodeUrl(episodeId string) (*url.URL, error) {
   // Iterate through all seasons and episodes to find the matching ID
   for _, season := range s.Vod[0].Seasons {
      for _, episode := range season.Episodes {
         if episode.Id == episodeId {
            // Directly access the path based on the data guarantees
            return build_stitcher(
               s.SessionToken, episode.Stitched.Paths[0].Path,
            ), nil
         }
      }
   }
   return nil, errors.New("episode not found")
}

// It assumes Vod and Stitched.Paths always have at least one entry
func (s *Series) GetMovieUrl() *url.URL {
   // Directly access the required path based on the data guarantees
   return build_stitcher(s.SessionToken, s.Vod[0].Stitched.Paths[0].Path)
}

type Stitched struct {
   Paths []struct {
      Path string
   }
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
         fmt.Fprintln(data, "season:", season.Number)
         fmt.Fprintln(data, "episode:", episode.Number)
         fmt.Fprintln(data, "name:", episode.Name)
         fmt.Fprint(data, "id: ", episode.Id)
      }
   }
   return data.String()
}
