package hboMax

import (
   "bytes"
   "encoding/json"
   "errors"
   "fmt"
   "io"
   "net/http"
   "net/url"
   "slices"
   "strconv"
   "strings"
)

func (l Login) Movie(showData *Show) (*Videos, error) {
   var req http.Request
   req.URL = &url.URL{
      Scheme: "https",
      Host:   api_host,
      Path: fmt.Sprintf(
         "/cms/routes/%v/%v",
         strings.TrimSuffix(showData.Category, "s"), showData.Id,
      ),
      RawQuery: url.Values{
         "include":          {"default"},
         "page[items.size]": {"1"},
      }.Encode(),
   }
   req.Header = http.Header{}
   req.Header.Set("authorization", "Bearer "+l.Data.Attributes.Token)
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

func (l Login) Season(showData *Show, number int) (*Videos, error) {
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
         "pf[show.id]":      {showData.Id},
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

func (l *Login) Widevine(editId string) (*Playback, error) {
   return l.playback(editId, "widevine")
}

// 1080p SL2000
// 1440p SL3000
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

func (p *Playback) Dash() (*Dash, error) {
   resp, err := http.Get(
      strings.Replace(p.Fallback.Manifest.Url, "_fallback", "", 1),
   )
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

// https://hbomax.com/at/en/movies/austin-powers-international-man-of-mystery/a979fb8b-f713-4de3-a625-d16ad4d37448
// https://hbomax.com/movies/one-battle-after-another/bebe611d-8178-481a-a4f2-de743b5b135a
// https://hbomax.com/shows/white-lotus/14f9834d-bc23-41a8-ab61-5c8abdbea505
// https://play.hbomax.com/movie/b7b66574-c6e3-4ed3-a266-6bc44180252e
// https://play.hbomax.com/show/31cb4b84-951a-4daf-8925-746fcdcddcb8
func ParseShow(inputUrl string) (*Show, error) {
   parsedUrl, err := url.Parse(inputUrl)
   if err != nil {
      return nil, err
   }
   path := strings.TrimPrefix(parsedUrl.Path, "/")
   segments := strings.Split(path, "/")
   count := len(segments)
   if count < 2 {
      return nil, errors.New("invalid url path")
   }
   // Create the instance
   showData := Show{
      Id: segments[count-1],
   }
   // Check immediate parent (e.g., /movie/id)
   if count >= 2 && isCategory(segments[count-2]) {
      showData.Category = segments[count-2]
      return &showData, nil
   }
   // Check grandparent (e.g., /movies/slug/id)
   if count >= 3 && isCategory(segments[count-3]) {
      showData.Category = segments[count-3]
      return &showData, nil
   }
   return nil, errors.New("category not found")
}

type Show struct {
   Category string
   Id       string
}

type Videos struct {
   Errors   []Error
   Included []*Video
}

func FetchSt() (*http.Cookie, error) {
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
      return nil, err
   }
   defer resp.Body.Close()
   for _, cookie := range resp.Cookies() {
      if cookie.Name == "st" {
         return cookie, nil
      }
   }
   return nil, http.ErrNoCookie
}

// you must
// /authentication/linkDevice/initiate
// first or this will always fail
func FetchLogin(st *http.Cookie) (*Login, error) {
   var req http.Request
   req.Method = "POST"
   req.URL = &url.URL{
      Scheme: "https",
      Host:   api_host, // Refactored
      Path:   "/authentication/linkDevice/login",
   }
   req.Header = http.Header{}
   req.AddCookie(st)
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

func FetchInitiate(st *http.Cookie, market string) (*Initiate, error) {
   var req http.Request
   req.Method = "POST"
   req.URL = &url.URL{
      Scheme: "https",
      Host: fmt.Sprintf("default.beam-%v.prd.api.discomax.com", market),
      Path:   "/authentication/linkDevice/initiate",
   }
   req.Header = http.Header{}
   req.Header.Set("x-device-info", device_info)
   req.AddCookie(st)
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

func (l *Login) PlayReady(editId string) (*Playback, error) {
   return l.playback(editId, "playready")
}

type Login struct {
   Data struct {
      Attributes struct {
         Token string
      }
   }
}

type Scheme struct {
   LicenseUrl string
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

// validVideoTypes acts as a set to hold the video types we want to keep.
var validVideoTypes = []string{
   "EPISODE",
   "MOVIE",
   "STANDALONE_EVENT",
}

func isCategory(segment string) bool {
   switch segment {
   case "movies", "shows", "movie", "show":
      return true
   default:
      return false
   }
}

type Dash struct {
   Body []byte
   Url  *url.URL
}

func (e *Error) Error() string {
   var data strings.Builder
   data.WriteString("code = ")
   data.WriteString(e.Code)
   if e.Detail != "" {
      data.WriteString("\ndetail = ")
      data.WriteString(e.Detail)
   } else {
      data.WriteString("\nmessage = ")
      data.WriteString(e.Message)
   }
   return data.String()
}

type Error struct {
   Code    string
   Detail  string // show was filtered by validator
   Message string // Token is missing or not valid
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
