// playback.go
package amazon

import (
   "bytes"
   "context"
   "crypto/rand"
   "encoding/json"
   "fmt"
   "net/http"
   "strconv"
   "time"
)

type PlaybackParams struct {
   BaseURL          string
   DeviceID         string
   DeviceTypeID     string
   GascEnabled      bool
   MarketplaceID    string
   TitleID          string
   DeviceToken      string
   PlaybackEnvelope string
   Quality          string // SD, HD, UHD
   VideoCodec       string // H264, H265
   BitrateMode      string // CVBR, CBR, CVBR+CBR
   HDR              string // SDR, HDR10, DV
   IsPlayReady      bool
   PlayerType       string // html5, xp
}

// GenerateNerid generates a Network Edge Request ID
func GenerateNerid(e int) string {
   const base64Chars = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/"

   timestamp := time.Now().UnixMilli()
   tsChars := make([]byte, 7)
   for i := 0; i < 7; i++ {
      tsChars[i] = base64Chars[timestamp%64]
      timestamp /= 64
   }

   for i, j := 0, len(tsChars)-1; i < j; i, j = i+1, j-1 {
      tsChars[i], tsChars[j] = tsChars[j], tsChars[i]
   }

   randBytes := make([]byte, 15)
   _, _ = rand.Read(randBytes)
   randPart := make([]byte, 15)
   for i := 0; i < 15; i++ {
      randPart[i] = base64Chars[int(randBytes[i])%64]
   }

   suffix := fmt.Sprintf("%02d", e%100)
   return string(tsChars) + string(randPart) + suffix
}

