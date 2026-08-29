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
   LicenseWidevine  = "https://lic.drmtoday.com/license-proxy-widevine/cenc/"
   LicensePlayReady = "https://lic.staging.drmtoday.com/license-proxy-headerauth/drmtoday/RightsManager.asmx"
)

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
      CustomData string
      KeyId      string
   }
   VideoUrl string
}

func (m *Media) License(url string, data []byte) ([]byte, error) {
   req, err := http.NewRequest(
      // final slash is needed
      "POST", url,
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

func (m *Media) LicensePlayReady(data []byte) ([]byte, error) {
   return m.License(LicensePlayReady, data)
}

func (m *Media) LicenseWidevine(data []byte) ([]byte, error) {
   return m.License(LicenseWidevine, data)
}

// Profile represents a Stan user profile.
type Profile struct {
   Id            string
   Name          string
   IsKidsProfile bool
}

// quality.go
