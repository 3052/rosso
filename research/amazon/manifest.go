package amazon

import (
   "bytes"
   "encoding/json"
   "fmt"
   "io"
   "net/http"
   "net/url"
)

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

// GetManifest requests the Playback Resources.
// It forces a SegmentBase MPD (~5MB instead of ~30MB) by manipulating the payload.
// Returns the MPD URL and any error encountered.
func (c *Client) GetManifest(p DeviceProfile, titleID, marketplaceID, envelope string) (string, error) {
   if p.AuthBearer == "" {
      return "", fmt.Errorf("AuthBearer is required")
   }

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
   req.Header.Set("Authorization", "Bearer "+p.AuthBearer)

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
