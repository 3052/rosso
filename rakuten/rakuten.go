package rakuten

import (
   "41.neocities.org/maya"
   "encoding/json"
   "errors"
   "fmt"
   "io"
   "net/url"
   "strconv"
   "strings"
)

type SeasonDetails struct {
   Id       string    `json:"id"`
   Title    string    `json:"title"`
   Episodes []Episode `json:"episodes"`
}

type Episode struct {
   Id          string      `json:"id"`
   Title       string      `json:"title"`
   ViewOptions ViewOptions `json:"view_options"`
}

func FetchSeason(id string, rating *Classification, region *Market) (*SeasonDetails, error) {
   target := &url.URL{
      Scheme: "https",
      Host:   "gizmo.rakuten.tv",
      Path:   "/v3/seasons/" + id,
   }

   query := url.Values{}
   query.Set("classification_id", strconv.Itoa(rating.NumericalId))
   query.Set("device_identifier", "atvui40")
   query.Set("market_code", region.Code)
   target.RawQuery = query.Encode()

   resp, err := maya.Get(target, nil)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()

   var respWrapper struct {
      Data SeasonDetails `json:"data"`
   }

   if err := json.NewDecoder(resp.Body).Decode(&respWrapper); err != nil {
      return nil, err
   }

   return &respWrapper.Data, nil
}

type TvShow struct {
   Id      string   `json:"id"`
   Title   string   `json:"title"`
   Seasons []Season `json:"seasons"`
}

type Season struct {
   Id string `json:"id"`
}

func FetchTvShow(tvShowId string, rating *Classification, region *Market) (*TvShow, error) {
   target := &url.URL{
      Scheme: "https",
      Host:   "gizmo.rakuten.tv",
      Path:   "/v3/tv_shows/" + tvShowId,
   }

   query := url.Values{}
   query.Set("classification_id", strconv.Itoa(rating.NumericalId))
   query.Set("device_identifier", "atvui40")
   query.Set("market_code", region.Code)
   target.RawQuery = query.Encode()

   resp, err := maya.Get(target, nil)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()

   var respWrapper struct {
      Data TvShow `json:"data"`
   }

   if err := json.NewDecoder(resp.Body).Decode(&respWrapper); err != nil {
      return nil, err
   }

   return &respWrapper.Data, nil
}

func (s *StreamInfo) FetchLicense(challenge []byte) ([]byte, error) {
   target, err := url.Parse(s.LicenseUrl)
   if err != nil {
      return nil, err
   }
   resp, err := maya.Post(target, nil, challenge)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   return io.ReadAll(resp.Body)
}

func (m *Movie) String() string {
   return formatPlayableDetails(m.Id, m.Title, m.ViewOptions.Private.Streams)
}

func (e *Episode) String() string {
   return formatPlayableDetails(e.Id, e.Title, e.ViewOptions.Private.Streams)
}

func formatPlayableDetails(identifier string, title string, playbackStreams []Stream) string {
   seenLanguages := make(map[string]bool)
   var audioLanguages []string
   for _, currentStream := range playbackStreams {
      for _, audioLanguage := range currentStream.AudioLanguages {
         if !seenLanguages[audioLanguage.Id] {
            seenLanguages[audioLanguage.Id] = true
            audioLanguages = append(audioLanguages, audioLanguage.Id)
         }
      }
   }
   formattedAudio := strings.Join(audioLanguages, ", ")
   return fmt.Sprintf("%s (%s) - Audio: %s", title, identifier, formattedAudio)
}

func (t *TvShow) String() string {
   var data strings.Builder
   for i, season_data := range t.Seasons {
      if i >= 1 {
         data.WriteByte('\n')
      }
      data.WriteString("season id = ")
      data.WriteString(season_data.Id)
   }
   return data.String()
}

func (p *ParsedUrl) IsMovie() bool {
   return p.ContentType == "movies"
}

func (p *ParsedUrl) IsTvShow() bool {
   return p.ContentType == "tv_shows"
}

type ParsedUrl struct {
   MarketCode  string
   ContentType string
   ContentId   string
}

func ParseUrl(targetUrl string) (*ParsedUrl, error) {
   parsed, err := url.Parse(targetUrl)
   if err != nil {
      return nil, err
   }

   if !strings.HasSuffix(parsed.Host, "rakuten.tv") {
      return nil, errors.New("invalid host")
   }

   segments := strings.Split(strings.Trim(parsed.Path, "/"), "/")
   if len(segments) == 0 || segments[0] == "" {
      return nil, errors.New("missing market code in path")
   }

   info := &ParsedUrl{
      MarketCode: segments[0],
   }

   if len(segments) == 3 {
      info.ContentType = segments[1]
      info.ContentId = segments[2]
   } else if len(segments) >= 5 && segments[1] == "player" {
      info.ContentType = segments[2]
      info.ContentId = segments[4]
   }

   query := parsed.Query()
   if info.ContentType == "" {
      info.ContentType = query.Get("content_type")
   }

   if info.ContentId == "" {
      if id := query.Get("content_id"); id != "" {
         info.ContentId = id
      } else if id := query.Get("tv_show_id"); id != "" {
         info.ContentId = id
         if info.ContentType == "" {
            info.ContentType = "tv_shows"
         }
      } else if id := query.Get("movie_id"); id != "" {
         info.ContentId = id
         if info.ContentType == "" {
            info.ContentType = "movies"
         }
      }
   }

   if info.MarketCode == "" || info.ContentType == "" || info.ContentId == "" {
      return nil, errors.New("could not extract all required components from url")
   }

   return info, nil
}

func (s *StreamInfo) GetManifest() (*url.URL, error) {
   return url.Parse(s.Url)
}
