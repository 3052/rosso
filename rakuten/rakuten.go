package rakuten

import (
   "encoding/json"
   "errors"
   "net/http"
   "net/url"
   "strconv"
   "strings"
)

func (s *StreamInfo) GetManifest() (*url.URL, error) {
   return url.Parse(s.Url)
}

type StreamInfo struct {
   LicenseUrl string `json:"license_url"`
   Url        string `json:"url"`
}

var classificationMap = map[string]int{
   "cz": 272,
   "es": 5,
   "fr": 23,
   "ie": 41,
   "nl": 69,
   "pl": 277,
   "pt": 64,
   "se": 282,
   "uk": 18,
}

// Parse extracts metadata from a Rakuten URL and returns a new Content struct
func ParseContent(urlData string) (*Content, error) {
   urlParse, err := url.Parse(urlData)
   if err != nil {
      return nil, err
   }

   c := &Content{}

   // Trim prefix once and extract the market code
   path := strings.TrimPrefix(urlParse.Path, "/")
   c.MarketCode, _, _ = strings.Cut(path, "/")

   // Check if the market code exists in the map and set ClassificationId
   var ok bool
   c.ClassificationId, ok = classificationMap[c.MarketCode]
   if !ok {
      return nil, errors.New("unknown market code")
   }

   // 1. Check Query Parameters
   query := urlParse.Query()
   contentType := query.Get("content_type")
   switch contentType {
   case "movies":
      c.Id = query.Get("content_id")
      c.Type = contentType
      return c, nil
   case "tv_shows":
      c.Id = query.Get("tv_show_id")
      c.Type = contentType
      return c, nil
   }

   // 2. Check Path Segments
   segments := strings.Split(path, "/")
   for _, segment := range segments {
      switch segment {
      case "movies", "tv_shows":
         c.Id = segments[len(segments)-1]
         c.Type = segment
         return c, nil
      }
   }

   return nil, errors.New("not a movie or tv show url")
}

// String implementation for MovieOrEpisode to pretty print details
func (m *MovieOrEpisode) String() string {
   seen := make(map[string]bool)
   var data strings.Builder
   data.WriteString("title = ")
   data.WriteString(m.Title)
   data.WriteString("\nid = ")
   data.WriteString(m.Id)
   for _, streamData := range m.ViewOptions.Private.Streams {
      for _, language := range streamData.AudioLanguages {
         if !seen[language.Id] {
            seen[language.Id] = true
            data.WriteString("\naudio language = ")
            data.WriteString(language.Id)
         }
      }
   }
   return data.String()
}

func (t TvShow) String() string {
   var data strings.Builder
   for i, season := range t.Seasons {
      if i >= 1 {
         data.WriteByte('\n')
      }
      data.WriteString("season id = ")
      data.WriteString(season.Id)
   }
   return data.String()
}

// Content represents the parsed Rakuten URL data
type Content struct {
   Id               string
   MarketCode       string
   Type             string
   ClassificationId int
}

// Constants for device and player configuration
const DeviceId = "atvui40"

const (
   PlayReady Player = DeviceId + ":DASH-CENC:PR"
   Widevine  Player = DeviceId + ":DASH-CENC:WVM"
)

const (
   Fhd VideoQuality = "FHD"
   Hd  VideoQuality = "HD"
)

type VideoQuality string

type Player string

type Season struct {
   Episodes []MovieOrEpisode `json:"episodes"`
}

type TvShow struct {
   Seasons []struct {
      Id string `json:"id"`
   } `json:"seasons"`
}

type MovieOrEpisode struct {
   Title       string `json:"title"`
   Id          string `json:"id"`
   ViewOptions struct {
      Private struct {
         Streams []struct {
            AudioLanguages []struct {
               Id string `json:"id"`
            } `json:"audio_languages"`
         } `json:"streams"`
      } `json:"private"`
   } `json:"view_options"`
}

func (c *Content) IsMovie() bool {
   return c.Type == "movies"
}

func (c *Content) IsTvShow() bool {
   return c.Type == "tv_shows"
}

func (c *Content) TvShow() (*TvShow, error) {
   urlData := url.URL{
      Scheme: "https",
      Host:   "gizmo.rakuten.tv",
      Path:   "/v3/tv_shows/" + c.Id,
      RawQuery: url.Values{
         "classification_id": {strconv.Itoa(c.ClassificationId)},
         "device_identifier": {DeviceId},
         "market_code":       {c.MarketCode},
      }.Encode(),
   }

   resp, err := http.Get(urlData.String())
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()

   if resp.StatusCode != http.StatusOK {
      return nil, errors.New(resp.Status)
   }

   var result struct {
      Data TvShow
   }
   if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
      return nil, err
   }
   return &result.Data, nil
}

func (c *Content) Movie() (*MovieOrEpisode, error) {
   urlData := url.URL{
      Scheme: "https",
      Host:   "gizmo.rakuten.tv",
      Path:   "/v3/movies/" + c.Id,
      RawQuery: url.Values{
         "classification_id": {strconv.Itoa(c.ClassificationId)},
         "device_identifier": {DeviceId},
         "market_code":       {c.MarketCode},
      }.Encode(),
   }

   resp, err := http.Get(urlData.String())
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()

   if resp.StatusCode != http.StatusOK {
      return nil, errors.New(resp.Status)
   }

   var result struct {
      Data MovieOrEpisode
   }
   if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
      return nil, err
   }
   return &result.Data, nil
}
