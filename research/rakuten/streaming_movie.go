package rakuten

import (
   "encoding/json"
   "net/url"

   "41.neocities.org/maya"
)

type MovieLicenseUuid string

type MovieStreamInfo struct {
   Wrid MovieLicenseUuid `json:"wrid"`
}

type MovieStreaming struct {
   StreamInfos []MovieStreamInfo `json:"stream_infos"`
}

func CreateMovieStreaming(contentId MovieContentId, classId ClassificationId, audioLang LanguageId) (*MovieStreaming, error) {
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

   headers := map[string]string{
      "content-type": "application/json",
   }

   resp, err := maya.Post(&endpoint, headers, body)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()

   var wrapper struct {
      Data MovieStreaming `json:"data"`
   }
   decoder := json.NewDecoder(resp.Body)
   if err := decoder.Decode(&wrapper); err != nil {
      return nil, err
   }

   return &wrapper.Data, nil
}
