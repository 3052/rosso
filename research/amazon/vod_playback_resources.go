// post_get_vod_playback_resources.go
package amazon

import (
   "encoding/json"
   "fmt"
   "io"
   "net/http"
   "net/url"
   "strings"
)

func GetVodPlaybackResources(s *Session) (string, error) {
   if s.DeviceID == "" {
      s.DeviceID = GenerateUUID()
   }

   targetURL := fmt.Sprintf("https://atv-ps.amazon.com/playback/prs/GetVodPlaybackResources?deviceID=%s&deviceTypeID=AOAGZA014O5RE&gascEnabled=false&marketplaceID=ATVPDKIKX0DER&uxLocale=en_US&firmware=1&titleId=%s",
      s.DeviceID,
      url.QueryEscape(s.TargetTitleID),
   )

   payloadMap := map[string]interface{}{
      "globalParameters": map[string]interface{}{
         "deviceCapabilityFamily": "WebPlayer",
         "playbackEnvelope":       s.PlaybackEnvelope,
         "capabilityDiscriminators": map[string]interface{}{
            "operatingSystem":   map[string]string{"name": "Windows", "version": "10.0"},
            "middleware":        map[string]string{"name": "Firefox64", "version": "140.0"},
            "nativeApplication": map[string]string{"name": "Firefox64", "version": "140.0"},
            "hfrControlMode":    "Legacy",
            "displayResolution": map[string]int{"height": 1080, "width": 1920},
         },
      },
      "vodPlaylistedPlaybackUrlsRequest": map[string]interface{}{
         "device": map[string]interface{}{
            "hdcpLevel":                      "1.4",
            "maxVideoResolution":             "1080p",
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
                  "frameRates":                       []string{"Standard", "High"},
                  "stitchType":                       "MultiPeriod",
                  "segmentInfoType":                  "Base",
                  "timedTextRepresentations":         []string{"NotInManifestNorStream", "SeparateStreamInManifest"},
                  "trickplayRepresentations":         []string{"NotInManifestNorStream"},
                  "variableAspectRatio":              "supported",
               },
            },
            "displayWidth":  1920,
            "displayHeight": 1080,
         },
      },
   }

   payloadBytes, err := json.Marshal(payloadMap)
   if err != nil {
      return "", err
   }

   req, err := http.NewRequest("POST", targetURL, strings.NewReader(string(payloadBytes)))
   if err != nil {
      return "", err
   }

   req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:140.0) Gecko/20100101 Firefox/140.0")
   req.Header.Set("Content-Type", "text/plain")
   req.Header.Set("Accept", "*/*")
   req.Header.Set("Origin", "https://www.amazon.com")
   req.Header.Set("Referer", "https://www.amazon.com/")

   resp, err := s.Client.Do(req)
   if err != nil {
      return "", err
   }
   defer resp.Body.Close()

   body, err := io.ReadAll(resp.Body)
   if err != nil {
      return "", err
   }

   return string(body), nil
}
