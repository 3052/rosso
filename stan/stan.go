package stan

import (
   "bytes"
   "encoding/json"
   "fmt"
   "io"
   "net/http"
   "net/url"
   "strings"
)

// license.go
// quality.go
const (
   StanName    = "Stan-AndroidTV"
   StanVersion = "4.32.1"
)

var BaseUrl = []string{
   "aws.stan.video",
   "gec.stan.video",
   // these are geo block
   "023-stan.akamaized.net",
   "666-stan.akamaized.net",
}

type ActivationCode struct {
   Code string
   Url  string
}

func FetchActivationCode() (*ActivationCode, error) {
   // FIX: send generate=true as query param, not form body (matches Python device_code())
   req, err := http.NewRequest(
      "POST", "https://api.stan.com.au/login/v1/activation-codes/", nil,
   )
   if err != nil {
      return nil, err
   }
   req.URL.RawQuery = url.Values{"generate": {"true"}}.Encode()
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

// DTError is an error returned by the DRM Today license server.
type DTError struct {
   RespCode string // x-dt-resp-code header, e.g. "20101"
   Message  string // x-dt-error-message header, e.g. "not_granted"
}

func (e *DTError) Error() string {
   return fmt.Sprintf("drmtoday: %s (resp-code %s)", e.Message, e.RespCode)
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
   payload := url.Values{
      "audioCodecs":  {"aac"},
      "colorSpace":   {"hdr10"},
      "drm":          {"playready,widevine"},
      "hdcpVersion":  {"2.2"},
      "jwToken":      {w.JwToken},
      "manufacturer": {"NVIDIA"},
      "model":        {"SHIELD Android TV"},
      "os":           {"Android-9"},
      "screenSize":   {"3840x2160"},
      "stanName":     {StanName},
      "stanVersion":  {StanVersion},
      "type":         {"console"},
      "videoCodecs":  {"h264,decode,dovi,h263,h265,hevc,mjpeg,mpeg2v,mp4,mpeg4,vc1,vp8,vp9"},
   }.Encode()

   req, err := http.NewRequest(
      "POST", "https://api.stan.com.au/login/v1/sessions/app",
      strings.NewReader(payload),
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

// LicensePlayReady requests a PlayReady license. The response is raw XML.
func (m *Media) LicensePlayReady(data []byte) ([]byte, error) {
   req, err := http.NewRequest(
      "POST", m.Drm.LicenseServerUrl, bytes.NewReader(data),
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
   return io.ReadAll(resp.Body)
}

// LicenseWidevine requests a Widevine license. The response is JSON-wrapped.
func (m *Media) LicenseWidevine(data []byte) ([]byte, error) {
   req, err := http.NewRequest(
      "POST", m.Drm.LicenseServerUrl, bytes.NewReader(data),
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

   if resp.StatusCode != http.StatusOK {
      return nil, &DTError{
         RespCode: resp.Header.Get("x-dt-resp-code"),
         Message:  resp.Header.Get("x-dt-error-message"),
      }
   }

   var result struct {
      License []byte
   }
   if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
      return nil, err
   }
   return result.License, nil
}

func (*AppSession) CachePath() string {
   return "rosso/stan/AppSession"
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

// stan.go
