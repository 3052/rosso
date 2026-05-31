// get_playback_resources.go
package amazon

import (
   "fmt"
   "io/ioutil"
   "net/http"
   "net/url"
   "strings"
)

// GetPlaybackResources requests the playback authorization and manifest URLs for a video.
func GetPlaybackResources(client *http.Client) ([]byte, error) {
   baseURL := "https://atv-ps.amazon.com/cdp/catalog/GetPlaybackResources"

   // Construct the query parameters
   q := url.Values{}
   q.Set("deviceID", "9e531fce-fd2a-4957-956f-7061e53be8b9")
   q.Set("deviceTypeID", "AOAGZA014O5RE")
   q.Set("gascEnabled", "false")
   q.Set("marketplaceID", "ATVPDKIKX0DER")
   q.Set("uxLocale", "en_US")
   q.Set("firmware", "1")
   q.Set("playerType", "xp")
   q.Set("operatingSystemName", "Windows")
   q.Set("operatingSystemVersion", "10.0")
   q.Set("deviceApplicationName", "Firefox64")
   q.Set("asin", "amzn1.dv.gti.a5ecdc15-befb-4ce3-952f-daedae5d34d7")
   q.Set("consumptionType", "Streaming")
   q.Set("desiredResources", "PlaybackUrls,SubtitleUrls,ForcedNarratives,TrickplayUrls,TransitionTimecodes,PlaybackSettings,CatalogMetadata,XRayMetadata")
   q.Set("resourceUsage", "CacheResources")
   q.Set("videoMaterialType", "Trailer")
   q.Set("clientId", "f22dbddb-ef2c-48c5-8876-bed0d47594fd")
   q.Set("userWatchSessionId", "7d7e6ce8-cc8b-4f45-92bd-918f19314836")
   q.Set("sitePageUrl", "https://www.amazon.com/gp/video/detail/B075RND57T")
   q.Set("displayWidth", "1920")
   q.Set("displayHeight", "1080")
   q.Set("supportsVariableAspectRatio", "true")
   q.Set("supportsEmbeddedTimedTextForVod", "true")
   q.Set("deviceProtocolOverride", "Https")
   q.Set("vodStreamSupportOverride", "Auxiliary")
   q.Set("deviceStreamingTechnologyOverride", "DASH")
   q.Set("deviceDrmOverride", "CENC")
   q.Set("deviceHdrFormatsOverride", "None")
   q.Set("deviceVideoCodecOverride", "H264")
   q.Set("deviceVideoQualityOverride", "HD")
   q.Set("deviceBitrateAdaptationsOverride", "CVBR,CBR")
   q.Set("supportsEmbeddedTrickplayForVod", "false")
   q.Set("audioTrackId", "all")
   q.Set("languageFeature", "MLFv2")
   q.Set("liveManifestType", "patternTemplate,accumulating,live")
   q.Set("supportedDRMKeyScheme", "DUAL_KEY")
   q.Set("supportsEmbeddedTrickplay", "true")
   q.Set("daiSupportsEmbeddedTrickplay", "true")
   q.Set("daiLiveManifestType", "patternTemplate,accumulating,live")
   q.Set("ssaiSegmentInfoSupport", "Base")
   q.Set("ssaiStitchType", "MultiPeriod")
   q.Set("gdprEnabled", "false")
   q.Set("subtitleFormat", "TTMLv2")
   q.Set("playbackSettingsFormatVersion", "1.0.0")
   q.Set("titleDecorationScheme", "primary-content")
   q.Set("xrayToken", "XRAY_WEB_2023_V2")
   q.Set("xrayPlaybackMode", "playback")
   q.Set("xrayDeviceClass", "normal")
   q.Set("nerid", "Z5/+NWUhffrJ5C7F+a4Koi00")
   reqUrl := baseURL + "?" + q.Encode()
   req, err := http.NewRequest("POST", reqUrl, strings.NewReader(data4))
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

const data4 = `
{
  "globalParameters": {
    "deviceCapabilityFamily": "WebPlayer",
    "capabilityDiscriminators": {
      "operatingSystem": {
        "name": "Windows",
        "version": "10.0"
      },
      "middleware": {
        "name": "Firefox64",
        "version": "140.0"
      },
      "nativeApplication": {
        "name": "Firefox64",
        "version": "140.0"
      },
      "hfrControlMode": "Legacy",
      "displayResolution": {
        "height": 1080,
        "width": 1920
      }
    }
  },
  "auditPingsRequest": {},
  "widevineServiceCertificateRequest": {},
  "playbackDataRequest": {},
  "timedTextUrlsRequest": {
    "supportedTimedTextFormats": [
      "TTMLv2",
      "DFXP"
    ]
  },
  "trickplayUrlsRequest": {},
  "transitionTimecodesRequest": {},
  "vodPlaylistedPlaybackUrlsRequest": {
    "device": {
      "hdcpLevel": "1.4",
      "maxVideoResolution": "1080p",
      "supportedStreamingTechnologies": [
        "DASH"
      ],
      "streamingTechnologies": {
        "DASH": {
          "bitrateAdaptations": [
            "CBR",
            "CVBR"
          ],
          "codecs": [
            "H264"
          ],
          "drmKeyScheme": "DualKey",
          "drmType": "Widevine",
          "dynamicRangeFormats": [
            "None"
          ],
          "edgeDeliveryAuthorizationSchemes": [
            "PVExchangeV1",
            "Transparent"
          ],
          "fragmentRepresentations": [
            "ByteOffsetRange",
            "SeparateFile"
          ],
          "frameRates": [
            "Standard",
            "High"
          ],
          "stitchType": "MultiPeriod",
          "segmentInfoType": "Base",
          "timedTextRepresentations": [
            "NotInManifestNorStream",
            "SeparateStreamInManifest"
          ],
          "trickplayRepresentations": [
            "NotInManifestNorStream"
          ],
          "variableAspectRatio": "supported"
        }
      },
      "displayWidth": 1920,
      "displayHeight": 1080
    },
    "ads": {
      "sitePageUrl": "https://www.amazon.com/gp/video/detail/B075RND57T?ref_=nav_custrec_signin",
      "gdpr": {
        "enabled": false,
        "consentMap": {}
      },
      "mainContentResumeOffsetHintMillis": 8849,
      "playerContractVersion": 1
    },
    "playbackCustomizations": {},
    "playbackSettingsRequest": {
      "firmware": "UNKNOWN",
      "playerType": "xp",
      "responseFormatVersion": "1.0.0",
      "titleId": "amzn1.dv.gti.a5ecdc15-befb-4ce3-952f-daedae5d34d7"
    }
  },
  "vodXrayMetadataRequest": {
    "xrayDeviceClass": "normal",
    "xrayPlaybackMode": "playback",
    "xrayToken": "XRAY_WEB_2023_V2"
  }
}
`
