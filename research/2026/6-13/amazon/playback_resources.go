package amazon

import (
   "bytes"
   "encoding/json"
   "fmt"
   "io"
   "net/http"
)

// GetVodPlaybackResources fetches the final MPD URL for playback.
func GetVodPlaybackResources(actorAccessToken, titleId, playbackEnvelope string) (string, error) {
   url := "https://ab8mt4dd97et.na.api.amazonvideo.com/playback/prs/GetVodPlaybackResources"

   req, err := http.NewRequest("POST", url, nil)
   if err != nil {
      return "", err
   }

   q := req.URL.Query()
   q.Add("deviceTypeID", "A2SNKIF736WF4T")
   q.Add("deviceID", "uuidcbb2f9705f13437e9e515622dce02106")
   q.Add("firmware", "1")
   q.Add("titleId", titleId)
   req.URL.RawQuery = q.Encode()

   payload := map[string]interface{}{
      "globalParameters": map[string]interface{}{
         "deviceCapabilityFamily": "LivingRoomPlayer",
         "playbackEnvelope":       playbackEnvelope,
         "capabilityDiscriminators": map[string]interface{}{
            "operatingSystem": map[string]string{"name": "Android", "version": "11"},
            "deviceModel":     map[string]string{"name": "sdk_gphone_x86", "version": "UNKNOWN"},
            "middleware":      map[string]string{"name": "Ignite", "version": "15.5.2026042820-android"},
         },
      },
      "auditPingsRequest":                 map[string]interface{}{},
      "widevineServiceCertificateRequest": map[string]interface{}{},
      "playbackDataRequest":               map[string]interface{}{},
      "timedTextUrlsRequest": map[string]interface{}{
         "supportedTimedTextFormats": []string{"TTMLv2", "DFXP"},
      },
      "trickplayUrlsRequest":       map[string]interface{}{},
      "transitionTimecodesRequest": map[string]interface{}{},
      "vodPlaylistedPlaybackUrlsRequest": map[string]interface{}{
         "device": map[string]interface{}{
            "hdcpLevel":                      "1.4",
            "maxVideoResolution":             "576p",
            "supportedStreamingTechnologies": []string{"DASH"},
            "streamingTechnologies": map[string]interface{}{
               "DASH": map[string]interface{}{
                  "bitrateAdaptations":               []string{"CBR", "CVBR"},
                  "codecs":                           []string{"H264"},
                  "drmKeyScheme":                     "DualKey",
                  "drmType":                          "Widevine",
                  "dynamicRangeFormats":              []string{"None"},
                  "edgeDeliveryAuthorizationSchemes": []string{"PVExchangeV1", "Transparent"},
                  "fragmentRepresentations":          []string{"ByteOffsetRange", "SeparateFile"},
                  "frameRates":                       []string{"Standard"},
                  "stitchType":                       "MultiPeriod",
                  "segmentInfoType":                  "Base",
                  "timedTextRepresentations":         []string{"NotInManifestNorStream", "SeparateStreamInManifest"},
                  "trickplayRepresentations":         []string{"NotInManifestNorStream"},
                  "variableAspectRatio":              "supported",
               },
            },
            "displayWidth":         1080,
            "displayHeight":        1080,
            "acceptedCreativeApis": []int{1006, 1008},
         },
         "ads": map[string]interface{}{
            "advertisingId":      "aff7331b-3bdf-476f-ae78-386b5d55e0e5",
            "appBundle":          "com.primevideo.Google",
            "appStoreUrl":        nil,
            "optOutOfAdTracking": false,
            "gdpr": map[string]interface{}{
               "enabled":    false,
               "consentMap": map[string]interface{}{},
            },
            "mainContentResumeOffsetHintMillis": 0,
            "playerContractVersion":             1,
         },
         "playbackCustomizations": map[string]interface{}{
            "capVideoDefinition": "SD",
         },
         "playbackSettingsRequest": map[string]interface{}{
            "deviceModel":           "sdk_gphone_x86",
            "firmware":              "google/sdk_gphone_x86/generic_x86_arm:11/RSR1.240422.006/12134477:userdebug/dev-keys",
            "heuristicProfile":      "{\"Quality\":\"High\",\"Buffering_Risk\":\"Low\",\"Startup_Time\":\"Priority\"}",
            "playerType":            "xp",
            "responseFormatVersion": "1.0.0",
            "titleId":               titleId,
         },
      },
      "vodXrayMetadataRequest": map[string]interface{}{
         "xrayDeviceClass":  "television",
         "xrayPlaybackMode": "playback",
         "xrayToken":        "XRAY_REIGN_3PLR_2025_V1",
      },
   }

   body, err := json.Marshal(payload)
   if err != nil {
      return "", err
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
      return "", err
   }
   defer resp.Body.Close()

   if resp.StatusCode != http.StatusOK {
      return "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
   }

   var result struct {
      VodPlaylistedPlaybackUrls struct {
         Result struct {
            PlaybackUrls struct {
               IntraTitlePlaylist []struct {
                  Type string `json:"type"`
                  Urls []struct {
                     Url string `json:"url"`
                  } `json:"urls"`
               } `json:"intraTitlePlaylist"`
            } `json:"playbackUrls"`
         } `json:"result"`
         Error struct {
            Message string `json:"message"`
         } `json:"error"`
      } `json:"vodPlaylistedPlaybackUrls"`
   }

   if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
      return "", err
   }

   if result.VodPlaylistedPlaybackUrls.Error.Message != "" {
      return "", fmt.Errorf("API error: %s", result.VodPlaylistedPlaybackUrls.Error.Message)
   }

   for _, playlist := range result.VodPlaylistedPlaybackUrls.Result.PlaybackUrls.IntraTitlePlaylist {
      if playlist.Type == "Main" && len(playlist.Urls) > 0 {
         return playlist.Urls[0].Url, nil
      }
   }

   return "", fmt.Errorf("mpd url not found in response")
}
