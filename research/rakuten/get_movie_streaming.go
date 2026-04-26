package rakuten

import (
   "encoding/json"
   "net/url"

   "41.neocities.org/maya"
)

type MovieStreaming struct {
   Id          string            `json:"id"`
   StreamInfos []MovieStreamInfo `json:"stream_infos"`
}

type MovieStreamInfo struct {
   Player     string `json:"player"`
   LicenseUrl string `json:"license_url"`
   Url        string `json:"url"`
}

func GetMovieStreaming(film *Movie) (*MovieStreaming, error) {
   payload := map[string]string{
      "audio_language":              "ENG",
      "audio_quality":               "2.0",
      "classification_id":           "41",
      "content_id":                  film.Id,
      "content_type":                "movies",
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
      Data *MovieStreaming `json:"data"`
   }
   if err := json.NewDecoder(resp.Body).Decode(&wrapper); err != nil {
      return nil, err
   }
   return wrapper.Data, nil
}
