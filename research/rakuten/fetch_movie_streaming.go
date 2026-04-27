package rakuten

import (
   "encoding/json"
   "net/url"

   "41.neocities.org/maya"
)

type Streamings struct {
   StreamInfos []StreamInfo `json:"stream_infos"`
}

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

func FetchMovieStreaming(film *Movie, rating *Classification, audio *Language) (*Streamings, error) {
   target := &url.URL{
      Scheme: "https",
      Host:   "gizmo.rakuten.tv",
      Path:   "/v3/avod/streamings",
   }

   payload := StreamingRequest{
      AudioLanguage:            audio.Id,
      AudioQuality:             "2.0",
      ClassificationId:         rating.NumericalId,
      ContentId:                film.Id,
      ContentType:              "movies",
      DeviceIdentifier:         "atvui40",
      DeviceSerial:             "not implemented",
      DeviceStreamVideoQuality: "UHD",
      Player:                   "atvui40:DASH-CENC:PR",
      SubtitleLanguage:         "MIS",
      VideoType:                "stream",
   }

   reqBytes, err := json.Marshal(payload)
   if err != nil {
      return nil, err
   }

   headers := map[string]string{
      "content-type": "application/json",
   }

   resp, err := maya.Post(target, headers, reqBytes)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()

   var respWrapper struct {
      Data Streamings `json:"data"`
   }

   if err := json.NewDecoder(resp.Body).Decode(&respWrapper); err != nil {
      return nil, err
   }

   return &respWrapper.Data, nil
}
