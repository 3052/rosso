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

const (
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

func FetchSt() (*http.Cookie, error) {
   req := http.Request{
      URL: &url.URL{
         Scheme:   "https",
         Host:     "default.prd.api.hbomax.com", // Refactored
         Path:     "/token",
         RawQuery: "realm=bolt",
      },
      Header: http.Header{},
   }
   req.Header.Set("x-device-info", device_info)
   req.Header.Set("x-disco-client", disco_client)
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

func isCategory(segment string) bool {
   switch segment {
   case "movies", "shows", "movie", "show":
      return true
   default:
      return false
   }
}

// https://hbomax.com/at/en/movies/austin-powers-international-man-of-mystery/a979fb8b-f713-4de3-a625-d16ad4d37448
// https://hbomax.com/movies/one-battle-after-another/bebe611d-8178-481a-a4f2-de743b5b135a
// https://hbomax.com/shows/white-lotus/14f9834d-bc23-41a8-ab61-5c8abdbea505
// https://play.hbomax.com/movie/b7b66574-c6e3-4ed3-a266-6bc44180252e
// https://play.hbomax.com/show/31cb4b84-951a-4daf-8925-746fcdcddcb8
func ParseShowId(urlData string) (string, error) {
   url_parse, err := url.Parse(urlData)
   if err != nil {
      return "", err
   }
   path := strings.TrimPrefix(url_parse.Path, "/")
   segments := strings.Split(path, "/")
   count := len(segments)
   if count < 2 {
      return "", errors.New("invalid url path")
   }
   // Extract the ID which is always the last segment
   id := segments[count-1]
   // Check immediate parent (e.g., /movie/id)
   if count >= 2 && isCategory(segments[count-2]) {
      return id, nil
   }
   // Check grandparent (e.g., /movies/slug/id)
   if count >= 3 && isCategory(segments[count-3]) {
      return id, nil
   }
   return "", errors.New("category not found")
}

type Dash struct {
   Body []byte
   Url  *url.URL
}

type Error struct {
   Code    string
   Detail  string // show was filtered by validator
   Message string // Token is missing or not valid
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

func (i *Included) String() string {
   var data strings.Builder
   if i.Attributes.SeasonNumber >= 1 {
      data.WriteString("season number = ")
      data.WriteString(strconv.Itoa(i.Attributes.SeasonNumber))
   }
   if i.Attributes.EpisodeNumber >= 1 {
      data.WriteString("\nepisode number = ")
      data.WriteString(strconv.Itoa(i.Attributes.EpisodeNumber))
   }
   if data.Len() >= 1 {
      data.WriteByte('\n')
   }
   data.WriteString("name = ")
   data.WriteString(i.Attributes.Name)
   data.WriteString("\nvideo type = ")
   data.WriteString(i.Attributes.VideoType)
   data.WriteString("\nedit id = ")
   data.WriteString(i.Relationships.Edit.Data.Id)
   return data.String()
}

type Included struct {
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

func FetchInitiate(st *http.Cookie, market string) (*Initiate, error) {
   req := http.Request{
      Method: "POST",
      URL: &url.URL{
         Scheme: "https",
         Host:   fmt.Sprintf("default.beam-%v.prd.api.discomax.com", market),
         Path:   "/authentication/linkDevice/initiate",
      },
      Header: http.Header{},
   }
   req.AddCookie(st)
   req.Header.Set("x-device-info", device_info)
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

type Initiate struct {
   LinkingCode string
   TargetUrl   string
}

func (i *Initiate) String() string {
   var data strings.Builder
   data.WriteString("target URL = ")
   data.WriteString(i.TargetUrl)
   data.WriteString("\nlinking code = ")
   data.WriteString(i.LinkingCode)
   return data.String()
}

func (l *Login) PlayReady(editId string) (*Playback, error) {
   return l.playback(editId, "playready")
}

// you must
// /authentication/linkDevice/initiate
// first or this will always fail
func FetchLogin(st *http.Cookie) (*Login, error) {
   req := http.Request{
      Method: "POST",
      URL: &url.URL{
         Scheme: "https",
         Host:   "default.prd.api.hbomax.com", // Refactored
         Path:   "/authentication/linkDevice/login",
      },
      Header: http.Header{},
   }
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

func (l *Login) Widevine(editId string) (*Playback, error) {
   return l.playback(editId, "widevine")
}

func (l *Login) playback(edit_id, drm string) (*Playback, error) {
   body, err := json.Marshal(map[string]any{
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
   if err != nil {
      return nil, err
   }
   req, err := http.NewRequest(
      "POST", "https://default.prd.api.hbomax.com", bytes.NewReader(body),
   )
   if err != nil {
      return nil, err
   }
   req.URL.Path = "/playback-orchestrator/any/playback-orchestrator/v1/playbackInfo"
   req.Header.Set("authorization", "Bearer "+l.Data.Attributes.Token)
   req.Header.Set("content-type", "application/json")
   resp, err := http.DefaultClient.Do(req)
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

type Login struct {
   Data struct {
      Attributes struct {
         Token string
      }
   }
}

func (l Login) Movie(showId string) (*Page, error) {
   req := http.Request{
      URL: &url.URL{
         Scheme: "https",
         Host:   "default.prd.api.hbomax.com",
         Path:   "/cms/routes/movie/" + showId,
         RawQuery: url.Values{
            "include":          {"default"},
            "page[items.size]": {"1"},
         }.Encode(),
      },
      Header: http.Header{},
   }
   req.Header.Set("authorization", "Bearer "+l.Data.Attributes.Token)
   resp, err := http.DefaultClient.Do(&req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   var result Page
   err = json.NewDecoder(resp.Body).Decode(&result)
   if err != nil {
      return nil, err
   }
   if len(result.Errors) >= 1 {
      return nil, &result.Errors[0]
   }
   return &result, nil
}

func (l Login) Season(showId string, number int) (*Page, error) {
   req := http.Request{
      URL: &url.URL{
         Scheme: "https",
         Host:   "default.prd.api.hbomax.com",
         Path:   "/cms/collections/generic-show-page-rail-episodes-tabbed-content",
         RawQuery: url.Values{
            "include":          {"default"},
            "pf[seasonNumber]": {strconv.Itoa(number)},
            "pf[show.id]":      {showId},
         }.Encode(),
      },
      Header: http.Header{},
   }
   req.Header.Set("authorization", "Bearer "+l.Data.Attributes.Token)
   resp, err := http.DefaultClient.Do(&req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   result := &Page{}
   err = json.NewDecoder(resp.Body).Decode(result)
   if err != nil {
      return nil, err
   }
   return result, nil
}

type Page struct {
   Errors   []Error
   Included []*Included
}

func (p *Page) FilterAndSort() {
   p.Included = slices.DeleteFunc(p.Included, func(i *Included) bool {
      if i.Attributes == nil {
         return true // Remove videos with nil attributes.
      }
      // We return 'true' to delete if the video's type is NOT in our slice.
      return !slices.Contains(validVideoTypes, i.Attributes.VideoType)
   })
   slices.SortFunc(p.Included, func(a, b *Included) int {
      if a.Attributes == nil || b.Attributes == nil {
         return 0 // Consider them equal if attributes are missing.
      }
      return a.Attributes.EpisodeNumber - b.Attributes.EpisodeNumber
   })
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

// 1080p SL2000
// 1440p SL3000
func (p *Playback) PlayReady(body []byte) ([]byte, error) {
   resp, err := http.Post(
      p.Drm.Schemes.PlayReady.LicenseUrl, "text/xml",
      bytes.NewReader(body),
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

func (p *Playback) Widevine(body []byte) ([]byte, error) {
   resp, err := http.Post(
      p.Drm.Schemes.Widevine.LicenseUrl, "application/x-protobuf",
      bytes.NewReader(body),
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

type Scheme struct {
   LicenseUrl string
}
