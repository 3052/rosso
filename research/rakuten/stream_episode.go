package rakuten

import (
   "encoding/json"
   "net/url"

   "41.neocities.org/maya"
)

type EpisodeStreamData struct {
   StreamInfos []EpisodeStreamInfo `json:"stream_infos"`
}

type EpisodeStreamInfo struct {
   Url        string     `json:"url"`
   LicenseUrl LicenseUrl `json:"license_url"`
   Wrid       string     `json:"wrid"`
}

type EpisodeStreamPayload struct {
   AudioLanguage            string `json:"audio_language"`
   AudioQuality             string `json:"audio_quality"`
   ClassificationId         string `json:"classification_id"`
   ContentId                string `json:"content_id"`
   ContentType              string `json:"content_type"`
   DeviceIdentifier         string `json:"device_identifier"`
   DeviceSerial             string `json:"device_serial"`
   DeviceStreamVideoQuality string `json:"device_stream_video_quality"`
   Player                   string `json:"player"`
   SubtitleLanguage         string `json:"subtitle_language"`
   VideoType                string `json:"video_type"`
}

func GetEpisodeStreamData(seasonEpisode *SeasonEpisode, sessionData *SessionData) (*EpisodeStreamData, error) {
   endpoint := &url.URL{
      Scheme: "https",
      Host:   "gizmo.rakuten.tv",
      Path:   "/v3/avod/streamings",
   }

   episodeStreamPayload := EpisodeStreamPayload{
      AudioLanguage:            sessionData.Profile.AudioLanguage.Id,
      AudioQuality:             "2.0",
      ClassificationId:         sessionData.Profile.Classification.Id,
      ContentId:                seasonEpisode.Id,
      ContentType:              "episodes",
      DeviceIdentifier:         "atvui40",
      DeviceSerial:             "not implemented",
      DeviceStreamVideoQuality: "UHD",
      Player:                   "atvui40:DASH-CENC:PR",
      SubtitleLanguage:         "MIS",
      VideoType:                "stream",
   }

   payloadData, err := json.Marshal(episodeStreamPayload)
   if err != nil {
      return nil, err
   }

   resp, err := maya.Post(endpoint, nil, payloadData)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()

   var response struct {
      Data EpisodeStreamData `json:"data"`
   }
   if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
      return nil, err
   }

   return &response.Data, nil
}
