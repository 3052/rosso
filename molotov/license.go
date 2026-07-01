// license.go
package molotov

import (
   "bytes"
   "fmt"
   "io"
   "net/http"
)

// GetLicense requests the DRM license. As a method on *AssetResponse,
// it can be used directly as a closure: func([]byte) ([]byte, error).
func (a *AssetResponse) GetLicense(challenge []byte) ([]byte, error) {
   req, err := http.NewRequest("POST", a.DRM.LicenseURL, bytes.NewReader(challenge))
   if err != nil {
      return nil, err
   }
   resp, err := doRequest(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()

   if resp.StatusCode != http.StatusOK {
      return nil, fmt.Errorf("license request failed with status: %d", resp.StatusCode)
   }

   licenseData, err := io.ReadAll(resp.Body)
   if err != nil {
      return nil, err
   }

   return licenseData, nil
}
