package kanopy

import (
   "bytes"
   "fmt"
   "io"
   "net/http"
)

// GetWidevineLicense submits the CDM challenge to the DRM endpoint and returns the Widevine License payload.
func (c *Client) GetWidevineLicense(drmLicenseID string, challenge []byte) ([]byte, error) {
   url := fmt.Sprintf("%s/kapi/licenses/widevine/%s", BaseURL, drmLicenseID)

   req, err := http.NewRequest("POST", url, bytes.NewBuffer(challenge))
   if err != nil {
      return nil, err
   }

   req.Header.Set("Authorization", "Bearer "+c.Token)
   req.Header.Set("User-Agent", c.UserAgent)
   req.Header.Set("X-Version", c.XVersion)

   // Explicitly using http.DefaultClient
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
