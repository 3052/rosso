package kanopy

import (
   "bytes"
   "fmt"
   "io"
   "net/http"
)

// GetWidevineLicense submits the binary CDM challenge to the Widevine DRM endpoint.
// It explicitly requires a Manifest (obtained from CreatePlay) to extract the DrmLicenseID.
func (s *Session) GetWidevineLicense(manifest *Manifest, challenge []byte) ([]byte, error) {
   if manifest == nil {
      return nil, fmt.Errorf("a valid stream manifest is required to request a DRM license")
   }
   if manifest.DrmLicenseID == "" {
      return nil, fmt.Errorf("manifest does not contain a DRM license ID")
   }

   url := fmt.Sprintf("%s/kapi/licenses/widevine/%s", BaseURL, manifest.DrmLicenseID)

   req, err := http.NewRequest("POST", url, bytes.NewBuffer(challenge))
   if err != nil {
      return nil, err
   }

   req.Header.Set("Authorization", "Bearer "+s.JWT)
   req.Header.Set("User-Agent", UserAgent)
   req.Header.Set("X-Version", XVersion)

   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()

   if resp.StatusCode != http.StatusOK {
      return nil, fmt.Errorf("widevine license request failed with status: %d", resp.StatusCode)
   }

   return io.ReadAll(resp.Body)
}
