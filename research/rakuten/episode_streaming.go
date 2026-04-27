package rakuten

import (
   "encoding/json"
   "net/url"

   "41.neocities.org/maya"
)

type EpisodeStreamingInfo struct {
   Url  string        `json:"url"`
   Wrid StreamingUuid `json:"wrid"`
}

type EpisodeStreaming struct {
   StreamInfos []EpisodeStreamingInfo `json:"stream_infos"`
}

type EpisodeStreamingPayload struct {
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

func CreateEpisodeStreaming(session *UserSession, episode *SeasonEpisode) (*EpisodeStreaming, error) {
   endpoint := &url.URL{
      Scheme: "https",
      Host:   "gizmo.rakuten.tv",
      Path:   "/v3/avod/streamings",
   }

   payload := EpisodeStreamingPayload{
      AudioLanguage:            session.Profile.AudioLanguage.Id,
      AudioQuality:             "2.0",
      ClassificationId:         session.Profile.Classification.NumericalId,
      ContentId:                episode.Id,
      ContentType:              "episodes",
      DeviceIdentifier:         "atvui40",
      DeviceSerial:             "not implemented",
      DeviceStreamVideoQuality: "UHD",
      Player:                   "atvui40:DASH-CENC:PR",
      SubtitleLanguage:         session.Profile.SubtitlesLanguage.Id,
      VideoType:                "stream",
   }

   body, err := json.Marshal(payload)
   if err != nil {
      return nil, err
   }

   resp, err := maya.Post(endpoint, nil, body)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()

   var wrapper struct {
      Data EpisodeStreaming `json:"data"`
   }
   if err := json.NewDecoder(resp.Body).Decode(&wrapper); err != nil {
      return nil, err
   }

   return &wrapper.Data, nil
}
