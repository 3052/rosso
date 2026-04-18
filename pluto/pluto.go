package pluto

import (
   "errors"
   "fmt"
   "net/url"
   "strings"
)

// It assumes Vod and Stitched.Paths always have at least one entry
func (s *Series) GetMovieUrl() *url.URL {
   // Directly access the required path based on the data guarantees
   return build_stitcher(s.SessionToken, s.Vod[0].Stitched.Paths[0].Path)
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

type Series struct {
   SessionToken string
   Vod          []Vod
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

// Define constants for the hardcoded URL parts
const (
   stitcherScheme = "https"
   stitcherHost   = "cfd-v4-service-stitcher-dash-use1-1.prd.pluto.tv"
)

var (
   app_name         = "androidtv"
   drm_capabilities = "widevine:L1"
)

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
