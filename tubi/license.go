// request_license_challenge.go
package tubi

import (
   "bytes"
   "fmt"
   "io"
   "net/http"
)

func (s *LicenseServer) PostLicense(payload []byte) ([]byte, error) {
   if s == nil || s.Url == "" {
      return nil, fmt.Errorf("invalid or missing server URL")
   }

   req, err := http.NewRequest("POST", s.Url, bytes.NewReader(payload))
   if err != nil {
      return nil, fmt.Errorf("failed to create request: %w", err)
   }

   req.Header.Set("content-type", "application/x-protobuf")
   req.Header.Set("user-agent", "Go-http-client/2.0")

   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, fmt.Errorf("request failed: %w", err)
   }
   defer resp.Body.Close()

   if resp.StatusCode != http.StatusOK {
      return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
   }

   body, err := io.ReadAll(resp.Body)
   if err != nil {
      return nil, fmt.Errorf("failed to read response body: %w", err)
   }

   return body, nil
}
