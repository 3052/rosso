package rakuten

import (
   "encoding/json"
   "net/url"

   "41.neocities.org/maya"
)

type EpisodeStreaming struct {
   Id          string              `json:"id"`
   StreamInfos []EpisodeStreamInfo `json:"stream_infos"`
}

type EpisodeStreamInfo struct {
   Player     string `json:"player"`
   LicenseUrl string `json:"license_url"`
   Url        string `json:"url"`
}

func GetEpisodeStreaming(episodeItem *Episode) (*EpisodeStreaming, error) {
   payload := map[string]string{
      "audio_language":              "FRA",
      "audio_quality":               "2.0",
      "classification_id":           "23",
      "content_id":                  episodeItem.Id,
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

   target := &url.URL{
      Scheme: "https",
      Host:   "gizmo.rakuten.tv",
      Path:   "/v3/avod/streamings",
   }

   resp, err := maya.Post(target, nil, body)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()

   var wrapper struct {
      Data *EpisodeStreaming `json:"data"`
   }
   if err := json.NewDecoder(resp.Body).Decode(&wrapper); err != nil {
      return nil, err
   }
   return wrapper.Data, nil
}
