package amazon

import (
   "bytes"
   "encoding/base64"
   "encoding/json"
   "fmt"
   "io"
   "net/http"
   "net/url"
)

const defaultAPIHost = "atv-ps.primevideo.com"

// DeviceProfile holds the specific capabilities and identities to test against the API.
type DeviceProfile struct {
   DeviceID      string
   DRMType       string   // "Widevine" or "PlayReady"
   DRMKeyScheme  string   // Optional: "DualKey", "SingleKey", or leave empty to omit
   HDCPLevel     string   // e.g. "1.4", "2.2", "2.3"
   MaxResolution string   // e.g. "480p", "720p", "1080p", "1440p", "2160p"
   HDRFormats    []string // e.g. "None", "HDR10", "DolbyVision"
   AuthBearer    string   // Required for authorization
}

// ManifestResponse defines the structure to extract the MPD URL.
type ManifestResponse struct {
   VodPlaylistedPlaybackUrls struct {
      Result struct {
         PlaybackUrls struct {
            IntraTitlePlaylist []struct {
               Urls []struct {
                  URL string `json:"url"`
               } `json:"urls"`
            } `json:"intraTitlePlaylist"`
         } `json:"playbackUrls"`
      } `json:"result"`
   } `json:"vodPlaylistedPlaybackUrls"`
}

// LicenseResponse defines the structure to extract the base64 license.
type LicenseResponse struct {
   WidevineLicense struct {
      License string `json:"license"`
   } `json:"widevineLicense"`
   PlayReadyLicense struct {
      License string `json:"license"`
   } `json:"playReadyLicense"`
}

// Client handles the communication with Amazon APIs.
type Client struct {
   HTTPClient *http.Client
}

// NewClient creates a new Amazon API client.
func NewClient(httpClient *http.Client) *Client {
   if httpClient == nil {
      httpClient = http.DefaultClient
   }
   return &Client{HTTPClient: httpClient}
}

// GetManifest requests the Playback Resources.
// It forces a SegmentBase MPD (~5MB instead of ~30MB) by manipulating the payload.
// Returns the MPD URL and any error encountered.
func (c *Client) GetManifest(p DeviceProfile, titleID, marketplaceID, envelope string) (string, error) {
   u := url.URL{
      Scheme: "https",
      Host:   defaultAPIHost,
      Path:   "/playback/prs/GetVodPlaybackResources",
   }
   q := u.Query()
   q.Set("deviceID", p.DeviceID)
   q.Set("deviceTypeID", "A3NM0WFSU3DLT5") // Hardcoded per requirements
   q.Set("marketplaceID", marketplaceID)
   q.Set("titleId", titleID)
   q.Set("uxLocale", "en_US")
   q.Set("firmware", "1")
   u.RawQuery = q.Encode()

   dashSettings := map[string]any{
      "bitrateAdaptations":  []string{"CBR", "CVBR"},
      "codecs":              []string{"H265"}, // Hardcoded per requirements
      "drmType":             p.DRMType,
      "dynamicRangeFormats": p.HDRFormats,
      // IMPORTANT: Forces the smaller SegmentBase MPD format by only allowing ByteOffsetRange
      "fragmentRepresentations": []string{"ByteOffsetRange"},
      "segmentInfoType":         "Base",
      "stitchType":              "MultiPeriod",
   }

   // Only append drmKeyScheme if it is explicitly provided
   if p.DRMKeyScheme != "" {
      dashSettings["drmKeyScheme"] = p.DRMKeyScheme
   }

   var width, height int
   switch p.MaxResolution {
   case "2160p":
      width, height = 3840, 2160
   case "1440p":
      width, height = 2560, 1440
   case "1080p":
      width, height = 1920, 1080
   case "720p":
      width, height = 1280, 720
   case "480p":
      width, height = 854, 480
   default:
      width, height = 1920, 1080
   }

   // Build the payload optimizing for SegmentBase
   payload := map[string]any{
      "globalParameters": map[string]any{
         "deviceCapabilityFamily": "LivingRoomPlayer", // Hardcoded per requirements
         "playbackEnvelope":       envelope,
      },
      "vodPlaylistedPlaybackUrlsRequest": map[string]any{
         "device": map[string]any{
            "hdcpLevel":                      p.HDCPLevel,
            "maxVideoResolution":             p.MaxResolution,
            "supportedStreamingTechnologies": []string{"DASH"},
            "streamingTechnologies": map[string]any{
               "DASH": dashSettings,
            },
            "displayWidth":  width,
            "displayHeight": height,
         },
      },
   }

   bodyBytes, err := json.Marshal(payload)
   if err != nil {
      return "", err
   }

   req, err := http.NewRequest("POST", u.String(), bytes.NewBuffer(bodyBytes))
   if err != nil {
      return "", err
   }

   req.Header.Set("Content-Type", "application/json")
   if p.AuthBearer != "" {
      req.Header.Set("Authorization", "Bearer "+p.AuthBearer)
   }

   resp, err := c.HTTPClient.Do(req)
   if err != nil {
      return "", err
   }
   defer resp.Body.Close()

   if resp.StatusCode != http.StatusOK {
      body, _ := io.ReadAll(resp.Body)
      return "", fmt.Errorf("bad status %d: %s", resp.StatusCode, string(body))
   }

   var result ManifestResponse
   if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
      return "", err
   }

   var mpdURL string

   playlists := result.VodPlaylistedPlaybackUrls.Result.PlaybackUrls.IntraTitlePlaylist
   if len(playlists) > 0 && len(playlists[0].Urls) > 0 {
      mpdURL = playlists[0].Urls[0].URL
   }

   if mpdURL == "" {
      return "", fmt.Errorf("failed to extract MPD from response")
   }

   return mpdURL, nil
}

