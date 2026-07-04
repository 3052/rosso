package hboMax

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

const (
   disco_client = "!:!:beam:!"
   device_info  = "!/!(!/!;!/!;!/!)"
)

const Markets = "amer apac emea latam"

// doReq handles executing the HTTP request and logging the method/URL
func doReq(req *http.Request) (*http.Response, error) {
   log.Println(req.Method, req.URL)
   return http.DefaultClient.Do(req)
}

// APIError represents a single error object from the Max API
type APIError struct {
   Code   string `json:"code"`
   Detail string `json:"detail"`
}

// APIErrors represents a collection of API errors and implements the error interface
type APIErrors []APIError

func (e APIErrors) Error() string {
   var b strings.Builder
   for i, err := range e {
      if i > 0 {
         b.WriteString(", ")
      }
      b.WriteString(err.Code)
      b.WriteString(": ")
      b.WriteString(err.Detail)
   }
   return b.String()
}

type Cookie struct {
   Name  string
   Value string
}

func StRequest() (*Cookie, error) {
   req, err := http.NewRequest(http.MethodGet, "https://default.prd.api.hbomax.com/token?realm=bolt", nil)
   if err != nil {
      return nil, err
   }
   req.Header.Set("x-device-info", device_info)
   req.Header.Set("x-disco-client", disco_client)

   resp, err := doReq(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()

   for _, each := range resp.Cookies() {
      if each.Name == "st" {
         return &Cookie{Name: each.Name, Value: each.Value}, nil
      }
   }
   return nil, errors.New("named cookie not present")
}

func (*Cookie) CachePath() string {
   return "rosso/hboMax/Cookie"
}

func (c *Cookie) String() string {
   return fmt.Sprintf("%v=%v", c.Name, c.Value)
}

type Initiate struct {
   LinkingCode string
   TargetUrl   string
}

func InitiateRequest(st *Cookie, market string) (*Initiate, error) {
   endpoint := fmt.Sprintf("https://default.beam-%v.prd.api.discomax.com/authentication/linkDevice/initiate", market)
   req, err := http.NewRequest(http.MethodPost, endpoint, nil)
   if err != nil {
      return nil, err
   }
   req.Header.Set("cookie", st.String())
   req.Header.Set("x-device-info", device_info)

   resp, err := doReq(req)
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
   }
   err = json.NewDecoder(resp.Body).Decode(&result)
   if err != nil {
      return nil, err
   }
   return &result.Data.Attributes, nil
}

func (i *Initiate) String() string {
   var data strings.Builder
   data.WriteString("target URL: ")
   data.WriteString(i.TargetUrl)
   data.WriteString("\nlinking code: ")
   data.WriteString(i.LinkingCode)
   return data.String()
}

type Login struct {
   Token string
}

// you must
// /authentication/linkDevice/initiate
// first or this will always fail
func LoginRequest(st *Cookie) (*Login, error) {
   req, err := http.NewRequest(http.MethodPost, "https://default.prd.api.hbomax.com/authentication/linkDevice/login", nil)
   if err != nil {
      return nil, err
   }
   req.Header.Set("cookie", st.String())

   resp, err := doReq(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()

   var result struct {
      Data struct {
         Attributes Login
      }
   }
   err = json.NewDecoder(resp.Body).Decode(&result)
   if err != nil {
      return nil, err
   }
   return &result.Data.Attributes, nil
}

func (*Login) CachePath() string {
   return "rosso/hboMax/Login"
}

type Playback struct {
   Drm struct {
      Schemes struct {
         PlayReady *Scheme
         Widevine  *Scheme
      }
   }
   Errors   APIErrors `json:"errors"`
   Fallback struct {
      Manifest struct {
         Url string // _fallback.mpd:1080p, .mpd:4K
      }
   }
   Manifest struct {
      Url string // 1080p
   }
}

func PlayReadyRequest(token, editId string) (*Playback, error) {
   return playback_request(token, editId, "playready")
}

func WidevineRequest(token, editId string) (*Playback, error) {
   return playback_request(token, editId, "widevine")
}

func playback_request(token, edit_id, drm string) (*Playback, error) {
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
      http.MethodPost,
      "https://default.prd.api.hbomax.com/playback-orchestrator/any/playback-orchestrator/v1/playbackInfo",
      bytes.NewReader(body),
   )
   if err != nil {
      return nil, err
   }
   req.Header.Set("authorization", "Bearer "+token)
   req.Header.Set("content-type", "application/json")
   req.Header.Set("x-disco-params", "realm=bolt")
   req.Header.Set("x-disco-client", disco_client)
   req.Header.Set("x-device-info", device_info)

   resp, err := doReq(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()

   var result Playback
   err = json.NewDecoder(resp.Body).Decode(&result)
   if err != nil {
      return nil, err
   }
   if len(result.Errors) > 0 {
      return nil, result.Errors
   }
   return &result, nil
}

func (*Playback) CachePath() string {
   return "rosso/hboMax/Playback"
}

func (p *Playback) GetManifest() (*url.URL, error) {
   manifest, err := url.Parse(p.Fallback.Manifest.Url)
   if err != nil {
      return nil, err
   }
   manifest.Path = strings.Replace(manifest.Path, "_fallback", "", 1)
   return manifest, nil
}

// SL2000 max 1080p
// SL3000 max 2160p
func (p *Playback) PlayReadyRequest(body []byte) ([]byte, error) {
   req, err := http.NewRequest(http.MethodPost, p.Drm.Schemes.PlayReady.LicenseUrl, bytes.NewReader(body))
   if err != nil {
      return nil, err
   }
   req.Header.Set("content-type", "text/xml")

   resp, err := doReq(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()

   if resp.StatusCode != http.StatusOK {
      return nil, errors.New(resp.Status)
   }
   return io.ReadAll(resp.Body)
}

func (p *Playback) WidevineRequest(body []byte) ([]byte, error) {
   req, err := http.NewRequest(http.MethodPost, p.Drm.Schemes.Widevine.LicenseUrl, bytes.NewReader(body))
   if err != nil {
      return nil, err
   }
   req.Header.Set("content-type", "application/x-protobuf")

   resp, err := doReq(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()

   if resp.StatusCode != http.StatusOK {
      return nil, errors.New(resp.Status)
   }
   return io.ReadAll(resp.Body)
}

type Scheme struct {
   LicenseUrl string
}
