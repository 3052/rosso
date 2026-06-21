package amazon

import (
   "bytes"
   "encoding/json"
   "fmt"
   "net/http"
   "net/url"
   "strings"
)

// GetVodPlaybackResources fetches the final MPD URL for playback.
// Pass "H264" or "H265" as the videoCodec.
// Pass "Widevine" or "PlayReady" as the drmType.
// Pass "CBR" or "CVBR" as the bitrateAdaptation.
// Pass "None", "DolbyVision", or "HDR10" as the dynamicRangeFormat.
func GetVodPlaybackResources(actorAccessToken, titleId, playbackEnvelope, videoCodec, drmType, bitrateAdaptation, dynamicRangeFormat string) (*PlaybackResource, error) {
   payload := map[string]any{
      "globalParameters": map[string]any{
         "playbackEnvelope":       playbackEnvelope,
         "deviceCapabilityFamily": "LivingRoomPlayer",
      },
      "vodPlaylistedPlaybackUrlsRequest": map[string]any{
         "device": map[string]any{
            "supportedStreamingTechnologies": []string{"DASH"},
            "streamingTechnologies": map[string]any{
               "DASH": map[string]any{
                  "bitrateAdaptations": []string{
                     bitrateAdaptation, // dynamically set ("CBR" or "CVBR")
                  },
                  "drmType": drmType, // dynamically set ("Widevine" or "PlayReady")
                  "dynamicRangeFormats": []string{
                     dynamicRangeFormat, // dynamically set ("None", "DolbyVision", or "HDR10")
                  },
                  "codecs": []string{
                     videoCodec, // dynamically set (e.g. "H264" or "H265")
                  },
               },
            },
            
            //FHD
            //"hdcpLevel": "2.1", //IIA
            
            //UHD
            "hdcpLevel": "2.3", // at least 2.2 is needed for UHD with hev1
            
            //"maxVideoResolution": "480p", // L3
            "maxVideoResolution": "2160p", // SL3000
         },
         "playbackSettingsRequest": map[string]any{
            "firmware": DeviceFirmware,
            "titleId":  titleId,
         },
      },
   }
   body, err := json.Marshal(payload)
   if err != nil {
      return nil, err
   }
   urlStr := "https://ab8mt4dd97et.na.api.amazonvideo.com/playback/prs/GetVodPlaybackResources"
   req, err := http.NewRequest("POST", urlStr, bytes.NewReader(body))
   if err != nil {
      return nil, err
   }
   query := url.Values{}
   query.Add("deviceID", DeviceID)
   query.Add("deviceTypeID", DeviceTypeID)
   req.URL.RawQuery = query.Encode()
   req.Header.Set("Authorization", "Bearer "+actorAccessToken)
   client := &http.Client{}
   resp, err := client.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()

   if resp.StatusCode != http.StatusOK {
      return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
   }

   var result struct {
      GlobalError struct {
         Code    string `json:"code"`
         Message string `json:"message"`
      } `json:"globalError"`
      VodPlaylistedPlaybackUrls struct {
         Result struct {
            PlaybackUrls struct {
               IntraTitlePlaylist []struct {
                  Type string             `json:"type"`
                  Urls []PlaybackResource `json:"urls"`
               } `json:"intraTitlePlaylist"`
            } `json:"playbackUrls"`
         } `json:"result"`
         Error struct {
            Message string `json:"message"`
         } `json:"error"`
      } `json:"vodPlaylistedPlaybackUrls"`
   }

   if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
      return nil, err
   }

   if result.GlobalError.Code != "" {
      return nil, fmt.Errorf("global API error: [%s] %s", result.GlobalError.Code, result.GlobalError.Message)
   }

   if result.VodPlaylistedPlaybackUrls.Error.Message != "" {
      return nil, fmt.Errorf("API error: %s", result.VodPlaylistedPlaybackUrls.Error.Message)
   }

   for _, playlist := range result.VodPlaylistedPlaybackUrls.Result.PlaybackUrls.IntraTitlePlaylist {
      if playlist.Type == "Main" && len(playlist.Urls) > 0 {
         res := playlist.Urls[0]
         return &res, nil
      }
   }

   return nil, fmt.Errorf("mpd url not found in response")
}

func (*PlaybackResource) CachePath() string {
   return "rosso/amazon/PlaybackResource"
}

type PlaybackResource struct {
   Url string
}

func (p *PlaybackResource) Clean() (*url.URL, error) {
   parsedURL, err := url.Parse(p.Url)
   if err != nil {
      return nil, err
   }
   parts := strings.Split(parsedURL.Path, "/")
   // Handle "/dm/3$..." structure
   if len(parts) > 4 && parts[1] == "dm" && strings.HasPrefix(parts[2], "3$") {
      // parts[0] = ""
      // parts[1] = "dm"
      // parts[2] = "3$..."
      // parts[3] = "iad_2"
      // parts[4:] = raw path
      parsedURL.Path = "/" + strings.Join(parts[4:], "/")
      // Handle "/3$..." structure
   } else if len(parts) > 3 && strings.HasPrefix(parts[1], "3$") {
      // parts[0] = ""
      // parts[1] = "3$..."
      // parts[2] = "iad_2"
      // parts[3:] = raw path
      parsedURL.Path = "/" + strings.Join(parts[3:], "/")
   }
   return parsedURL, nil
}
