package amazon

import (
   "bytes"
   "encoding/json"
   "fmt"
   "net/http"
   "net/url"
)

// GetWidevineLicense requests a Widevine DRM license from the Amazon endpoint.
func GetWidevineLicense(actorAccessToken, playbackEnvelope string, licenseChallenge []byte) ([]byte, error) {
   reqURL := "https://ab8mt4dd97et.na.api.amazonvideo.com/playback/drm-vod/GetWidevineLicense"
   payload := map[string]interface{}{
      "playbackEnvelope":   playbackEnvelope,
      "licenseChallenge":   licenseChallenge,
   }
   query := url.Values{}
   query.Add("deviceTypeID", DeviceTypeID)
   query.Add("deviceID", DeviceID)
   return fetchDRMLicense(reqURL, actorAccessToken, query, payload)
}

// GetPlayReadyLicense fetches the PlayReady DRM license for the given title.
func GetPlayReadyLicense(actorAccessToken, playbackEnvelope string, licenseChallenge []byte) ([]byte, error) {
   reqURL := "https://atv-ps.primevideo.com/playback/drm-vod/GetPlayReadyLicense"
   payload := map[string]interface{}{
      "playbackEnvelope": playbackEnvelope,
      "licenseChallenge": licenseChallenge,
   }
   query := url.Values{}
   query.Add("deviceTypeID", DeviceTypeID)
   query.Add("deviceID", DeviceID)
   return fetchDRMLicense(reqURL, actorAccessToken, query, payload)
}

// fetchDRMLicense is the shared base function for making DRM requests
func fetchDRMLicense(reqURL, actorAccessToken string, query url.Values, payload map[string]interface{}) ([]byte, error) {
   body, err := json.Marshal(payload)
   if err != nil {
      return nil, fmt.Errorf("failed to marshal payload: %w", err)
   }
   // bytes.NewReader automatically handles ContentLength and Read closures
   req, err := http.NewRequest(http.MethodPost, reqURL, bytes.NewReader(body))
   if err != nil {
      return nil, fmt.Errorf("failed to create request: %w", err)
   }
   req.URL.RawQuery = query.Encode()
   // Standardized headers for both Widevine and PlayReady
   req.Header.Set("Authorization", "Bearer "+actorAccessToken)
   client := &http.Client{}
   resp, err := client.Do(req)
   if err != nil {
      return nil, fmt.Errorf("request failed: %w", err)
   }
   defer resp.Body.Close()

   if resp.StatusCode != http.StatusOK {
      return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
   }
   // By using a map, we can dynamically handle either "widevineLicense" or "playReadyLicense" top-level keys.
   // Go automatically base64-decodes JSON strings when unmarshaling into a []byte!
   var result map[string]struct {
      License []byte `json:"license"`
   }
   if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
      return nil, fmt.Errorf("failed to decode response: %w", err)
   }
   // Extract and return whichever license was provided in the response
   if wv, ok := result["widevineLicense"]; ok && len(wv.License) > 0 {
      return wv.License, nil
   }
   if pr, ok := result["playReadyLicense"]; ok && len(pr.License) > 0 {
      return pr.License, nil
   }
   return nil, fmt.Errorf("license not found in response")
}
