package rakuten

import (
   "41.neocities.org/maya"
   "encoding/json"
   "errors"
   "io"
   "net/url"
   "strconv"
   "strings"
   _ "embed"
)

//go:embed classification.json
var classification_json []byte

func (c *Content) FetchClassification() (*Classification, error) {
   target := url.URL{
      Scheme:   "https",
      Host:     "gizmo.rakuten.tv",
      Path:     "/v3/me/start",
      RawQuery: url.Values{"market_code": {c.MarketCode}}.Encode(),
   }

   header := map[string]string{"content-type": "application/json"}

   resp, err := maya.Post(&target, header, classification_json)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   var result struct {
      Data struct {
         Profile struct {
            Classification Classification
         }
      }
   }
   err = json.NewDecoder(resp.Body).Decode(&result)
   if err != nil {
      return nil, err
   }
   return &result.Data.Profile.Classification, nil
}

type Classification struct {
   NumericalId int `json:"numerical_id"`
}
// For TV Shows, 'id' should be the Episode ID.
// For Movies, 'id' is ignored (uses c.Id).
func (c *Content) FetchStreamInfo(class *Classification, id, audioLanguage string, playerData Player, quality VideoQuality) (*StreamInfo, error) {
   value := map[string]any{
      "audio_language":              audioLanguage,
      "audio_quality":               "2.0",
      "classification_id":           class.NumericalId,
      "device_identifier":           DeviceId,
      "device_serial":               "not implemented",
      "device_stream_video_quality": string(quality),
      "player":                      string(playerData),
      "subtitle_language":           "MIS",
      "video_type":                  "stream",
   }
   switch c.Type {
   case "tv_shows":
      value["content_id"] = id
      value["content_type"] = "episodes"
   case "movies":
      value["content_id"] = c.Id
      value["content_type"] = "movies"
   }
   body, err := json.Marshal(value)
   if err != nil {
      return nil, err
   }
   resp, err := maya.Post(
      &url.URL{
         Scheme: "https",
         Host:   "gizmo.rakuten.tv",
         Path:   "/v3/avod/streamings",
      },
      map[string]string{"content-type": "application/json"},
      body,
   )
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   var result struct {
      Data struct {
         StreamInfos []StreamInfo `json:"stream_infos"`
      }
      Errors []struct {
         Message string // 2026-04-25
      }
   }
   err = json.NewDecoder(resp.Body).Decode(&result)
   if err != nil {
      return nil, err
   }
   if len(result.Errors) >= 1 {
      return nil, errors.New(result.Errors[0].Message)
   }
   return &result.Data.StreamInfos[0], nil
}

func (s *StreamInfo) GetManifest() (*url.URL, error) {
   return url.Parse(s.Url)
}

type StreamInfo struct {
   LicenseUrl string `json:"license_url"`
   Url        string `json:"url"`
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

// Constants for device and player configuration
const DeviceId = "atvui40"

const (
   PlayReady Player = DeviceId + ":DASH-CENC:PR"
   Widevine  Player = DeviceId + ":DASH-CENC:WVM"
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

// Content represents the parsed Rakuten URL data
type Content struct {
   Id         string
   MarketCode string
   Type       string
}

const (
   Fhd VideoQuality = "FHD"
   Hd  VideoQuality = "HD"
   Uhd VideoQuality = "UHD"
)

func (s *StreamInfo) FetchPlayReady(body []byte) ([]byte, error) {
   target, err := url.Parse(s.LicenseUrl)
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

func (s *StreamInfo) FetchWidevine(body []byte) ([]byte, error) {
   target, err := url.Parse(s.LicenseUrl)
   if err != nil {
      return nil, err
   }
   resp, err := maya.Post(
      target, map[string]string{"content-type": "application/x-protobuf"}, body,
   )
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   return io.ReadAll(resp.Body)
}

// Parse extracts metadata from a Rakuten URL and returns a new Content struct
func ParseContent(urlData string) (*Content, error) {
   url_parse, err := url.Parse(urlData)
   if err != nil {
      return nil, err
   }

   c := &Content{}

   // Trim prefix once and extract the market code
   path := strings.TrimPrefix(url_parse.Path, "/")
   c.MarketCode, _, _ = strings.Cut(path, "/")

   // 1. Check Query Parameters
   query := url_parse.Query()
   content_type := query.Get("content_type")
   switch content_type {
   case "movies":
      c.Id = query.Get("content_id")
      c.Type = content_type
      return c, nil
   case "tv_shows":
      c.Id = query.Get("tv_show_id")
      c.Type = content_type
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

func (c *Content) Movie(class *Classification) (*MovieOrEpisode, error) {
   resp, err := maya.Get(
      &url.URL{
         Scheme: "https",
         Host:   "gizmo.rakuten.tv",
         Path:   "/v3/movies/" + c.Id,
         RawQuery: url.Values{
            "classification_id": {strconv.Itoa(class.NumericalId)},
            "device_identifier": {DeviceId},
            "market_code":       {c.MarketCode},
         }.Encode(),
      },
      nil,
   )
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()

   if resp.StatusCode != 200 {
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

func (c *Content) TvShow(class *Classification) (*TvShow, error) {
   resp, err := maya.Get(
      &url.URL{
         Scheme: "https",
         Host:   "gizmo.rakuten.tv",
         Path:   "/v3/tv_shows/" + c.Id,
         RawQuery: url.Values{
            "classification_id": {strconv.Itoa(class.NumericalId)},
            "device_identifier": {DeviceId},
            "market_code":       {c.MarketCode},
         }.Encode(),
      },
      nil,
   )
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()

   if resp.StatusCode != 200 {
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

// Season fetches episodes for a specific season (GET).
func (c *Content) Season(class *Classification, seasonId string) (*Season, error) {
   resp, err := maya.Get(
      &url.URL{
         Scheme: "https",
         Host:   "gizmo.rakuten.tv",
         Path:   "/v3/seasons/" + seasonId,
         RawQuery: url.Values{
            "classification_id": {strconv.Itoa(class.NumericalId)},
            "device_identifier": {DeviceId},
            "market_code":       {c.MarketCode},
         }.Encode(),
      },
      nil,
   )
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   if resp.StatusCode != 200 {
      return nil, errors.New(resp.Status)
   }

   var result struct {
      Data Season
   }
   if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
      return nil, err
   }
   return &result.Data, nil
}
