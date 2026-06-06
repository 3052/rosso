package amazon

import (
   "bytes"
   "encoding/json"
   "fmt"
   "net/http"
)

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

func GetMPDUrl(actorAccessToken, titleId, deviceId, playbackEnvelope string) (string, error) {
   url := fmt.Sprintf("https://abzq7aq4866p.na.api.amazonvideo.com/playback/prs/GetVodPlaybackResources?consumptionType=STREAMING&deviceID=%s&deviceTypeID=A43PXU4ZN2AL1&format=json&titleId=%s&version=1&videoMaterialType=Feature", deviceId, titleId)

   payload := map[string]interface{}{
      "globalParameters": map[string]interface{}{
         "capabilityDiscriminators": map[string]interface{}{
            "discriminators": map[string]interface{}{
               "software": map[string]interface{}{
                  "player":   map[string]string{"name": "Android Player", "version": "3.0.458.357"},
                  "renderer": map[string]string{"drmScheme": "WIDEVINE", "name": "MCMD"},
               },
            },
         },
         "version": 1,
      },
      "deviceCapabilityFamily": "AndroidPlayer",
      "playbackEnvelope":       playbackEnvelope,
      "vodPlaybackUrlsRequest": map[string]interface{}{
         "device": map[string]interface{}{
            "streamingTechnologies": map[string]interface{}{
               "DASH": map[string]interface{}{
                  "codecs":              []string{"H264", "H265"},
                  "drmKeyScheme":        "DualKey",
                  "drmType":             "WIDEVINE",
                  "dynamicRangeFormats": []string{"None"},
               },
            },
            "supportedStreamingTechnologies": []string{"DASH"},
         },
      },
   }

   body, _ := json.Marshal(payload)
   req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
   if err != nil {
      return "", err
   }

   req.Header.Set("Authorization", "Bearer "+actorAccessToken)
   req.Header.Set("Content-Type", "application/json; charset=utf-8")
   req.Header.Set("User-Agent", "Dalvik/2.1.0 (Linux; U; Android 11; sdk_gphone_x86_64 Build/RSR1.240422.006)")

   client := &http.Client{}
   resp, err := client.Do(req)
   if err != nil {
      return "", err
   }
   defer resp.Body.Close()

   if resp.StatusCode != http.StatusOK {
      return "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
   }

   var prsResp PlaybackResourcesResponse
   if err := json.NewDecoder(resp.Body).Decode(&prsResp); err != nil {
      return "", err
   }

   if len(prsResp.VodPlaybackUrls.Result.PlaybackUrls.UrlSets) == 0 {
      return "", fmt.Errorf("no URL sets found in playback resources response")
   }

   return prsResp.VodPlaybackUrls.Result.PlaybackUrls.UrlSets[0].Url, nil
}
