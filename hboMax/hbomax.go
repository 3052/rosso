package hboMax

import (
   "bytes"
   "encoding/json"
   "errors"
   "io"
   "net/http"
   "net/url"
   "slices"
   "strconv"
   "strings"
)

type ShowKey struct {
   Type string
   Id   string
}

// https://hbomax.com/at/en/movies/austin-powers-international-man-of-mystery/a979fb8b-f713-4de3-a625-d16ad4d37448
// https://hbomax.com/movies/one-battle-after-another/bebe611d-8178-481a-a4f2-de743b5b135a
// https://hbomax.com/shows/white-lotus/14f9834d-bc23-41a8-ab61-5c8abdbea505
// https://play.hbomax.com/show/31cb4b84-951a-4daf-8925-746fcdcddcb8
func (s *ShowKey) Parse(rawLink string) error {
   link, err := url.Parse(rawLink)
   if err != nil {
      return err
   }
   segments := strings.Split(strings.TrimPrefix(link.Path, "/"), "/")
   count := len(segments)
   if count < 2 {
      return errors.New("invalid URL format: not enough path segments")
   }
   s.Id = segments[count-1]
   // 1. Check for standard catalog types (plural) at position -3
   // URL: .../movies/slug/id
   if count >= 3 {
      segmentType := segments[count-3]
      if segmentType == "movies" || segmentType == "shows" {
         s.Type = segmentType
         return nil
      }
   }
   // 2. Check for player types (singular) at position -2
   // URL: .../show/id
   segmentType := segments[count-2]
   if segmentType == "show" {
      s.Type = segmentType
      return nil
   }
   return errors.New("unrecognized content type")
}

