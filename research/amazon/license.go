package amazon

import (
   "bytes"
   "encoding/base64"
   "encoding/json"
   "fmt"
   "io"
   "net/http"
   "net/url"
)

// LicenseResponse defines the structure to extract the base64 license.
type LicenseResponse struct {
   WidevineLicense struct {
      License string `json:"license"`
   } `json:"widevineLicense"`
   PlayReadyLicense struct {
      License string `json:"license"`
   } `json:"playReadyLicense"`
}

// GetLicense submits the CDM challenge and retrieves the base64 encoded license.
// `challenge` expects raw bytes for Widevine or raw XML/SOAP bytes for PlayReady.
func (c *Client) GetLicense(p DeviceProfile, envelope string, challenge []byte) (string, error) {
   if p.AuthBearer == "" {
      return "", fmt.Errorf("AuthBearer is required")
   }

   endpoint := "/playback/drm-vod/GetWidevineLicense"
   if p.DRMType == "PlayReady" {
      endpoint = "/playback/drm-vod/GetPlayReadyLicense"
   }

   u := url.URL{
      Scheme: "https",
      Host:   defaultAPIHost,
      Path:   endpoint,
   }
   q := u.Query()
   q.Set("deviceID", p.DeviceID)
   q.Set("deviceTypeID", defaultDeviceTypeID) // Centralized
   u.RawQuery = q.Encode()

   payload := map[string]any{
      "playbackEnvelope": envelope,
      "licenseChallenge": base64.StdEncoding.EncodeToString(challenge),
   }

   bodyBytes, err := json.Marshal(payload)
   if err != nil {
      return "", err
   }

   req, err := http.NewRequest("POST", u.String(), bytes.NewBuffer(bodyBytes))
   if err != nil {
      return "", err
   }

   req.Header.Set("Content-Type", "application/json")
   req.Header.Set("Authorization", "Bearer "+p.AuthBearer)

   resp, err := c.HTTPClient.Do(req)
   if err != nil {
      return "", err
   }
   defer resp.Body.Close()

   if resp.StatusCode != http.StatusOK {
      body, _ := io.ReadAll(resp.Body)
      return "", fmt.Errorf("bad status %d: %s", resp.StatusCode, string(body))
   }

   var result LicenseResponse
   if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
      return "", err
   }

   license := result.WidevineLicense.License
   if p.DRMType == "PlayReady" {
      license = result.PlayReadyLicense.License
   }

   if license == "" {
      return "", fmt.Errorf("could not find license string in JSON response")
   }

   return license, nil
}
