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

   req.Header.Set("Content-Type", "application/octet-stream")
   req.Header.Set("x-dt-auth-token", a.DRM.Token)
   req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:140.0) Gecko/20100101 Firefox/140.0")

   client := &http.Client{}
   resp, err := client.Do(req)
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
