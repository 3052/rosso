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

// GetPlayReadyLicense fetches the PlayReady DRM license for the given title,
// unwraps the JSON response, decodes the base64, and returns the raw XML.
func GetPlayReadyLicense(actorAccessToken, titleId, playbackEnvelope string, licenseChallenge []byte) ([]byte, error) {
   payload := map[string]interface{}{
      "licenseChallenge": licenseChallenge,
      "playbackEnvelope": playbackEnvelope,
   }
   body, err := json.Marshal(payload)
   if err != nil {
      return nil, err
   }
   req, err := http.NewRequest(
      "POST",
      "https://atv-ps.primevideo.com/playback/drm-vod/GetPlayReadyLicense",
      bytes.NewReader(body),
   )
   if err != nil {
      return nil, err
   }
   query := url.Values{}
   query.Add("deviceTypeID", DeviceTypeID)
   query.Add("deviceID", DeviceID)
   req.Header.Set("Authorization", "Bearer "+actorAccessToken)
   req.URL.RawQuery = query.Encode()
   client := &http.Client{}
   resp, err := client.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()

   if resp.StatusCode != http.StatusOK {
      return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
   }

   var result struct {
      PlayReadyLicense struct {
         License string `json:"license"`
      } `json:"playReadyLicense"`
   }

   if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
      return nil, err
   }

   if result.PlayReadyLicense.License == "" {
      return nil, fmt.Errorf("empty license returned from API")
   }

   // The license is returned as a Base64 encoded XML string, decode it
   xmlBytes, err := base64.StdEncoding.DecodeString(result.PlayReadyLicense.License)
   if err != nil {
      return nil, fmt.Errorf("failed to decode base64 license: %w", err)
   }

   return xmlBytes, nil
}
