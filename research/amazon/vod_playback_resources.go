// --- get_vod_playback_resources.go ---
// Posts to the Playback Resource Service (PRS) using the envelope to retrieve the MPD URL.
package amazon

import (
   "bytes"
   "encoding/json"
   "fmt"
   "net/http"
   "net/url"
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

   payload := map[string]interface{}{
      "auditPingsRequest": map[string]interface{}{
         "device": map[string]string{
            "category": "Phone",
            "platform": "Android",
         },
      },
      "globalParameters": map[string]interface{}{
         "capabilityDiscriminators": map[string]interface{}{
            "discriminators": map[string]interface{}{
               "hardware": map[string]string{
                  "chipset":      "goldfish_x86_64",
                  "manufacturer": "Google",
                  "modelName":    "sdk_gphone_x86_64",
               },
               "software": map[string]interface{}{
                  "application": map[string]string{
                     "name":    "com.amazon.avod.thirdpartyclient",
                     "version": "458",
                  },
                  "client": map[string]interface{}{"id": nil},
                  "firmware": map[string]string{
                     "version": "google/sdk_gphone_x86_64/generic_x86_64_arm64:11/RSR1.240422.006/12134477:userdebug/dev-keys",
                  },
                  "operatingSystem": map[string]string{
                     "name":    "Android",
                     "version": "11",
                  },
                  "player": map[string]string{
                     "name":    "Android Player",
                     "version": "3.0.458.357",
                  },
                  "renderer": map[string]string{
                     "drmScheme": "WIDEVINE",
                     "name":      "MCMD",
                  },
               },
            },
            "version": 1,
         },
      },
      "deviceCapabilityFamily": "AndroidPlayer",
      "playbackEnvelope":       playbackEnvelope,
      "playbackDataRequest":    map[string]interface{}{},
      "timedTextUrlsRequest": map[string]interface{}{
         "supportedTimedTextFormats": []string{"TTMLv2", "DFXP"},
      },
      "transitionTimecodesRequest": map[string]interface{}{},
      "trickplayUrlsRequest":       map[string]interface{}{},
      "vodPlaybackUrlsRequest": map[string]interface{}{
         "ads": map[string]interface{}{
            "advertisingId":      "738e5ee9-5d04-49b3-80fd-1db41971a255",
            "appBundle":          "ATVAndroid3P",
            "appStoreUrl":        "http://www.samsungapps.com/appquery/appDetail.as?appId=ATVAndroid3P",
            "gdpr":               map[string]interface{}{"consentMap": nil, "enabled": false},
            "optOutOfAdTracking": false,
         },
         "device": map[string]interface{}{
            "displayBasedVending": "supported",
            "displayHeight":       1080,
            "displayWidth":        2340,
            "streamingTechnologies": map[string]interface{}{
               "DASH": map[string]interface{}{
                  "edgeDeliveryAuthorizationSchemes":      nil,
                  "fragmentRepresentations":               []string{"ByteOffsetRange", "SeparateFile"},
                  "manifestThinningToSupportedResolution": "Forbidden",
                  "segmentInfoType":                       "List",
                  "stitchType":                            "MultiPeriod",
                  "timedTextRepresentations":              []string{"BurnedIn", "NotInManifestNorStream", "SeparateStreamInManifest"},
                  "trickplayRepresentations":              []string{"NotInManifestNorStream"},
                  "variableAspectRatio":                   "supported",
                  "vastTimelineType":                      "Absolute",
                  "bitrateAdaptations":                    []string{"CBR", "CVBR"},
                  "codecs":                                []string{"H264", "H265"},
                  "drmKeyScheme":                          "DualKey",
                  "drmStrength":                           "L10",
                  "drmType":                               "WIDEVINE",
                  "dynamicRangeFormats":                   []string{"None"},
                  "frameRates":                            []string{"Standard"},
               },
            },
            "acceptedCreativeApis":           []string{},
            "category":                       "Phone",
            "hdcpLevel":                      "no_ports",
            "maxVideoResolution":             "576p",
            "operatingSystem":                "Android11",
            "platform":                       "Android",
            "supportedStreamingTechnologies": []string{"DASH"},
         },
      },
      "playbackCustomizations": map[string]interface{}{
         "desiredAudioTracks": map[string]interface{}{
            "languageCodes":     []string{"en-us", "en"},
            "audioSubTypes":     []string{"dialog"},
            "maximumTrackCount": 1,
         },
      },
      "playbackSettingsRequest": map[string]interface{}{
         "chipset":               "goldfish_x86_64",
         "deviceModel":           "sdk_gphone_x86_64",
         "firmware":              "fmw:30-app:3.0.458.357",
         "responseFormatVersion": "1.0.0",
         "heuristicProfile":      `{"STARTUP_TIME":"PRIORITY","BUFFERING_RISK":"LOW","QUALITY":"LOW"}`,
         "playerType":            "Android Player",
         "softwareVersion":       "458",
         "titleId":               titleID,
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

   req.Header.Set("x-gasc-enabled", "true")
   req.Header.Set("x-request-priority", "CRITICAL")
   req.Header.Set("Accept", "application/json")
   req.Header.Set("User-Agent", "Dalvik/2.1.0 (Linux; U; Android 11; sdk_gphone_x86_64 Build/RSR1.240422.006)")
   req.Header.Set("Accept-Language", "en_US")
   req.Header.Set("x-retry-count", "0")
   req.Header.Set("Authorization", "Bearer "+bearerToken)
   req.Header.Set("Content-Type", "application/json; charset=utf-8")
   req.Header.Set("Connection", "Keep-Alive")
   req.Header.Set("Accept-Encoding", "identity")

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

   urlSets := resourceResp.VodPlaybackUrls.Result.PlaybackUrls.UrlSets
   if len(urlSets) > 0 && urlSets[0].Url != "" {
      return urlSets[0].Url, nil
   }

   return "", fmt.Errorf("MPD URL not found in PRS response")
}
