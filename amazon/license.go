package amazon

import (
   "bytes"
   "encoding/base64"
   "encoding/json"
   "fmt"
   "io"
   "net/http"
)

// GetPlayReadyLicense fetches the PlayReady DRM license for the given title,
// unwraps the JSON response, decodes the base64, and returns the raw XML.
func GetPlayReadyLicense(actorAccessToken, titleId, playbackEnvelope string, licenseChallenge []byte) ([]byte, error) {
   url := "https://atv-ps.primevideo.com/playback/drm-vod/GetPlayReadyLicense"

   req, err := http.NewRequest("POST", url, nil)
   if err != nil {
      return nil, err
   }

   query := req.URL.Query()
   query.Add("deviceTypeID", DeviceTypeID)
   query.Add("deviceID", DeviceID)
   query.Add("gascEnabled", "false")
   query.Add("marketplaceID", "ATVPDKIKX0DER")
   query.Add("uxLocale", "en_US")
   query.Add("firmware", "1")
   query.Add("titleId", titleId)
   req.URL.RawQuery = query.Encode()

   payload := map[string]interface{}{
      "packagingFormat":  "MPEG_DASH",
      "playbackEnvelope": playbackEnvelope,
      "licenseChallenge": licenseChallenge,
   }

   body, err := json.Marshal(payload)
   if err != nil {
      return nil, err
   }

   req.Body = io.NopCloser(bytes.NewBuffer(body))
   req.ContentLength = int64(len(body))

   req.Header.Set("User-Agent", UserAgent)
   req.Header.Set("Content-Type", "text/plain") // Required as text/plain for the DRM endpoint
   req.Header.Set("Accept", "*/*")
   req.Header.Set("Authorization", "Bearer "+actorAccessToken)

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

// GetWidevineLicense requests a Widevine DRM license from the Amazon endpoint.
func GetWidevineLicense(actorAccessToken, titleId, playbackEnvelope string, licenseChallenge []byte) ([]byte, error) {
   url := "https://ab8mt4dd97et.na.api.amazonvideo.com/playback/drm-vod/GetWidevineLicense"

   req, err := http.NewRequest("POST", url, nil)
   if err != nil {
      return nil, err
   }

   query := req.URL.Query()
   query.Add("deviceID", DeviceID)
   query.Add("deviceTypeID", DeviceTypeID)
   query.Add("gascEnabled", "false")
   query.Add("marketplaceID", "ATVPDKIKX0DER")
   query.Add("uxLocale", "en-US")
   query.Add("firmware", "1")
   query.Add("titleId", titleId)
   req.URL.RawQuery = query.Encode()

   payload := map[string]interface{}{
      "includeHdcpTestKey": true,
      "playbackEnvelope":   playbackEnvelope,
      "licenseChallenge":   licenseChallenge, // Go automatically encodes []byte to a base64 string
   }

   body, err := json.Marshal(payload)
   if err != nil {
      return nil, err
   }

   req.Body = io.NopCloser(bytes.NewBuffer(body))
   req.ContentLength = int64(len(body))

   req.Header.Set("User-Agent", UserAgent)
   req.Header.Set("Content-Type", "text/plain")
   req.Header.Set("Accept", "*/*")
   req.Header.Set("Authorization", "Bearer "+actorAccessToken)

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
      WidevineLicense struct {
         License []byte `json:"license"`
      } `json:"widevineLicense"`
   }

   if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
      return nil, err
   }

   if len(result.WidevineLicense.License) == 0 {
      return nil, fmt.Errorf("license not found in response")
   }

   return result.WidevineLicense.License, nil
}
