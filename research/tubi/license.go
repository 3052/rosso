// request_license_challenge.go
package tubi

import (
   "bytes"
   "fmt"
   "io"
   "net/http"
)

// PostLicenseChallenge sends the Widevine DRM challenge using the URL from the LicenseServer descendant.
func PostLicenseChallenge(server *LicenseServer, payload []byte) ([]byte, error) {
   if server == nil || server.URL == "" {
      return nil, fmt.Errorf("invalid or missing server URL")
   }

   request, err := http.NewRequest("POST", server.URL, bytes.NewReader(payload))
   if err != nil {
      return nil, fmt.Errorf("failed to create request: %w", err)
   }

   request.Header.Set("content-type", "application/x-protobuf")
   request.Header.Set("user-agent", "Go-http-client/2.0")

   response, err := http.DefaultClient.Do(request)
   if err != nil {
      return nil, fmt.Errorf("request failed: %w", err)
   }
   defer response.Body.Close()

   if response.StatusCode != http.StatusOK {
      return nil, fmt.Errorf("unexpected status code: %d", response.StatusCode)
   }

   body, err := io.ReadAll(response.Body)
   if err != nil {
      return nil, fmt.Errorf("failed to read response body: %w", err)
   }

   return body, nil
}
