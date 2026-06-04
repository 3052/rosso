// --- get_vod_playback_resources.go ---
package amazon

import (
   "bytes"
   "encoding/json"
   "fmt"
   "net/http"
   "net/url"
)

// PlaybackResourcesResponse represents the structure needed to extract the MPD URL
type PlaybackResourcesResponse struct {
   VodPlaybackUrls struct {
      Result struct {
         PlaybackUrls struct {
            UrlSets []struct {
               Url string `json:"url"`
            } `json:"urlSets"`
         } `json:"playbackUrls"`
      } `json:"result"`
   } `json:"vodPlaybackUrls"`
}

// GetMPDUrl makes the POST request to PRS (Playback Resource Service) to get the MPD URL
func GetMPDUrl(titleID, deviceID, bearerToken, playbackEnvelope string) (string, error) {
   baseURL := "https://abzq7aq4866p.na.api.amazonvideo.com/playback/prs/GetVodPlaybackResources"

   q := url.Values{}
   q.Add("consumptionType", "STREAMING")
   q.Add("deviceID", deviceID)
   q.Add("deviceTypeID", "A43PXU4ZN2AL1")
   q.Add("firmware", "fmw:30-app:3.0.458.357")
   q.Add("format", "json")
   q.Add("osLocale", "en_US")
   q.Add("softwareVersion", "458")
   q.Add("titleId", titleID)
   q.Add("uxLocale", "en_US")
   q.Add("version", "1")
   q.Add("videoMaterialType", "Feature")

   // Construct the JSON payload containing the required playback envelope
   payload := map[string]interface{}{
      "deviceCapabilityFamily": "AndroidPlayer",
      "playbackEnvelope":       playbackEnvelope,
      "vodPlaybackUrlsRequest": map[string]interface{}{
         "device": map[string]interface{}{
            "supportedStreamingTechnologies": []string{"DASH"},
         },
      },
      "playbackSettingsRequest": map[string]interface{}{
         "titleId": titleID,
      },
   }

   bodyBytes, err := json.Marshal(payload)
   if err != nil {
      return "", err
   }

   req, err := http.NewRequest("POST", fmt.Sprintf("%s?%s", baseURL, q.Encode()), bytes.NewBuffer(bodyBytes))
   if err != nil {
      return "", err
   }

   req.Header.Set("Authorization", "Bearer "+bearerToken)
   req.Header.Set("Content-Type", "application/json; charset=utf-8")
   req.Header.Set("Accept", "application/json")
   req.Header.Set("User-Agent", "Dalvik/2.1.0 (Linux; U; Android 11; sdk_gphone_x86_64 Build/RSR1.240422.006)")

   client := &http.Client{}
   resp, err := client.Do(req)
   if err != nil {
      return "", err
   }
   defer resp.Body.Close()

   if resp.StatusCode != http.StatusOK {
      return "", fmt.Errorf("failed to get playback resources, status: %d", resp.StatusCode)
   }

   var resourceResp PlaybackResourcesResponse
   if err := json.NewDecoder(resp.Body).Decode(&resourceResp); err != nil {
      return "", err
   }

   // Extract the first MPD URL from the urlSets array
   urlSets := resourceResp.VodPlaybackUrls.Result.PlaybackUrls.UrlSets
   if len(urlSets) > 0 && urlSets[0].Url != "" {
      return urlSets[0].Url, nil
   }

   return "", fmt.Errorf("MPD URL not found in PRS response")
}