// GetVodPlaybackResources makes a request to /playback/prs/GetVodPlaybackResources
func GetVodPlaybackResources(ctx context.Context, client *http.Client, params PlaybackParams) (map[string]interface{}, error) {
   urlPath := fmt.Sprintf("https://%s/playback/prs/GetVodPlaybackResources", params.BaseURL)

   hdrMap := map[string]string{
      "SDR":   "None",
      "HDR10": "Hdr10",
      "DV":    "DolbyVision",
   }
   hdrFormat := "None"
   if val, ok := hdrMap[params.HDR]; ok {
      hdrFormat = val
   }

   var bitrateAdaptations []string
   if params.BitrateMode == "CVBR+CBR" || params.BitrateMode == "CVBR,CBR" {
      bitrateAdaptations = []string{"CVBR", "CBR"}
   } else {
      bitrateAdaptations = []string{params.BitrateMode}
   }

   var globalParams map[string]interface{}
   var vodPlaybackUrlsRequest map[string]interface{}

   if params.DeviceToken == "" {
      globalParams = map[string]interface{}{
         "deviceCapabilityFamily": "WebPlayer",
         "capabilityDiscriminators": map[string]interface{}{
            "operatingSystem":   map[string]string{"name": "Windows", "version": "10.0"},
            "middleware":        map[string]string{"name": "EdgeNext", "version": "136.0.0.0"},
            "nativeApplication": map[string]string{"name": "EdgeNext", "version": "136.0.0.0"},
            "hfrControlMode":    "Legacy",
            "displayResolution": map[string]int{"height": 2304, "width": 4096},
         },
      }

      if params.PlaybackEnvelope != "" {
         globalParams["playbackEnvelope"] = params.PlaybackEnvelope
      }

      drmType := "Widevine"
      drmKeyScheme := "DualKey"
      if params.IsPlayReady {
         drmType = "PlayReady"
         drmKeyScheme = "SingleKey"
      }

      hdcpLevel := "1.4"
      maxRes := "1080p"
      if params.Quality == "UHD" {
         hdcpLevel = "2.2"
         maxRes = "2160p"
      } else if params.Quality == "SD" {
         maxRes = "480p"
      }

      vodPlaybackUrlsRequest = map[string]interface{}{
         "device": map[string]interface{}{
            "hdcpLevel":                      hdcpLevel,
            "maxVideoResolution":             maxRes,
            "supportedStreamingTechnologies": []string{"DASH"},
            "streamingTechnologies": map[string]interface{}{
               "DASH": map[string]interface{}{
                  "bitrateAdaptations":       bitrateAdaptations,
                  "codecs":                   []string{params.VideoCodec},
                  "drmKeyScheme":             drmKeyScheme,
                  "drmType":                  drmType,
                  "dynamicRangeFormats":      hdrFormat,
                  "fragmentRepresentations":  []string{"ByteOffsetRange", "SeparateFile"},
                  "frameRates":               []string{"Standard"},
                  "segmentInfoType":          "Base",
                  "timedTextRepresentations": []string{"NotInManifestNorStream", "SeparateStreamInManifest"},
                  "trickplayRepresentations": []string{"NotInManifestNorStream"},
                  "variableAspectRatio":      "supported",
               },
            },
            "displayWidth":  4096,
            "displayHeight": 2304,
         },
         "ads": map[string]interface{}{
            "sitePageUrl": "",
            "gdpr": map[string]interface{}{
               "enabled":    "false",
               "consentMap": map[string]interface{}{},
            },
         },
         "playbackCustomizations": map[string]interface{}{},
         "playbackSettingsRequest": map[string]interface{}{
            "firmware":              "UNKNOWN",
            "playerType":            params.PlayerType,
            "responseFormatVersion": "1.0.0",
            "titleId":               params.TitleID,
         },
      }
   } else {
      globalParams = map[string]interface{}{
         "deviceCapabilityFamily": "AndroidPlayer",
         "capabilityDiscriminators": map[string]interface{}{
            "discriminators": map[string]interface{}{
               "software": map[string]interface{}{},
               "version":  1,
            },
         },
      }

      if params.PlaybackEnvelope != "" {
         globalParams["playbackEnvelope"] = params.PlaybackEnvelope
      }

      drmType := "Widevine"
      if params.IsPlayReady {
         drmType = "PlayReady"
      }

      techProfile := map[string]interface{}{
         "fragmentRepresentations":               []string{"ByteOffsetRange", "SeparateFile"},
         "manifestThinningToSupportedResolution": "Forbidden",
         "segmentInfoType":                       "List",
         "timedTextRepresentations":              []string{"BurnedIn", "NotInManifestNorStream", "SeparateStreamInManifest"},
         "trickplayRepresentations":              []string{"NotInManifestNorStream"},
         "variableAspectRatio":                   "supported",
         "vastTimelineType":                      "Absolute",
         "bitrateAdaptations":                    bitrateAdaptations,
         "codecs":                                []string{params.VideoCodec},
         "drmKeyScheme":                          "SingleKey",
         "drmStrength":                           "L40",
         "drmType":                               drmType,
         "dynamicRangeFormats":                   []string{hdrFormat},
         "frameRates":                            []string{"Standard"},
      }

      techProfileSmooth := map[string]interface{}{}
      for k, v := range techProfile {
         techProfileSmooth[k] = v
      }
      techProfileSmooth["drmType"] = "PlayReady"

      vodPlaybackUrlsRequest = map[string]interface{}{
         "ads": map[string]interface{}{},
         "device": map[string]interface{}{
            "displayBasedVending": "supported",
            "displayHeight":       2304,
            "displayWidth":        4096,
            "streamingTechnologies": map[string]interface{}{
               "DASH":            techProfile,
               "SmoothStreaming": techProfileSmooth,
            },
            "acceptedCreativeApis":           []string{},
            "category":                       "Tv",
            "hdcpLevel":                      "2.2",
            "maxVideoResolution":             "2160p",
            "platform":                       "Android",
            "supportedStreamingTechnologies": []string{"DASH", "SmoothStreaming"},
         },
         "playbackCustomizations": map[string]interface{}{},
         "playbackSettingsRequest": map[string]interface{}{
            "firmware":              "UNKNOWN",
            "playerType":            params.PlayerType,
            "responseFormatVersion": "1.0.0",
            "titleId":               params.TitleID,
         },
      }
   }

   auditPingsRequest := map[string]interface{}{}
   if params.DeviceToken != "" {
      auditPingsRequest = map[string]interface{}{
         "device": map[string]string{
            "category": "Tv",
            "platform": "Android",
         },
      }
   }

   payload := map[string]interface{}{
      "globalParameters":           globalParams,
      "auditPingsRequest":          auditPingsRequest,
      "playbackDataRequest":        map[string]interface{}{},
      "timedTextUrlsRequest":       map[string]interface{}{"supportedTimedTextFormats": []string{"TTMLv2", "DFXP"}},
      "trickplayUrlsRequest":       map[string]interface{}{},
      "transitionTimecodesRequest": map[string]interface{}{},
      "vodPlaybackUrlsRequest":     vodPlaybackUrlsRequest,
      "vodXrayMetadataRequest": map[string]string{
         "xrayDeviceClass":  "normal",
         "xrayPlaybackMode": "playback",
         "xrayToken":        "XRAY_WEB_2023_V2",
      },
   }

   body, err := json.Marshal(payload)
   if err != nil {
      return nil, err
   }

   req, err := http.NewRequestWithContext(ctx, http.MethodPost, urlPath, bytes.NewReader(body))
   if err != nil {
      return nil, err
   }

   q := req.URL.Query()
   q.Add("deviceID", params.DeviceID)
   q.Add("deviceTypeID", params.DeviceTypeID)
   q.Add("gascEnabled", strconv.FormatBool(params.GascEnabled))
   q.Add("marketplaceID", params.MarketplaceID)
   q.Add("uxLocale", "en_EN")
   q.Add("firmware", "1")
   q.Add("titleId", params.TitleID)
   q.Add("nerid", GenerateNerid(0))
   req.URL.RawQuery = q.Encode()

   req.Header.Set("Content-Type", "application/json")
   if params.DeviceToken != "" {
      req.Header.Set("Authorization", "Bearer "+params.DeviceToken)
   }

   resp, err := client.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()

   var result map[string]interface{}
   if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
      return nil, err
   }

   return result, nil
}
