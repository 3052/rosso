package stan

import (
   "bytes"
   "encoding/json"
   "fmt"
   "log"
   "net/http"
   "net/url"
   "strconv"
   "strings"
)

const (
   StanName    = "Stan-AndroidTV"
   StanVersion = "4.32.1"
)

var BaseUrl = []string{
   "023-stan.akamaized.net",
   "666-stan.akamaized.net", // geo block
   "aws.stan.video",
   "gec.stan.video",
}

// deviceData returns the device payload that must be sent with the
// session request (matches Python _device_data()).
func deviceData() url.Values {
   return url.Values{
      "type":         {"console"},
      "screenSize":   {"3840x2160"},
      "stanName":     {StanName},
      "stanVersion":  {StanVersion},
      "manufacturer": {"NVIDIA"},
      "model":        {"SHIELD Android TV"},
      "os":           {"Android-9"},
      "videoCodecs":  {"h264,decode,dovi,h263,h265,hevc,mjpeg,mpeg2v,mp4,mpeg4,vc1,vp8,vp9"},
      "audioCodecs":  {"aac"},
      "drm":          {"widevine"},
      "hdcpVersion":  {"2.2"},
      "colorSpace":   {"hdr10"},
   }
}

func do(req *http.Request) (*http.Response, error) {
   log.Println(req.Method, req.URL)
   return http.DefaultClient.Do(req)
}

type APIError struct {
   Code string
}

type APIErrors []APIError

func (e APIErrors) Error() string {
   codes := make([]string, len(e))
   for i, err := range e {
      codes[i] = err.Code
   }
   return fmt.Sprintf("stan: %s", strings.Join(codes, ", "))
}

type ActivationCode struct {
   Code string
   Url  string
}

func FetchActivationCode() (*ActivationCode, error) {
   // FIX: send generate=true as query param, not form body (matches Python device_code())
   req, err := http.NewRequest(
      "POST", "https://api.stan.com.au/login/v1/activation-codes/",
      strings.NewReader(""),
   )
   if err != nil {
      return nil, err
   }
   req.URL.RawQuery = url.Values{
      "generate": {"true"},
   }.Encode()
   resp, err := do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   var code ActivationCode
   err = json.NewDecoder(resp.Body).Decode(&code)
   if err != nil {
      return nil, err
   }
   return &code, nil
}

func (*ActivationCode) CachePath() string {
   return "rosso/stan/ActivationCode"
}

func (a *ActivationCode) String() string {
   var data strings.Builder
   data.WriteString("Stan.\n")
   data.WriteString("Log in with code\n")
   data.WriteString("1. Visit stan.com.au/activate\n")
   data.WriteString("2. Enter the code:\n")
   data.WriteString(a.Code)
   return data.String()
}

func (a *ActivationCode) Token() (*WebToken, error) {
   req, err := http.NewRequest("GET", a.Url, nil)
   if err != nil {
      return nil, err
   }
   resp, err := do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()

   // FIX: check HTTP status code (matches Python device_login())
   if resp.StatusCode != http.StatusOK {
      return nil, fmt.Errorf("stan: activation code not yet authorized (HTTP %d)", resp.StatusCode)
   }

   var web WebToken
   err = json.NewDecoder(resp.Body).Decode(&web)
   if err != nil {
      return nil, err
   }
   return &web, nil
}

// AppSession is the full session response from /login/v1/sessions/app.
// FIX: expanded to match Python _oauth() response fields.
type AppSession struct {
   JwToken string
   Renew   int64
   Now     int64
   UserId  string
   Profile Profile
   Errors  APIErrors
}

func (*AppSession) CachePath() string {
   return "rosso/stan/AppSession"
}

func (a *AppSession) FetchMedia(id int) (*Media, error) {
   req, err := http.NewRequest(
      "GET", "https://api.stan.com.au/concurrency/v1/streams", nil,
   )
   if err != nil {
      return nil, err
   }
   req.Header.Set("x-forwarded-for", "1.128.0.0")
   req.URL.RawQuery = url.Values{
      "capabilities.drm": {"widevine"},
      "format":           {"hls,dash"},
      "jwToken":          {a.JwToken},
      "programId":        {strconv.Itoa(id)},
      "quality":          {"sd"}, // auto
   }.Encode()
   resp, err := do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()

   var result struct {
      Media  Media
      Errors APIErrors
   }
   if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
      return nil, err
   }

   if len(result.Errors) > 0 {
      return nil, result.Errors
   }

   return &result.Media, nil
}

type Media struct {
   Drm *struct {
      CustomData string
      KeyId      string
   }
   VideoUrl string
}

func (m *Media) BaseUrl(host string) (*url.URL, error) {
   video, err := url.Parse(m.VideoUrl)
   if err != nil {
      return nil, err
   }
   video.Host = host
   return video, nil
}

func (*Media) CachePath() string {
   return "rosso/stan/Media"
}

func (m *Media) License(data []byte) ([]byte, error) {
   req, err := http.NewRequest(
      // final slash is needed
      "POST", "https://lic.drmtoday.com/license-proxy-widevine/cenc/",
      bytes.NewReader(data),
   )
   if err != nil {
      return nil, err
   }
   req.Header.Set("dt-custom-data", m.Drm.CustomData)
   resp, err := do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   var result struct {
      License []byte
   }
   err = json.NewDecoder(resp.Body).Decode(&result)
   if err != nil {
      return nil, err
   }
   return result.License, nil
}

// Profile represents a Stan user profile.
type Profile struct {
   Id            string
   Name          string
   IconImage     ProfileIcon
   IsKidsProfile bool
}

// ProfileIcon represents the icon image for a profile.
type ProfileIcon struct {
   Url string
}

// WebToken is the intermediate token returned from the activation URL.
// It is used as the payload for the session request.
type WebToken struct {
   JwToken string
}

func (*WebToken) CachePath() string {
   return "rosso/stan/WebToken"
}

// FetchSession exchanges the WebToken for a full AppSession.
// FIX: uses /login/v1/sessions/app (not /sessions/mobile/app),
//
//   includes device data in the payload,
//   and handles error responses.
func (w *WebToken) FetchSession() (*AppSession, error) {
   // Build the form payload: jwToken + device data (matches Python _oauth())
   payload := deviceData()
   payload.Set("jwToken", w.JwToken)

   req, err := http.NewRequest(
      "POST", "https://api.stan.com.au/login/v1/sessions/app",
      strings.NewReader(payload.Encode()),
   )
   if err != nil {
      return nil, err
   }
   req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
   // FIX: match Python _oauth() headers
   req.Header.Del("Accept")
   req.Header.Del("Connection")

   resp, err := do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()

   result := &AppSession{}
   err = json.NewDecoder(resp.Body).Decode(result)
   if err != nil {
      return nil, err
   }

   // FIX: handle error responses (matches Python _oauth() error handling)
   if len(result.Errors) > 0 {
      return nil, result.Errors
   }

   return result, nil
}

// stan.go
