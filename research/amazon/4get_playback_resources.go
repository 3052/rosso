// get_playback_resources.go
package amazon

import (
   "bytes"
   "encoding/json"
   "fmt"
   "io/ioutil"
   "net/http"
   "net/url"
)

// GetPlaybackResources requests the playback authorization and manifest URLs for a video.
func GetPlaybackResources(client *http.Client) ([]byte, error) {
   baseURL := "https://atv-ps.amazon.com/cdp/catalog/GetPlaybackResources"

   // Construct the query parameters
   q := url.Values{}
   q.Set("deviceTypeID", "AOAGZA014O5RE")
   q.Set("firmware", "1")
   q.Set("asin", "amzn1.dv.gti.a5ecdc15-befb-4ce3-952f-daedae5d34d7")
   q.Set("consumptionType", "Streaming")
   q.Set("desiredResources", "PlaybackUrls,SubtitleUrls,ForcedNarratives,TrickplayUrls,TransitionTimecodes,PlaybackSettings,CatalogMetadata,XRayMetadata")
   q.Set("resourceUsage", "CacheResources")
   q.Set("videoMaterialType", "Trailer")
   q.Set("deviceStreamingTechnologyOverride", "DASH")

   reqUrl := baseURL + "?" + q.Encode()

   // The payload containing the highly-specific playbackEnvelope JWT needed for authorization
   payloadMap := map[string]any{
      "globalParameters": map[string]any{
         "deviceCapabilityFamily": "WebPlayer",
         "capabilityDiscriminators": map[string]any{
            "operatingSystem": map[string]any{
               "name":    "Windows",
               "version": "10.0",
            },
            "middleware": map[string]any{
               "name":    "Firefox64",
               "version": "140.0",
            },
            "nativeApplication": map[string]any{
               "name":    "Firefox64",
               "version": "140.0",
            },
            "hfrControlMode": "Legacy",
            "displayResolution": map[string]any{
               "height": 1080,
               "width":  1920,
            },
         },
      },
      "auditPingsRequest":                 map[string]any{},
      "widevineServiceCertificateRequest": map[string]any{},
      "playbackDataRequest":               map[string]any{},
      "timedTextUrlsRequest": map[string]any{
         "supportedTimedTextFormats": []string{"TTMLv2", "DFXP"},
      },
      "trickplayUrlsRequest":       map[string]any{},
      "transitionTimecodesRequest": map[string]any{},
      "vodPlaylistedPlaybackUrlsRequest": map[string]any{
         "device": map[string]any{
            "hdcpLevel":                      "1.4",
            "maxVideoResolution":             "1080p",
            "supportedStreamingTechnologies": []string{"DASH"},
            "streamingTechnologies": map[string]any{
               "DASH": map[string]any{
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
      "ads": map[string]any{
         "sitePageUrl": "https://www.amazon.com/gp/video/detail/B075RND57T?ref_=nav_custrec_signin",
         "gdpr": map[string]any{
            "enabled":    false,
            "consentMap": map[string]any{},
         },
         "mainContentResumeOffsetHintMillis": 8849,
         "playerContractVersion":             1,
      },
      "playbackCustomizations": map[string]any{},
      "playbackSettingsRequest": map[string]any{
         "firmware":              "UNKNOWN",
         "playerType":            "xp",
         "responseFormatVersion": "1.0.0",
         "titleId":               "amzn1.dv.gti.a5ecdc15-befb-4ce3-952f-daedae5d34d7",
      },
      "vodXrayMetadataRequest": map[string]any{
         "xrayDeviceClass":  "normal",
         "xrayPlaybackMode": "playback",
         "xrayToken":        "XRAY_WEB_2023_V2",
      },
   }

   payloadJSON, err := json.Marshal(payloadMap)
   if err != nil {
      return nil, fmt.Errorf("error marshaling payload: %w", err)
   }

   req, err := http.NewRequest("POST", reqUrl, bytes.NewBuffer(payloadJSON))
   if err != nil {
      return nil, fmt.Errorf("error creating request: %w", err)
   }

   req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:140.0) Gecko/20100101 Firefox/140.0")
   req.Header.Set("Accept", "*/*")
   req.Header.Set("Accept-Language", "en-US,en;q=0.5")
   req.Header.Set("Accept-Encoding", "identity")
   req.Header.Set("Content-Type", "text/plain")
   req.Header.Set("Origin", "https://www.amazon.com")
   req.Header.Set("Referer", "https://www.amazon.com/")
   req.Header.Set("Sec-Fetch-Dest", "empty")
   req.Header.Set("Sec-Fetch-Mode", "cors")
   req.Header.Set("Sec-Fetch-Site", "same-site")
   req.Header.Set("Priority", "u=4")

   resp, err := client.Do(req)
   if err != nil {
      return nil, fmt.Errorf("error executing request: %w", err)
   }
   defer resp.Body.Close()

   if resp.StatusCode != http.StatusOK {
      return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
   }

   return ioutil.ReadAll(resp.Body)
}
