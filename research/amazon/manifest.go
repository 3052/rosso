// START OF FILE manifest.go
package amazon

import (
   "encoding/json"
   "fmt"
   "io"
   "net/http"
   "net/url"
   "regexp"
   "sort"
   "strings"
)

// ManifestResponse represents the JSON returned by the GetPlaybackResources endpoint.
type ManifestResponse struct {
   AudioVideoUrls struct {
      AvCdnUrlSets []AvCdnUrlSet `json:"avCdnUrlSets"`
   } `json:"audioVideoUrls"`
   ErrorsByResource map[string]struct {
      ErrorCode string `json:"errorCode"`
      Message   string `json:"message"`
   } `json:"errorsByResource"`
   ReturnedTitleRendition struct {
      ContentId           string                 `json:"contentId"`
      SelectedEntitlement map[string]interface{} `json:"selectedEntitlement"`
   } `json:"returnedTitleRendition"`
}

type AvCdnUrlSet struct {
   Cdn            string `json:"cdn"`
   CdnWeightsRank int    `json:"cdnWeightsRank"`
   AvUrlInfoList  []struct {
      Url string `json:"url"`
   } `json:"avUrlInfoList"`
}

// PlaybackOptions holds the customizable options for the manifest request.
type PlaybackOptions struct {
   VideoQuality string // SD, HD, UHD
   VideoCodec   string // H264, H265
   BitrateMode  string // CVBR, CBR, CVBR,CBR
   HDRFormat    string // None, Hdr10, DolbyVision
   IsPrimeVideo bool
}

func DefaultPlaybackOptions() PlaybackOptions {
   return PlaybackOptions{
      VideoQuality: "HD",
      VideoCodec:   "H264",
      BitrateMode:  "CVBR,CBR",
      HDRFormat:    "None",
      IsPrimeVideo: false,
   }
}

// GetPlaybackResources fetches the manifest metadata from Amazon.
func GetPlaybackResources(
   client *http.Client,
   endpoint string, // e.g. "https://atv-ps.amazon.com/cdp/catalog/GetPlaybackResources"
   accessToken string,
   asin string,
   marketplaceID string, // e.g. "ATVPDKIKX0DER" for US
   device map[string]string,
   opts PlaybackOptions,
) (*ManifestResponse, error) {

   reqURL, err := url.Parse(endpoint)
   if err != nil {
      return nil, err
   }

   gascEnabled := "false"
   if opts.IsPrimeVideo {
      gascEnabled = "true"
   }

   q := reqURL.Query()
   q.Set("asin", asin)
   q.Set("consumptionType", "Streaming")
   q.Set("desiredResources", "PlaybackUrls,AudioVideoUrls,CatalogMetadata,ForcedNarratives,SubtitlePresets,SubtitleUrls,TransitionTimecodes,TrickplayUrls,CuepointPlaylist,XRayMetadata,PlaybackSettings")
   q.Set("deviceID", device["device_serial"])
   q.Set("deviceTypeID", device["device_type"])
   q.Set("firmware", "1")
   q.Set("gascEnabled", gascEnabled)
   q.Set("marketplaceID", marketplaceID)
   q.Set("resourceUsage", "CacheResources")
   q.Set("videoMaterialType", "Feature")
   q.Set("playerType", "html5")
   q.Set("clientId", "f22dbddb-ef2c-48c5-8876-bed0d47594fd") // Browser client ID from python
   q.Set("deviceDrmOverride", "CENC")
   q.Set("deviceStreamingTechnologyOverride", "DASH")
   q.Set("deviceProtocolOverride", "Https")
   q.Set("deviceVideoCodecOverride", opts.VideoCodec)
   q.Set("deviceBitrateAdaptationsOverride", opts.BitrateMode)
   q.Set("deviceVideoQualityOverride", opts.VideoQuality)
   q.Set("deviceHdrFormatsOverride", opts.HDRFormat)
   q.Set("supportedDRMKeyScheme", "DUAL_KEY")
   q.Set("liveManifestType", "live,accumulating")
   q.Set("titleDecorationScheme", "primary-content")
   q.Set("subtitleFormat", "TTMLv2")
   q.Set("languageFeature", "MLFv2")
   q.Set("uxLocale", "en_US")
   q.Set("xrayDeviceClass", "normal")
   q.Set("xrayPlaybackMode", "playback")
   q.Set("xrayToken", "XRAY_WEB_2020_V1")
   q.Set("playbackSettingsFormatVersion", "1.0.0")
   q.Set("playerAttributes", `{"frameRate": "HFR"}`)

   reqURL.RawQuery = q.Encode()

   req, err := http.NewRequest(http.MethodGet, reqURL.String(), nil)
   if err != nil {
      return nil, err
   }
   req.Header.Set("Authorization", "Bearer "+accessToken)

   resp, err := client.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()

   bodyBytes, err := io.ReadAll(resp.Body)
   if err != nil {
      return nil, err
   }

   var manifestResp ManifestResponse
   if err := json.Unmarshal(bodyBytes, &manifestResp); err != nil {
      return nil, fmt.Errorf("failed to decode response: %v\nBody: %s", err, string(bodyBytes))
   }

   // Check for rights/entitlement exceptions
   if _, hasException := manifestResp.ReturnedTitleRendition.SelectedEntitlement["rightsException"]; hasException {
      return nil, fmt.Errorf("entitlement error: the profile used does not have the rights to this title")
   }

   // Check for Playback errors
   if pbErr, ok := manifestResp.ErrorsByResource["PlaybackUrls"]; ok && pbErr.ErrorCode != "PRS.NoRights.NotOwned" {
      return nil, fmt.Errorf("playback URLs error: %s [%s]", pbErr.Message, pbErr.ErrorCode)
   }

   return &manifestResp, nil
}

// GetBestMPDURL sorts available CDN manifests by rank and returns the highest priority URL.
func GetBestMPDURL(manifest *ManifestResponse) (string, error) {
   sets := manifest.AudioVideoUrls.AvCdnUrlSets
   if len(sets) == 0 {
      return "", fmt.Errorf("no DASH manifests available")
   }

   // Sort ascending by CdnWeightsRank (lower number = higher priority / rank 1 is best)
   sort.Slice(sets, func(i, j int) bool {
      return sets[i].CdnWeightsRank < sets[j].CdnWeightsRank
   })

   if len(sets[0].AvUrlInfoList) == 0 {
      return "", fmt.Errorf("CDN url list is empty")
   }

   return sets[0].AvUrlInfoList[0].Url, nil
}

var cleanRegex = regexp.MustCompile(`^(https?://.*/)d.?/.*~/(.*)$`)

// CleanMPDURL translates the Python MPD URL cleaning logic safely using Go's net/url.
func CleanMPDURL(mpdURL string) string {
   // Try regex match first: removes the proxying segments
   matches := cleanRegex.FindStringSubmatch(mpdURL)
   if len(matches) == 3 {
      return matches[1] + matches[2]
   }

   // Fallback logic equivalent to: re.split(r"(?i)(/)", mpd_url)[:5] + re.split(r"(?i)(/)", mpd_url)[9:]
   // Essentially removing the 1st and 2nd directories from the URL path.
   u, err := url.Parse(mpdURL)
   if err == nil {
      parts := strings.Split(u.Path, "/")
      // parts[0] is "" (before the first slash)
      // parts[1] is 1st dir
      // parts[2] is 2nd dir
      // parts[3:] is the rest of the path
      if len(parts) > 3 {
         u.Path = "/" + strings.Join(parts[3:], "/")
         return u.String()
      }
   }

   return mpdURL
}
