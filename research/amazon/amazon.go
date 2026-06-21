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
   DeviceTypeID  string
   Family        string   // "LivingRoomPlayer", "WebPlayer", or "AndroidPlayer"
   DRMType       string   // "Widevine" or "PlayReady"
   DRMKeyScheme  string   // "DualKey" or "SingleKey"
   HDCPLevel     string   // "1.4", "2.2", "2.3"
   MaxResolution string   // "1080p", "2160p"
   Codecs        []string // "H264", "H265"
   HDRFormats    []string // "None", "HDR10", "HDR10Plus", "DolbyVision"
   APIHost       string   // e.g., "ab8mt4dd97et.na.api.amazonvideo.com" or "atv-ps.primevideo.com"
   AuthBearer    string   // Required if using APIHost (LivingRoom/Android)
   Cookies       string   // Required if using Web APIHost (WebPlayer)
}

var (
   // WidevineL3Profile requests standard 1080p H264 content.
   WidevineL3Profile = DeviceProfile{
      DeviceTypeID:  "A2SNKIF736WF4T", // AndroidTV / LivingRoom
      Family:        "LivingRoomPlayer",
      DRMType:       "Widevine",
      DRMKeyScheme:  "DualKey",
      HDCPLevel:     "1.4",
      MaxResolution: "1080p",
      Codecs:        []string{"H264"},
      HDRFormats:    []string{"None"},
      APIHost:       "ab8mt4dd97et.na.api.amazonvideo.com",
   }

   // PlayReadySL2000Profile requests standard 1080p H264 content via PlayReady.
   PlayReadySL2000Profile = DeviceProfile{
      DeviceTypeID:  "AOAGZA014O5RE", // Edge / Web
      Family:        "WebPlayer",
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
      DeviceTypeID:  "A2SNKIF736WF4T", // AndroidTV / LivingRoom
      Family:        "LivingRoomPlayer",
      DRMType:       "PlayReady",
      DRMKeyScheme:  "SingleKey",
      HDCPLevel:     "2.3",
      MaxResolution: "2160p",
      Codecs:        []string{"H265"},
      HDRFormats:    []string{"HDR10", "HDR10Plus", "DolbyVision"},
      APIHost:       "ab8mt4dd97et.na.api.amazonvideo.com",
   }
)

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
func (c *Client) GetManifest(p DeviceProfile, titleID, marketplaceID, envelope string) (mpdURL string, handoffToken string, err error) {
   u := url.URL{
      Scheme: "https",
      Host:   p.APIHost,
      Path:   "/playback/prs/GetVodPlaybackResources",
   }
   q := u.Query()
   q.Set("deviceID", p.DeviceID)
   q.Set("deviceTypeID", p.DeviceTypeID)
   q.Set("marketplaceID", marketplaceID)
   q.Set("titleId", titleID)
   q.Set("uxLocale", "en_US")
   q.Set("firmware", "1")
   u.RawQuery = q.Encode()

   // Build the payload optimizing for SegmentBase
   payload := map[string]any{
      "globalParameters": map[string]any{
         "deviceCapabilityFamily": p.Family,
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

   bodyBytes, _ := json.Marshal(payload)
   req, err := http.NewRequest("POST", u.String(), bytes.NewBuffer(bodyBytes))
   if err != nil {
      return "", "", err
   }

   req.Header.Set("Content-Type", "application/json")
   if p.AuthBearer != "" {
      req.Header.Set("Authorization", "Bearer "+p.AuthBearer)
   }
   if p.Cookies != "" {
      req.Header.Set("Cookie", p.Cookies)
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

   var result map[string]any
   if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
      return "", "", err
   }

   if sess, ok := result["sessionization"].(map[string]any); ok {
      handoffToken, _ = sess["sessionHandoffToken"].(string)
   }

   // Deep navigate to find MPD URL
   if vppu, ok := result["vodPlaylistedPlaybackUrls"].(map[string]any); ok {
      if res, ok := vppu["result"].(map[string]any); ok {
         if pu, ok := res["playbackUrls"].(map[string]any); ok {
            if itp, ok := pu["intraTitlePlaylist"].([]any); ok && len(itp) > 0 {
               if firstPlaylist, ok := itp[0].(map[string]any); ok {
                  if urls, ok := firstPlaylist["urls"].([]any); ok && len(urls) > 0 {
                     if firstURL, ok := urls[0].(map[string]any); ok {
                        mpdURL, _ = firstURL["url"].(string)
                     }
                  }
               }
            }
         }
      }
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
   q.Set("deviceTypeID", p.DeviceTypeID)
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

   bodyBytes, _ := json.Marshal(payload)
   req, err := http.NewRequest("POST", u.String(), bytes.NewBuffer(bodyBytes))
   if err != nil {
      return "", err
   }

   req.Header.Set("Content-Type", "application/json")
   if p.AuthBearer != "" {
      req.Header.Set("Authorization", "Bearer "+p.AuthBearer)
   }
   if p.Cookies != "" {
      req.Header.Set("Cookie", p.Cookies)
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

   var result map[string]any
   if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
      return "", err
   }

   var license string
   if p.DRMType == "Widevine" {
      if wl, ok := result["widevineLicense"].(map[string]any); ok {
         license, _ = wl["license"].(string)
      }
   } else {
      if pl, ok := result["playReadyLicense"].(map[string]any); ok {
         license, _ = pl["license"].(string)
      }
   }

   if license == "" {
      return "", fmt.Errorf("could not find license string in JSON response")
   }

   return license, nil
}
