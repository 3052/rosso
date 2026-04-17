// request_license_challenge.go
package tubi

import (
   "bytes"
   "fmt"
   "io"
   "net/http"
   "strconv"
)

// PostLicenseChallenge sends the Widevine DRM challenge to the dynamic license URL
// extracted from the CMS request, using the payload provided by the caller.
func PostLicenseChallenge(licenseURL string, payload []byte) ([]byte, error) {
   req, err := http.NewRequest("POST", licenseURL, bytes.NewReader(payload))
   if err != nil {
      return nil, fmt.Errorf("failed to create request: %w", err)
   }

   req.Header.Set("content-type", "application/x-protobuf")
   req.Header.Set("content-length", strconv.Itoa(len(payload)))
   req.Header.Set("accept-encoding", "gzip")
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
