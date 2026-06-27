package amazon

import (
   "bytes"
   "encoding/json"
   "fmt"
   "net/http"
   "net/url"
)

// GetPlayReadyLicense fetches the PlayReady DRM license for the given title.
func GetPlayReadyLicense(actorAccessToken, playbackEnvelope string, licenseChallenge []byte) ([]byte, error) {
   reqURL := HostATVPS + "/playback/drm-vod/GetPlayReadyLicense"
   payload := map[string]any{
      "playbackEnvelope": playbackEnvelope,
      "licenseChallenge": licenseChallenge,
   }
   query := url.Values{}
   query.Add("deviceTypeID", DeviceTypeID)
   query.Add("deviceID", DeviceID)
   return fetchDRMLicense(reqURL, actorAccessToken, query, payload)
}

// GetWidevineLicense requests a Widevine DRM license from the Amazon endpoint.
func GetWidevineLicense(actorAccessToken, playbackEnvelope string, licenseChallenge []byte) ([]byte, error) {
   reqURL := HostATVPS + "/playback/drm-vod/GetWidevineLicense"
   payload := map[string]any{
      "playbackEnvelope": playbackEnvelope,
      "licenseChallenge": licenseChallenge,
   }
   query := url.Values{}
   query.Add("deviceTypeID", DeviceTypeID)
   query.Add("deviceID", DeviceID)
   return fetchDRMLicense(reqURL, actorAccessToken, query, payload)
}

// fetchDRMLicense is the shared base function for making DRM requests
func fetchDRMLicense(reqURL, actorAccessToken string, query url.Values, payload map[string]any) ([]byte, error) {
   body, err := json.Marshal(payload)
   if err != nil {
      return nil, fmt.Errorf("failed to marshal payload: %w", err)
   }

   req, err := http.NewRequest(http.MethodPost, reqURL, bytes.NewReader(body))
   if err != nil {
      return nil, fmt.Errorf("failed to create request: %w", err)
   }
   req.URL.RawQuery = query.Encode()
   req.Header.Set("Authorization", "Bearer "+actorAccessToken)

   client := &http.Client{}
   resp, err := client.Do(req)
   if err != nil {
      return nil, fmt.Errorf("request failed: %w", err)
   }
   defer resp.Body.Close()

   var result struct {
      WidevineLicense *struct {
         License []byte `json:"license"`
      } `json:"widevineLicense"`
      PlayReadyLicense *struct {
         License []byte `json:"license"`
      } `json:"playReadyLicense"`
      Message *struct {
         Body *struct {
            Code    string `json:"code"`
            Message string `json:"message"`
         } `json:"body"`
      } `json:"message"`
   }

   if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
      return nil, fmt.Errorf("failed to decode response (status %d): %w", resp.StatusCode, err)
   }

   // 1. Check for the structured JSON API error
   if result.Message != nil && result.Message.Body != nil {
      return nil, fmt.Errorf("API error [%s]: %s", result.Message.Body.Code, result.Message.Body.Message)
   }

   // 2. Check for standard HTTP errors if no JSON error message was provided
   if resp.StatusCode != http.StatusOK {
      return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
   }

   // 3. Extract and return whichever license was provided
   if result.WidevineLicense != nil && len(result.WidevineLicense.License) > 0 {
      return result.WidevineLicense.License, nil
   }
   if result.PlayReadyLicense != nil && len(result.PlayReadyLicense.License) > 0 {
      return result.PlayReadyLicense.License, nil
   }

   return nil, fmt.Errorf("license not found in response")
}
