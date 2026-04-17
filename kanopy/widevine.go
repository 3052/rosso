package kanopy

import (
   "bytes"
   "fmt"
   "io"
   "net/http"
)

func (s *Session) GetWidevine(manifest *Manifest, challenge []byte) ([]byte, error) {
   if manifest == nil {
      return nil, fmt.Errorf("a valid stream manifest is required to request a DRM license")
   }
   if manifest.DrmLicenseId == "" {
      return nil, fmt.Errorf("manifest does not contain a DRM license ID")
   }

   url := fmt.Sprintf("%s/kapi/licenses/widevine/%s", BaseUrl, manifest.DrmLicenseId)

   req, err := http.NewRequest("POST", url, bytes.NewBuffer(challenge))
   if err != nil {
      return nil, err
   }

   req.Header.Set("Authorization", "Bearer "+s.Jwt)
   req.Header.Set("User-Agent", UserAgent)
   req.Header.Set("X-Version", Xversion)

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
