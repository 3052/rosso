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

type Start struct {
   Profile Profile `json:"profile"`
   Market  Market  `json:"market"`
}

type Profile struct {
   Classification Classification `json:"classification"`
   AudioLanguage  Language       `json:"audio_language"`
}

type Classification struct {
   NumericalId int `json:"numerical_id"`
}

type Language struct {
   Id string `json:"id"`
}

type Market struct {
   Code string `json:"code"`
}

type StartRequest struct {
   DeviceIdentifier string         `json:"device_identifier"`
   DeviceMetadata   DeviceMetadata `json:"device_metadata"`
}

type DeviceMetadata struct {
   AppVersion   string `json:"app_version"`
   Brand        string `json:"brand"`
   Model        string `json:"model"`
   Os           string `json:"os"`
   SerialNumber string `json:"serial_number"`
   Uid          string `json:"uid"`
   Year         int    `json:"year"`
}

func FetchStart(marketCode string) (*Start, error) {
   target := &url.URL{
      Scheme: "https",
      Host:   "gizmo.rakuten.tv",
      Path:   "/v3/me/start",
   }

   query := url.Values{}
   query.Set("market_code", marketCode)
   target.RawQuery = query.Encode()

   payload := StartRequest{
      DeviceIdentifier: "web",
      DeviceMetadata: DeviceMetadata{
         AppVersion:   "app_version",
         Brand:        "brand",
         Model:        "model",
         Os:           "os",
         SerialNumber: "serial_number",
         Uid:          "uid",
         Year:         0,
      },
   }

   body, err := json.Marshal(payload)
   if err != nil {
      return nil, err
   }

   headers := map[string]string{
      "content-type": "application/json",
   }

   resp, err := maya.Post(target, headers, body)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()

   var apiResp struct {
      Data Start `json:"data"`
   }

   if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
      return nil, err
   }

   return &apiResp.Data, nil
}

type Movie struct {
   Id          string      `json:"id"`
   Title       string      `json:"title"`
   ViewOptions ViewOptions `json:"view_options"`
}

type ViewOptions struct {
   Private Private `json:"private"`
}

type Private struct {
   Streams []Stream `json:"streams"`
}

type Stream struct {
   AudioLanguages []Language `json:"audio_languages"`
}

func FetchMovie(movieId string, userClassification Classification, targetMarket Market) (*Movie, error) {
   target := &url.URL{
      Scheme: "https",
      Host:   "gizmo.rakuten.tv",
      Path:   "/v3/movies/" + movieId,
   }

   query := url.Values{}
   query.Set("classification_id", strconv.Itoa(userClassification.NumericalId))
   query.Set("device_identifier", "atvui40")
   query.Set("market_code", targetMarket.Code)
   target.RawQuery = query.Encode()

   resp, err := maya.Get(target, nil)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()

   var apiResp struct {
      Data Movie `json:"data"`
   }

   if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
      return nil, err
   }

   return &apiResp.Data, nil
}

func (targetMovie *Movie) String() string {
   return formatPlayableDetails(targetMovie.Id, targetMovie.Title, targetMovie.ViewOptions.Private.Streams)
}

func formatPlayableDetails(identifier string, title string, playbackStreams []Stream) string {
   seenLanguages := make(map[string]bool)
   var availableLanguages []string
   for _, currentStream := range playbackStreams {
      for _, audioLanguage := range currentStream.AudioLanguages {
         if !seenLanguages[audioLanguage.Id] {
            seenLanguages[audioLanguage.Id] = true
            availableLanguages = append(availableLanguages, audioLanguage.Id)
         }
      }
   }
   formattedAudio := strings.Join(availableLanguages, ", ")
   return fmt.Sprintf("%s (%s) - Audio: %s", title, identifier, formattedAudio)
}

type Address struct {
   MarketCode  string
   ContentType string
   ContentId   string
}

func ParseAddress(targetUrl string) (*Address, error) {
   target, err := url.Parse(targetUrl)
   if err != nil {
      return nil, err
   }

   if !strings.HasSuffix(target.Host, "rakuten.tv") {
      return nil, errors.New("invalid host")
   }

   segments := strings.Split(strings.Trim(target.Path, "/"), "/")
   if len(segments) == 0 || segments[0] == "" {
      return nil, errors.New("missing market code in path")
   }

   parsed := &Address{
      MarketCode: segments[0],
   }

   if len(segments) == 3 {
      parsed.ContentType = segments[1]
      parsed.ContentId = segments[2]
   } else if len(segments) >= 5 && segments[1] == "player" {
      parsed.ContentType = segments[2]
      parsed.ContentId = segments[4]
   }

   query := target.Query()
   if parsed.ContentType == "" {
      parsed.ContentType = query.Get("content_type")
   }

   if parsed.ContentId == "" {
      if id := query.Get("content_id"); id != "" {
         parsed.ContentId = id
      } else if id := query.Get("tv_show_id"); id != "" {
         parsed.ContentId = id
         if parsed.ContentType == "" {
            parsed.ContentType = "tv_shows"
         }
      } else if id := query.Get("movie_id"); id != "" {
         parsed.ContentId = id
         if parsed.ContentType == "" {
            parsed.ContentType = "movies"
         }
      }
   }

   if parsed.MarketCode == "" || parsed.ContentType == "" || parsed.ContentId == "" {
      return nil, errors.New("could not extract all required components from url")
   }

   return parsed, nil
}

