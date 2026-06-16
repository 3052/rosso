package amazon

import (
   "bytes"
   "encoding/json"
   "fmt"
   "io"
   "net/http"
)

// GetWidevineLicense requests a Widevine DRM license from the Amazon endpoint.
func GetWidevineLicense(actorAccessToken, titleId, playbackEnvelope string, licenseChallenge []byte) ([]byte, error) {
   url := "https://ab8mt4dd97et.na.api.amazonvideo.com/playback/drm-vod/GetWidevineLicense"

   req, err := http.NewRequest("POST", url, nil)
   if err != nil {
      return nil, err
   }

   q := req.URL.Query()
   q.Add("deviceID", "uuidcbb2f9705f13437e9e515622dce02106")
   q.Add("deviceTypeID", "A2SNKIF736WF4T")
   q.Add("gascEnabled", "false")
   q.Add("marketplaceID", "ATVPDKIKX0DER")
   q.Add("uxLocale", "en-US")
   q.Add("firmware", "1")
   q.Add("titleId", titleId)
   req.URL.RawQuery = q.Encode()

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

   req.Header.Set("User-Agent", "Android/google/sdk_gphone_x86/generic_x86_arm:11/RSR1.240422.006/12134477:userdebug/dev-keys, Ignition X/15.5.2026042820-android, Google")
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

// GetPlayReadyLicense requests a PlayReady DRM license from the Amazon endpoint.
func GetPlayReadyLicense(actorAccessToken, titleId, playbackEnvelope string, licenseChallenge []byte) ([]byte, error) {
   url := "https://ab8mt4dd97et.na.api.amazonvideo.com/playback/drm-vod/GetPlayReadyLicense"

   req, err := http.NewRequest("POST", url, nil)
   if err != nil {
      return nil, err
   }

   q := req.URL.Query()
   q.Add("deviceID", "uuidcbb2f9705f13437e9e515622dce02106")
   q.Add("deviceTypeID", "A2SNKIF736WF4T")
   q.Add("gascEnabled", "false")
   q.Add("marketplaceID", "ATVPDKIKX0DER")
   q.Add("uxLocale", "en-US")
   q.Add("firmware", "1")
   q.Add("titleId", titleId)
   req.URL.RawQuery = q.Encode()

   payload := map[string]interface{}{
      "packagingFormat":  "MPEG_DASH",
      "playbackEnvelope": playbackEnvelope,
      "licenseChallenge": licenseChallenge, // Go automatically encodes []byte to a base64 string
   }

   body, err := json.Marshal(payload)
   if err != nil {
      return nil, err
   }

   req.Body = io.NopCloser(bytes.NewBuffer(body))
   req.ContentLength = int64(len(body))

   req.Header.Set("User-Agent", "Android/google/sdk_gphone_x86/generic_x86_arm:11/RSR1.240422.006/12134477:userdebug/dev-keys, Ignition X/15.5.2026042820-android, Google")
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
      PlayReadyLicense struct {
         License []byte `json:"license"` // Go automatically decodes the base64 string back into []byte
      } `json:"playReadyLicense"`
   }

   if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
      return nil, err
   }

   if len(result.PlayReadyLicense.License) == 0 {
      return nil, fmt.Errorf("license not found in response")
   }

   return result.PlayReadyLicense.License, nil
}
