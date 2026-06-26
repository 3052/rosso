package amazon

import (
   "bytes"
   "encoding/json"
   "fmt"
   "io"
   "net/http"
   "net/url"
   "strings"
)

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
   q.Set("deviceTypeID", defaultDeviceTypeID) // Centralized
   q.Set("marketplaceID", marketplaceID)
   q.Set("titleId", titleID)
   q.Set("uxLocale", "en_US")
   q.Set("firmware", "1")
   u.RawQuery = q.Encode()

   dashSettings := map[string]any{
      "bitrateAdaptations":  []string{p.BitrateAdaptation}, // Dynamic based on profile loop
      "codecs":              []string{p.VideoCodec},        // Dynamic based on profile loop
      "drmType":             p.DRMType,
      "dynamicRangeFormats": []string{p.HDRFormats}, // Wrap the single string into a slice
      // IMPORTANT: Forces the smaller SegmentBase MPD format by only allowing ByteOffsetRange
      "fragmentRepresentations": []string{"ByteOffsetRange"},
      "segmentInfoType":         "Base",
      "stitchType":              "MultiPeriod",
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
            "displayWidth":  3840,
            "displayHeight": 2160,
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
   if len(playlists) > 0 {
      // Require Akamai to avoid the 30MB Cloudfront/Amazon MPD bloat
      for _, u := range playlists[0].Urls {
         if strings.ToLower(u.CDN) == "akamai" {
            mpdURL = u.URL
            break
         }
      }
   }

   if mpdURL == "" {
      return "", fmt.Errorf("failed to extract Akamai MPD from response (it may be missing or Cloudfront only)")
   }

   cleanedURL, err := trimURLPath(mpdURL)
   if err != nil {
      return "", fmt.Errorf("failed to trim MPD URL path: %w", err)
   }

   return cleanedURL.String(), nil
}

// trimURLPath removes Amazon's restrictive path segments from the MPD URL.
func trimURLPath(rawUrl string) (*url.URL, error) {
   parsedURL, err := url.Parse(rawUrl)
   if err != nil {
      return nil, err
   }

   parts := strings.Split(parsedURL.Path, "/")

   // Handle "/dm/3$..." structure
   if len(parts) > 4 && parts[1] == "dm" && strings.HasPrefix(parts[2], "3$") {
      parsedURL.Path = "/" + strings.Join(parts[4:], "/")
      // Handle "/3$..." structure
   } else if len(parts) > 3 && strings.HasPrefix(parts[1], "3$") {
      parsedURL.Path = "/" + strings.Join(parts[3:], "/")
   }

   return parsedURL, nil
}

// ManifestResponse defines the structure to extract the MPD URL.
type ManifestResponse struct {
   VodPlaylistedPlaybackUrls struct {
      Result struct {
         PlaybackUrls struct {
            IntraTitlePlaylist []struct {
               Urls []struct {
                  URL string `json:"url"`
                  CDN string `json:"cdn"` // Used to explicitly require Akamai
               } `json:"urls"`
            } `json:"intraTitlePlaylist"`
         } `json:"playbackUrls"`
      } `json:"result"`
   } `json:"vodPlaylistedPlaybackUrls"`
}
