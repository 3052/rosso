package molotov

import (
   "bytes"
   "fmt"
   "io"
   "log"
   "net/http"
)

// DeviceID is the centralized value used for the x-device-id header across all requests.
const DeviceID = "x-device-id"

const x_forwarded_for = "178.132.106.134"

// doRequest logs the method and URL, then performs the HTTP request.
func doRequest(req *http.Request) (*http.Response, error) {
   log.Println(req.Method, req.URL)
   client := &http.Client{}
   return client.Do(req)
}

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
