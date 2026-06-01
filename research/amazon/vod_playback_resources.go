// post_get_vod_playback_resources.go
package amazon

import (
   "encoding/json"
   "io"
   "net/http"
   "net/url"
   "strings"
)

func GetVodPlaybackResources(s *Session) (string, error) {
   payloadMap := map[string]interface{}{
      "globalParameters": map[string]interface{}{
         "deviceCapabilityFamily": "WebPlayer",
         "playbackEnvelope":       s.PlaybackEnvelope,
      },
      "vodPlaylistedPlaybackUrlsRequest": map[string]interface{}{
         "device": map[string]interface{}{
            "maxVideoResolution":             "1080p",
            "supportedStreamingTechnologies": []string{"DASH"},
            "streamingTechnologies": map[string]interface{}{
               "DASH": map[string]interface{}{
                  "bitrateAdaptations": []string{"CBR", "CVBR"},
                  "codecs":             []string{"H264"},
                  "drmType":            "Widevine",
                  // this is optional but changes the URL
                  "edgeDeliveryAuthorizationSchemes": []string{"PVExchangeV1", "Transparent"},
                  //"stitchType":                       "MultiPeriod",
                  //"segmentInfoType":                  "Base",
               },
            },
         },
      },
   }
   payloadBytes, err := json.Marshal(payloadMap)
   if err != nil {
      return "", err
   }
   u := &url.URL{
      Scheme: "https",
      Host:   "atv-ps.amazon.com",
      Path:   "/playback/prs/GetVodPlaybackResources",
   }
   q := u.Query()
   q.Set("deviceTypeID", "AOAGZA014O5RE")
   q.Set("deviceID", s.DeviceID)
   u.RawQuery = q.Encode()
   targetURL := u.String()
   req, err := http.NewRequest("POST", targetURL, strings.NewReader(string(payloadBytes)))
   if err != nil {
      return "", err
   }
   resp, err := s.Client.Do(req)
   if err != nil {
      return "", err
   }
   defer resp.Body.Close()

   body, err := io.ReadAll(resp.Body)
   if err != nil {
      return "", err
   }

   return string(body), nil
}
