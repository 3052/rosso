package rakuten

import (
   "encoding/json"
   "net/url"

   "41.neocities.org/maya"
)

func CreateStreamingInfoHd(episode *EpisodeItem) (*StreamingInfo, error) {
   location := &url.URL{
      Scheme: "https",
      Host:   "gizmo.rakuten.tv",
      Path:   "/v3/avod/streamings",
   }

   payload := map[string]string{
      "audio_language":              "ENG",
      "audio_quality":               "2.0",
      "classification_id":           "18",
      "content_id":                  episode.Id,
      "content_type":                "episodes",
      "device_identifier":           "atvui40",
      "device_serial":               "not implemented",
      "device_stream_video_quality": "HD",
      "player":                      "atvui40:DASH-CENC:WVM",
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

   resp, err := maya.Post(location, headers, body)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()

   var response struct {
      Data *StreamingInfo `json:"data"`
   }
   if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
      return nil, err
   }
   return response.Data, nil
}