// GetLicense submits the CDM challenge and retrieves the base64 encoded license.
// `challenge` expects raw bytes for Widevine or raw XML/SOAP bytes for PlayReady.
func (c *Client) GetLicense(p DeviceProfile, titleID, marketplaceID, envelope string, challenge []byte) (string, error) {
   endpoint := "/playback/drm-vod/GetWidevineLicense"
   if p.DRMType == "PlayReady" {
      endpoint = "/playback/drm-vod/GetPlayReadyLicense"
   }

   u := url.URL{
      Scheme: "https",
      Host:   defaultAPIHost,
      Path:   endpoint,
   }
   q := u.Query()
   q.Set("deviceID", p.DeviceID)
   q.Set("deviceTypeID", "A3NM0WFSU3DLT5") // Hardcoded per requirements
   q.Set("marketplaceID", marketplaceID)
   q.Set("titleId", titleID)
   u.RawQuery = q.Encode()

   payload := map[string]any{
      "playbackEnvelope": envelope,
      "licenseChallenge": base64.StdEncoding.EncodeToString(challenge),
   }

   if p.DRMType == "Widevine" {
      payload["includeHdcpTestKey"] = true
   } else if p.DRMType == "PlayReady" {
      payload["packagingFormat"] = "MPEG_DASH"
   }

   bodyBytes, err := json.Marshal(payload)
   if err != nil {
      return "", err
   }

   req, err := http.NewRequest("POST", u.String(), bytes.NewBuffer(bodyBytes))
   if err != nil {
      return "", err
   }

   req.Header.Set("Content-Type", "application/json")
   if p.AuthBearer != "" {
      req.Header.Set("Authorization", "Bearer "+p.AuthBearer)
   }

   resp, err := c.HTTPClient.Do(req)
   if err != nil {
      return "", err
   }
   defer resp.Body.Close()

   if resp.StatusCode != http.StatusOK {
      body, _ := io.ReadAll(resp.Body)
      return "", fmt.Errorf("bad status %d: %s", resp.StatusCode, string(body))
   }

   var result LicenseResponse
   if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
      return "", err
   }

   license := result.WidevineLicense.License
   if p.DRMType == "PlayReady" {
      license = result.PlayReadyLicense.License
   }

   if license == "" {
      return "", fmt.Errorf("could not find license string in JSON response")
   }

   return license, nil
}
