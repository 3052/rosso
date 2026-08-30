package stan

import (
   "bytes"
   "encoding/json"
   "fmt"
   "io"
   "log"
   "net/http"
)

func do(req *http.Request) (*http.Response, error) {
   log.Println(req.Method, req.URL)
   return http.DefaultClient.Do(req)
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
