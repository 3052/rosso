package rakuten

import (
   "encoding/json"
   "net/url"

   "41.neocities.org/maya"
)

type StreamInfo struct {
   LicenseUrl string `json:"license_url"`
   Wrid       string `json:"wrid"`
}

type StreamingsResponse struct {
   StreamInfos []*StreamInfo `json:"stream_infos"`
}

type StreamingPayload struct {
   AudioLanguage            string      `json:"audio_language"`
   AudioQuality             string      `json:"audio_quality"`
   ClassificationId         int         `json:"classification_id"`
   ContentId                ContentId   `json:"content_id"`
   ContentType              ContentType `json:"content_type"`
   DeviceIdentifier         string      `json:"device_identifier"`
   DeviceSerial             string      `json:"device_serial"`
   DeviceStreamVideoQuality string      `json:"device_stream_video_quality"`
   Player                   string      `json:"player"`
   SubtitleLanguage         string      `json:"subtitle_language"`
   VideoType                string      `json:"video_type"`
}

func CreateStreamings(targetId ContentId, targetType ContentType, sessionResp *SessionResponse) (*StreamingsResponse, error) {
   endpoint := &url.URL{
      Scheme: "https",
      Host:   "gizmo.rakuten.tv",
      Path:   "/v3/avod/streamings",
   }

   payload := &StreamingPayload{
      AudioLanguage:            sessionResp.Market.DefaultAudioLanguage.Abbr,
      AudioQuality:             "2.0",
      ClassificationId:         sessionResp.Profile.Classification.NumericalId,
      ContentId:                targetId,
      ContentType:              targetType,
      DeviceIdentifier:         "atvui40",
      DeviceSerial:             "not implemented",
      DeviceStreamVideoQuality: "UHD",
      Player:                   "atvui40:DASH-CENC:PR",
      SubtitleLanguage:         "MIS",
      VideoType:                "stream",
   }

   bodyBytes, err := json.Marshal(payload)
   if err != nil {
      return nil, err
   }

   headers := map[string]string{
      "content-type": "application/json",
   }

   resp, err := maya.Post(endpoint, headers, bodyBytes)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()

   var wrapper struct {
      Data *StreamingsResponse `json:"data"`
   }
   if err := json.NewDecoder(resp.Body).Decode(&wrapper); err != nil {
      return nil, err
   }

   return wrapper.Data, nil
}
