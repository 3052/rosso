package rakuten

import (
   "encoding/json"
   "net/url"

   "41.neocities.org/maya"
)

type EpisodeLicenseUuid string

type EpisodeStreamInfo struct {
   Wrid EpisodeLicenseUuid `json:"wrid"`
}

type EpisodeStreaming struct {
   StreamInfos []EpisodeStreamInfo `json:"stream_infos"`
}

func CreateEpisodeStreaming(contentId EpisodeContentId, classId ClassificationId, audioLang LanguageId) (*EpisodeStreaming, error) {
   endpoint := url.URL{
      Scheme: "https",
      Host:   "gizmo.rakuten.tv",
      Path:   "/v3/avod/streamings",
   }

   payload := map[string]any{
      "audio_language":              string(audioLang),
      "audio_quality":               "2.0",
      "classification_id":           int(classId),
      "content_id":                  string(contentId),
      "content_type":                "episodes",
      "device_identifier":           "atvui40",
      "device_serial":               "not implemented",
      "device_stream_video_quality": "UHD",
      "player":                      "atvui40:DASH-CENC:PR",
      "subtitle_language":           "MIS",
      "video_type":                  "stream",
   }

   body, err := json.Marshal(payload)
   if err != nil {
      return nil, err
   }

   headers := map[string]string{
      "content-type": "application/json",
   }

   resp, err := maya.Post(&endpoint, headers, body)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()

   var wrapper struct {
      Data EpisodeStreaming `json:"data"`
   }
   decoder := json.NewDecoder(resp.Body)
   if err := decoder.Decode(&wrapper); err != nil {
      return nil, err
   }

   return &wrapper.Data, nil
}