func (l Login) Movie(show *ShowKey) (*Videos, error) {
   var req http.Request
   req.Header = http.Header{}
   req.Header.Set("authorization", "Bearer "+l.Data.Attributes.Token)
   req.URL = &url.URL{
      Scheme: "https",
      Host:   api_host,
      RawQuery: url.Values{
         "include":          {"default"},
         "page[items.size]": {"1"},
      }.Encode(),
      Path: join(
         "/cms/routes/", strings.TrimSuffix(show.Type, "s"), "/", show.Id,
      ),
   }
   resp, err := http.DefaultClient.Do(&req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   var result Videos
   err = json.NewDecoder(resp.Body).Decode(&result)
   if err != nil {
      return nil, err
   }
   if len(result.Errors) >= 1 {
      return nil, &result.Errors[0]
   }
   return &result, nil
}
func (l *Login) playback(edit_id, drm string) (*Playback, error) {
   data, err := json.Marshal(map[string]any{
      "editId":               edit_id,
      "consumptionType":      "streaming",
      "appBundle":            "",         // required
      "applicationSessionId": "",         // required
      "firstPlay":            false,      // required
      "gdpr":                 false,      // required
      "playbackSessionId":    "",         // required
      "userPreferences":      struct{}{}, // required
      "capabilities": map[string]any{
         "contentProtection": map[string]any{
            "contentDecryptionModules": []any{
               map[string]string{
                  "drmKeySystem": drm,
               },
            },
         },
         "manifests": map[string]any{
            "formats": map[string]any{
               "dash": struct{}{}, // required
            }, // required
         }, // required
      }, // required
      "deviceInfo": map[string]any{
         "player": map[string]any{
            "mediaEngine": map[string]string{
               "name":    "", // required
               "version": "", // required
            }, // required
            "playerView": map[string]int{
               "height": 0, // required
               "width":  0, // required
            }, // required
            "sdk": map[string]string{
               "name":    "", // required
               "version": "", // required
            }, // required
         }, // required
      }, // required
   })
   var req http.Request
   if err != nil {
      return nil, err
   }
   req.Body = io.NopCloser(bytes.NewReader(data))
   req.Header = http.Header{}
   req.Header.Set("authorization", "Bearer "+l.Data.Attributes.Token)
   req.Header.Set("content-type", "application/json")
   req.Method = "POST"
   req.URL = &url.URL{
      Scheme: "https",
      Host:   api_host,
      Path:   "/playback-orchestrator/any/playback-orchestrator/v1/playbackInfo",
   }
   resp, err := http.DefaultClient.Do(&req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   if resp.StatusCode == 504 {
      return nil, errors.New(resp.Status) // bail since no response body
   }
   var result Playback
   err = json.NewDecoder(resp.Body).Decode(&result)
   if err != nil {
      return nil, err
   }
   if len(result.Errors) >= 1 {
      return nil, &result.Errors[0]
   }
   return &result, nil
}

const (
   api_host     = "default.prd.api.hbomax.com"
   disco_client = "!:!:beam:!"
   device_info  = "!/!(!/!;!/!;!/!)"
)

var Markets = []string{
   "amer",
   "apac",
   "emea",
   "latam",
}

func join(data ...string) string {
   return strings.Join(data, "")
}

func (l *Login) PlayReady(editId string) (*Playback, error) {
   return l.playback(editId, "playready")
}

func (l Login) Season(show *ShowKey, number int) (*Videos, error) {
   var req http.Request
   req.Header = http.Header{}
   req.Header.Set("authorization", "Bearer "+l.Data.Attributes.Token)
   req.URL = &url.URL{
      Scheme: "https",
      Host:   api_host,
      Path:   "/cms/collections/generic-show-page-rail-episodes-tabbed-content",
      RawQuery: url.Values{
         "include":          {"default"},
         "pf[seasonNumber]": {strconv.Itoa(number)},
         "pf[show.id]":      {show.Id},
      }.Encode(),
   }
   resp, err := http.DefaultClient.Do(&req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   result := &Videos{}
   err = json.NewDecoder(resp.Body).Decode(result)
   if err != nil {
      return nil, err
   }
   return result, nil
}

func (p *Playback) PlayReady(data []byte) ([]byte, error) {
   resp, err := http.Post(
      p.Drm.Schemes.PlayReady.LicenseUrl, "text/xml",
      bytes.NewReader(data),
   )
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   if resp.StatusCode != http.StatusOK {
      return nil, errors.New(resp.Status)
   }
   return io.ReadAll(resp.Body)
}

// you must
// /authentication/linkDevice/initiate
// first or this will always fail
func (s St) Login() (*Login, error) {
   var req http.Request
   req.Header = http.Header{}
   req.AddCookie(s.Cookie)
   req.Method = "POST"
   req.URL = &url.URL{
      Scheme: "https",
      Host:   api_host, // Refactored
      Path:   "/authentication/linkDevice/login",
   }
   resp, err := http.DefaultClient.Do(&req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   result := &Login{}
   err = json.NewDecoder(resp.Body).Decode(result)
   if err != nil {
      return nil, err
   }
   return result, nil
}

func (s *St) Fetch() error {
   var req http.Request
   req.Header = http.Header{}
   req.Header.Set("x-device-info", device_info)
   req.Header.Set("x-disco-client", disco_client)
   req.URL = &url.URL{
      Scheme:   "https",
      Host:     api_host, // Refactored
      Path:     "/token",
      RawQuery: "realm=bolt",
   }
   resp, err := http.DefaultClient.Do(&req)
   if err != nil {
      return err
   }
   defer resp.Body.Close()
   for _, cookie := range resp.Cookies() {
      if cookie.Name == "st" {
         s.Cookie = cookie
         return nil
      }
   }
   return http.ErrNoCookie
}

// validVideoTypes acts as a set to hold the video types we want to keep.
var validVideoTypes = []string{
   "EPISODE",
   "MOVIE",
   "STANDALONE_EVENT",
}

func (v *Videos) FilterAndSort() {
   v.Included = slices.DeleteFunc(v.Included, func(vid *Video) bool {
      if vid.Attributes == nil {
         return true // Remove videos with nil attributes.
      }
      // We return 'true' to delete if the video's type is NOT in our slice.
      return !slices.Contains(validVideoTypes, vid.Attributes.VideoType)
   })
   slices.SortFunc(v.Included, func(a, b *Video) int {
      if a.Attributes == nil || b.Attributes == nil {
         return 0 // Consider them equal if attributes are missing.
      }
      return a.Attributes.EpisodeNumber - b.Attributes.EpisodeNumber
   })
}

type Video struct {
   Attributes *struct {
      SeasonNumber  int
      EpisodeNumber int
      Name          string
      VideoType     string
   }
   Relationships *struct {
      Edit *struct {
         Data struct {
            Id string
         }
      }
   }
}

type Videos struct {
   Errors   []Error
   Included []*Video
}

func (s St) Initiate(market string) (*Initiate, error) {
   var req http.Request
   req.Header = http.Header{}
   req.Header.Set("x-device-info", device_info)
   req.AddCookie(s.Cookie)
   req.Method = "POST"
   req.URL = &url.URL{
      Scheme: "https",
      Host:   join("default.beam-", market, ".prd.api.discomax.com"),
      Path:   "/authentication/linkDevice/initiate",
   }
   resp, err := http.DefaultClient.Do(&req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   if resp.StatusCode != http.StatusOK {
      return nil, errors.New(resp.Status)
   }
   var result struct {
      Data struct {
         Attributes Initiate
      }
      Errors []Error
   }
   err = json.NewDecoder(resp.Body).Decode(&result)
   if err != nil {
      return nil, err
   }
   if len(result.Errors) >= 1 {
      return nil, &result.Errors[0]
   }
   return &result.Data.Attributes, nil
}

type Login struct {
   Data struct {
      Attributes struct {
         Token string
      }
   }
}

func (e *Error) Error() string {
   if e.Detail != "" {
      return e.Detail
   }
   return e.Message
}

type Error struct {
   Detail  string // show was filtered by validator
   Message string // Token is missing or not valid
}

type Scheme struct {
   LicenseUrl string
}

type Playback struct {
   Drm struct {
      Schemes struct {
         PlayReady *Scheme
         Widevine  *Scheme
      }
   }
   Errors   []Error
   Fallback struct {
      Manifest struct {
         Url string // _fallback.mpd:1080p, .mpd:4K
      }
   }
   Manifest struct {
      Url string // 1080p
   }
}

type St struct {
   Cookie *http.Cookie
}

func (p *Playback) Widevine(data []byte) ([]byte, error) {
   resp, err := http.Post(
      p.Drm.Schemes.Widevine.LicenseUrl, "application/x-protobuf",
      bytes.NewReader(data),
   )
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   if resp.StatusCode != http.StatusOK {
      return nil, errors.New(resp.Status)
   }
   return io.ReadAll(resp.Body)
}

func (l *Login) Widevine(editId string) (*Playback, error) {
   return l.playback(editId, "widevine")
}

func (i *Initiate) String() string {
   var data strings.Builder
   data.WriteString("target URL = ")
   data.WriteString(i.TargetUrl)
   data.WriteString("\nlinking code = ")
   data.WriteString(i.LinkingCode)
   return data.String()
}

type Initiate struct {
   LinkingCode string
   TargetUrl   string
}

func (v *Video) String() string {
   var data strings.Builder
   if v.Attributes.SeasonNumber >= 1 {
      data.WriteString("season number = ")
      data.WriteString(strconv.Itoa(v.Attributes.SeasonNumber))
   }
   if v.Attributes.EpisodeNumber >= 1 {
      data.WriteString("\nepisode number = ")
      data.WriteString(strconv.Itoa(v.Attributes.EpisodeNumber))
   }
   if data.Len() >= 1 {
      data.WriteByte('\n')
   }
   data.WriteString("name = ")
   data.WriteString(v.Attributes.Name)
   data.WriteString("\nvideo type = ")
   data.WriteString(v.Attributes.VideoType)
   data.WriteString("\nedit id = ")
   data.WriteString(v.Relationships.Edit.Data.Id)
   return data.String()
}
func (p *Playback) Dash() (*Dash, error) {
   resp, err := http.Get(
      strings.Replace(p.Fallback.Manifest.Url, "_fallback", "", 1),
   )
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   var result Dash
   result.Body, err = io.ReadAll(resp.Body)
   if err != nil {
      return nil, err
   }
   result.Url = resp.Request.URL
   return &result, nil
}

type Dash struct {
   Body []byte
   Url  *url.URL
}
