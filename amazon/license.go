package amazon

import (
   "bytes"
   "encoding/json"
   "fmt"
   "io"
   "net/http"
   "net/url"
)

// GetPlayReadyLicense fetches the PlayReady DRM license for the given title.
func GetPlayReadyLicense(actorAccessToken, playbackEnvelope string, licenseChallenge []byte, deviceTypeID string) ([]byte, error) {
   reqURL := HostATVPS + "/playback/drm-vod/GetPlayReadyLicense"
   payload := map[string]any{
      "playbackEnvelope": playbackEnvelope,
      "licenseChallenge": licenseChallenge,
   }
   query := url.Values{}
   query.Add("deviceTypeID", deviceTypeID) // Updated to use function input
   query.Add("deviceID", DeviceID)
   return fetchDRMLicense(reqURL, actorAccessToken, query, payload)
}

// GetWidevineLicense requests a Widevine DRM license from the Amazon endpoint.
func GetWidevineLicense(actorAccessToken, playbackEnvelope string, licenseChallenge []byte, deviceTypeID string) ([]byte, error) {
   reqURL := HostATVPS + "/playback/drm-vod/GetWidevineLicense"
   payload := map[string]any{
      "playbackEnvelope": playbackEnvelope,
      "licenseChallenge": licenseChallenge,
   }
   query := url.Values{}
   query.Add("deviceTypeID", deviceTypeID) // Updated to use function input
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

   // Read the body once so we can attempt multiple unmarshals
   respBytes, err := io.ReadAll(resp.Body)
   if err != nil {
      return nil, fmt.Errorf("failed to read response: %w", err)
   }

   // 1. Try the standard response format (contains licenses or a nested error object)
   var standardResp struct {
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

   if err := json.Unmarshal(respBytes, &standardResp); err == nil {
      if standardResp.Message != nil && standardResp.Message.Body != nil {
         return nil, fmt.Errorf("API error [%s]: %s", standardResp.Message.Body.Code, standardResp.Message.Body.Message)
      }
      if standardResp.WidevineLicense != nil && len(standardResp.WidevineLicense.License) > 0 {
         return standardResp.WidevineLicense.License, nil
      }
      if standardResp.PlayReadyLicense != nil && len(standardResp.PlayReadyLicense.License) > 0 {
         return standardResp.PlayReadyLicense.License, nil
      }
   }

   // 2. If the first unmarshal fails (e.g., "message" is a string causing a type error), try the flat error format
   var flatErrorResp struct {
      Code    string `json:"code"`
      ID      string `json:"id"`
      Message string `json:"message"`
   }

   if err := json.Unmarshal(respBytes, &flatErrorResp); err == nil && flatErrorResp.Message != "" {
      return nil, fmt.Errorf("code: %s, message: %s, id: %s", flatErrorResp.Code, flatErrorResp.Message, flatErrorResp.ID)
   }

   // 3. Check for standard HTTP errors if no JSON error message was extracted
   if resp.StatusCode != http.StatusOK {
      return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
   }

   return nil, fmt.Errorf("license not found in response")
}
