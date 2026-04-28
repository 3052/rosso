// FILE: rakuten/fetch_streaming.go
package rakuten

import (
   "encoding/json"
   "net/url"

   "41.neocities.org/maya"
)

type StreamInfo struct {
   LicenseUrl string `json:"license_url"`
   Url        string `json:"url"`
}

type StreamingRequest struct {
   AudioLanguage            string `json:"audio_language"`
   AudioQuality             string `json:"audio_quality"`
   ClassificationId         int    `json:"classification_id"`
   ContentId                string `json:"content_id"`
   ContentType              string `json:"content_type"`
   DeviceIdentifier         string `json:"device_identifier"`
   DeviceSerial             string `json:"device_serial"`
   DeviceStreamVideoQuality string `json:"device_stream_video_quality"`
   Player                   string `json:"player"`
   SubtitleLanguage         string `json:"subtitle_language"`
   VideoType                string `json:"video_type"`
}

func FetchMovieStreaming(contentId string, rating *Classification, audioLanguage string) (*StreamInfo, error) {
   return fetchStreaming(contentId, "movies", rating, audioLanguage)
}

func FetchEpisodeStreaming(contentId string, rating *Classification, audioLanguage string) (*StreamInfo, error) {
   return fetchStreaming(contentId, "episodes", rating, audioLanguage)
}

func fetchStreaming(contentId string, contentType string, rating *Classification, audioLanguage string) (*StreamInfo, error) {
   target := &url.URL{
      Scheme: "https",
      Host:   "gizmo.rakuten.tv",
      Path:   "/v3/avod/streamings",
   }

   payload := StreamingRequest{
      AudioLanguage:            audioLanguage,
      AudioQuality:             "2.0",
      ClassificationId:         rating.NumericalId,
      ContentId:                contentId,
      ContentType:              contentType,
      DeviceIdentifier:         "atvui40",
      DeviceSerial:             "not implemented",
      DeviceStreamVideoQuality: "UHD",
      Player:                   "atvui40:DASH-CENC:PR",
      SubtitleLanguage:         "MIS",
      VideoType:                "stream",
   }

   body, err := json.Marshal(payload)
   if err != nil {
      return nil, err
   }

   headers := map[string]string{
      "content-type": "application/json",
   }

   resp, err := maya.Post(target, headers, body)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   var result struct {
      Data struct {
         StreamInfos []StreamInfo `json:"stream_infos"`
      }
   }
   if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
      return nil, err
   }
   return &result.Data.StreamInfos[0], nil
}
