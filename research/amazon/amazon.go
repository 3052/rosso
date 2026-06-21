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

// DeviceProfile holds the specific capabilities and identities for a target device.
type DeviceProfile struct {
   DeviceID      string
   DRMType       string   // "Widevine" or "PlayReady"
   DRMKeyScheme  string   // "DualKey" or "SingleKey"
   HDCPLevel     string   // "1.4", "2.2", "2.3"
   MaxResolution string   // "480p", "1080p", "2160p"
   Codecs        []string // "H264", "H265"
   HDRFormats    []string // "None", "HDR10", "DolbyVision"
   APIHost       string   // e.g., "atv-ps.primevideo.com" (regional)
   AuthBearer    string   // Required for authorization
}

var (
   // WidevineL3Profile requests standard 1080p H264 content.
   WidevineL3Profile = DeviceProfile{
      DRMType:       "Widevine",
      DRMKeyScheme:  "DualKey",
      HDCPLevel:     "1.4",
      MaxResolution: "1080p",
      Codecs:        []string{"H264"},
      HDRFormats:    []string{"None"},
      APIHost:       "atv-ps.primevideo.com",
   }

   // PlayReadySL2000Profile requests standard 1080p H264 content via PlayReady.
   PlayReadySL2000Profile = DeviceProfile{
      DRMType:       "PlayReady",
      DRMKeyScheme:  "SingleKey",
      HDCPLevel:     "1.4",
      MaxResolution: "1080p",
      Codecs:        []string{"H264"},
      HDRFormats:    []string{"None"},
      APIHost:       "atv-ps.primevideo.com",
   }

   // PlayReadySL3000Profile requests 4K HDR/DV H265 content via PlayReady.
   PlayReadySL3000Profile = DeviceProfile{
      DRMType:       "PlayReady",
      DRMKeyScheme:  "SingleKey",
      HDCPLevel:     "2.3",
      MaxResolution: "2160p",
      Codecs:        []string{"H265"},
      HDRFormats:    []string{"HDR10", "DolbyVision"},
      APIHost:       "atv-ps.primevideo.com",
   }
)

// ManifestResponse defines the structure to extract the MPD URL and handoff token.
type ManifestResponse struct {
   Sessionization struct {
      SessionHandoffToken string `json:"sessionHandoffToken"`
   } `json:"sessionization"`
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
// Returns the MPD URL, Session Handoff Token, and any error encountered.
func (c *Client) GetManifest(p DeviceProfile, titleID, marketplaceID, envelope string) (string, string, error) {
   u := url.URL{
      Scheme: "https",
      Host:   p.APIHost,
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
               "DASH": map[string]any{
                  "bitrateAdaptations":  []string{"CBR", "CVBR"},
                  "codecs":              p.Codecs,
                  "drmKeyScheme":        p.DRMKeyScheme,
                  "drmType":             p.DRMType,
                  "dynamicRangeFormats": p.HDRFormats,
                  // IMPORTANT: Forces the smaller SegmentBase MPD format
                  "fragmentRepresentations": []string{"ByteOffsetRange", "SeparateFile"},
                  "segmentInfoType":         "Base",
                  "stitchType":              "MultiPeriod",
               },
            },
            "displayWidth":  3840,
            "displayHeight": 2160,
         },
      },
   }

   bodyBytes, err := json.Marshal(payload)
   if err != nil {
      return "", "", err
   }

   req, err := http.NewRequest("POST", u.String(), bytes.NewBuffer(bodyBytes))
   if err != nil {
      return "", "", err
   }

   req.Header.Set("Content-Type", "application/json")
   if p.AuthBearer != "" {
      req.Header.Set("Authorization", "Bearer "+p.AuthBearer)
   }

   resp, err := c.HTTPClient.Do(req)
   if err != nil {
      return "", "", err
   }
   defer resp.Body.Close()

   if resp.StatusCode != http.StatusOK {
      body, _ := io.ReadAll(resp.Body)
      return "", "", fmt.Errorf("bad status %d: %s", resp.StatusCode, string(body))
   }

   var result ManifestResponse
   if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
      return "", "", err
   }

   handoffToken := result.Sessionization.SessionHandoffToken
   var mpdURL string

   playlists := result.VodPlaylistedPlaybackUrls.Result.PlaybackUrls.IntraTitlePlaylist
   if len(playlists) > 0 && len(playlists[0].Urls) > 0 {
      mpdURL = playlists[0].Urls[0].URL
   }

   if mpdURL == "" || handoffToken == "" {
      return "", "", fmt.Errorf("failed to extract MPD or Handoff Token from response")
   }

   return mpdURL, handoffToken, nil
}

// GetLicense submits the CDM challenge and retrieves the base64 encoded license.
// `challenge` expects raw bytes for Widevine or raw XML/SOAP bytes for PlayReady.
func (c *Client) GetLicense(p DeviceProfile, titleID, marketplaceID, envelope, handoffToken string, challenge []byte) (string, error) {
   endpoint := "/playback/drm-vod/GetWidevineLicense"
   if p.DRMType == "PlayReady" {
      endpoint = "/playback/drm-vod/GetPlayReadyLicense"
   }

   u := url.URL{
      Scheme: "https",
      Host:   p.APIHost,
      Path:   endpoint,
   }
   q := u.Query()
   q.Set("deviceID", p.DeviceID)
   q.Set("deviceTypeID", "A3NM0WFSU3DLT5") // Hardcoded per requirements
   q.Set("marketplaceID", marketplaceID)
   q.Set("titleId", titleID)
   u.RawQuery = q.Encode()

   payload := map[string]any{
      "playbackEnvelope":    envelope,
      "sessionHandoffToken": handoffToken,
      "licenseChallenge":    base64.StdEncoding.EncodeToString(challenge),
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