func (a *Address) IsMovie() bool {
   return a.ContentType == "movies"
}

func (a *Address) IsTvShow() bool {
   return a.ContentType == "tv_shows"
}

type Season struct {
   Id       string    `json:"id"`
   Title    string    `json:"title"`
   Episodes []Episode `json:"episodes"`
}

type Episode struct {
   Id          string      `json:"id"`
   Title       string      `json:"title"`
   ViewOptions ViewOptions `json:"view_options"`
}

func FetchSeason(seasonId string, userClassification Classification, targetMarket Market) (*Season, error) {
   target := &url.URL{
      Scheme: "https",
      Host:   "gizmo.rakuten.tv",
      Path:   "/v3/seasons/" + seasonId,
   }

   query := url.Values{}
   query.Set("classification_id", strconv.Itoa(userClassification.NumericalId))
   query.Set("device_identifier", "atvui40")
   query.Set("market_code", targetMarket.Code)
   target.RawQuery = query.Encode()

   resp, err := maya.Get(target, nil)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()

   var apiResp struct {
      Data Season `json:"data"`
   }

   if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
      return nil, err
   }

   return &apiResp.Data, nil
}

func (targetEpisode *Episode) String() string {
   return formatPlayableDetails(targetEpisode.Id, targetEpisode.Title, targetEpisode.ViewOptions.Private.Streams)
}

type TvShow struct {
   Id      string   `json:"id"`
   Title   string   `json:"title"`
   Seasons []Season `json:"seasons"`
}

func FetchTvShow(tvShowId string, userClassification Classification, targetMarket Market) (*TvShow, error) {
   target := &url.URL{
      Scheme: "https",
      Host:   "gizmo.rakuten.tv",
      Path:   "/v3/tv_shows/" + tvShowId,
   }

   query := url.Values{}
   query.Set("classification_id", strconv.Itoa(userClassification.NumericalId))
   query.Set("device_identifier", "atvui40")
   query.Set("market_code", targetMarket.Code)
   target.RawQuery = query.Encode()

   resp, err := maya.Get(target, nil)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()

   var apiResp struct {
      Data TvShow `json:"data"`
   }

   if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
      return nil, err
   }

   return &apiResp.Data, nil
}

func (show *TvShow) String() string {
   var data strings.Builder
   for i, currentSeason := range show.Seasons {
      if i >= 1 {
         data.WriteByte('\n')
      }
      data.WriteString("season id: ")
      data.WriteString(currentSeason.Id)
   }
   return data.String()
}

type StreamingRequest struct {
   AudioLanguage            string `json:"audio_language"`
   AudioQuality             string `json:"audio_quality"`
   ClassificationId         int    `json:"classification_id"`
   ContentId                string `json:"content_id"`
   ContentType              string `json:"content_type"`
   DeviceIdentifier         string `json:"device_identifier"`
   DeviceSerial             string `json:"device_serial"`
   DeviceStreamVideoQuality string `json:"device_stream_video_quality"`
   Player                   string `json:"player"`
   SubtitleLanguage         string `json:"subtitle_language"`
   VideoType                string `json:"video_type"`
}

func FetchMovieStreaming(contentId string, userClassification Classification, audioLanguage string) (*StreamInfo, error) {
   return fetchStreaming(contentId, "movies", userClassification, audioLanguage)
}

func FetchEpisodeStreaming(contentId string, userClassification Classification, audioLanguage string) (*StreamInfo, error) {
   return fetchStreaming(contentId, "episodes", userClassification, audioLanguage)
}

func fetchStreaming(contentId string, contentType string, userClassification Classification, audioLanguage string) (*StreamInfo, error) {
   target := &url.URL{
      Scheme: "https",
      Host:   "gizmo.rakuten.tv",
      Path:   "/v3/avod/streamings",
   }

   payload := StreamingRequest{
      AudioLanguage:            audioLanguage,
      AudioQuality:             "2.0",
      ClassificationId:         userClassification.NumericalId,
      ContentId:                contentId,
      ContentType:              contentType,
      DeviceIdentifier:         "atvui40",
      DeviceSerial:             "not implemented",
      DeviceStreamVideoQuality: "UHD",
      Player:                   "atvui40:DASH-CENC:PR",
      SubtitleLanguage:         "MIS",
      VideoType:                "stream",
   }

   body, err := json.Marshal(payload)
   if err != nil {
      return nil, err
   }

   headers := map[string]string{
      "content-type": "application/json",
   }

   resp, err := maya.Post(target, headers, body)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()

   var apiResp struct {
      Data struct {
         StreamInfos []StreamInfo `json:"stream_infos"`
      } `json:"data"`
   }

   if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
      return nil, err
   }

   for _, info := range apiResp.Data.StreamInfos {
      return &info, nil
   }

   return nil, errors.New("no stream infos found")
}

type Url struct {
   Url url.URL
}

func (u *Url) UnmarshalText(text []byte) error {
   return u.Url.UnmarshalBinary(text)
}

func (u *Url) MarshalText() ([]byte, error) {
   return u.Url.MarshalBinary()
}

type StreamInfo struct {
   LicenseUrl Url `json:"license_url"`
   Url        Url // MPD
}

func (s *StreamInfo) FetchLicense(challenge []byte) ([]byte, error) {
   resp, err := maya.Post(&s.LicenseUrl.Url, nil, challenge)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   return io.ReadAll(resp.Body)
}
