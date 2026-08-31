package stan

import (
   "bytes"
   "encoding/json"
   "fmt"
   "io"
   "net/http"
   "net/url"
   "strconv"
   "strings"
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
   req, err := http.NewRequest(
      "POST", "https://api.stan.com.au/login/v1/activation-codes?generate=true", nil,
   )
   if err != nil {
      return nil, err
   }

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

func (a *ActivationCode) FetchToken() (*WebToken, error) {
   req, err := http.NewRequest("GET", a.Url, nil)
   if err != nil {
      return nil, err
   }
   resp, err := do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()

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

func (a *ActivationCode) String() string {
   var data strings.Builder
   data.WriteString("Stan.\n")
   data.WriteString("Log in with code\n")
   data.WriteString("1. Visit stan.com.au/activate\n")
   data.WriteString("2. Enter the code:\n")
   data.WriteString(a.Code)
   return data.String()
}

// DTError is an error returned by the DRM Today license server.
type DTError struct {
   RespCode string // x-dt-resp-code header, e.g. "20101"
   Message  string // x-dt-error-message header, e.g. "not_granted"
}

func (e *DTError) Error() string {
   return fmt.Sprintf("drmtoday: %s (resp-code %s)", e.Message, e.RespCode)
}

type Media struct {
   Drm *struct {
      CustomData       string
      KeyId            string
      LicenseServerUrl string
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

// DRM options: "widevine" or "playready".
// Quality options: "sd" (540p), "high" (1080p), "ultra" (2160p), "auto".
func (a *AppSession) FetchMedia(id int, quality, drm string) (*Media, error) {
   req, err := http.NewRequest(
      "GET", "https://api.stan.com.au/concurrency/v1/streams", nil,
   )
   if err != nil {
      return nil, err
   }
   req.Header.Set("x-forwarded-for", "1.128.0.0")
   req.URL.RawQuery = url.Values{
      "capabilities.drm": {drm},
      "format":           {"dash"}, // hls
      "jwToken":          {a.JwToken},
      "programId":        {strconv.Itoa(id)},
      "quality":          {quality},
   }.Encode()
   resp, err := do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()

   var result struct {
      Media  Media
      Errors ApiErrors
   }
   if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
      return nil, err
   }

   if len(result.Errors) > 0 {
      return nil, result.Errors
   }

   // Fail if the stream response didn't include DRM info
   if result.Media.Drm == nil {
      return nil, fmt.Errorf("stan: no DRM data in stream response")
   }

   return &result.Media, nil
}

func (*WebToken) CachePath() string {
   return "rosso/stan/WebToken"
}

// stan.go
