package rakuten

import (
   "encoding/json"
   "net/url"

   "41.neocities.org/maya"
)

type MovieStreamData struct {
   StreamInfos []MovieStreamInfo `json:"stream_infos"`
}

type MovieStreamInfo struct {
   Url        string     `json:"url"`
   LicenseUrl LicenseUrl `json:"license_url"`
   Wrid       string     `json:"wrid"`
}

type MovieStreamPayload struct {
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

func GetMovieStreamData(movie *Movie, sessionData *SessionData) (*MovieStreamData, error) {
   endpoint := &url.URL{
      Scheme: "https",
      Host:   "gizmo.rakuten.tv",
      Path:   "/v3/avod/streamings",
   }

   movieStreamPayload := MovieStreamPayload{
      AudioLanguage:            sessionData.Profile.AudioLanguage.Id,
      AudioQuality:             "2.0",
      ClassificationId:         sessionData.Profile.Classification.Id,
      ContentId:                movie.Id,
      ContentType:              "movies",
      DeviceIdentifier:         "atvui40",
      DeviceSerial:             "not implemented",
      DeviceStreamVideoQuality: "UHD",
      Player:                   "atvui40:DASH-CENC:PR",
      SubtitleLanguage:         "MIS",
      VideoType:                "stream",
   }

   payloadData, err := json.Marshal(movieStreamPayload)
   if err != nil {
      return nil, err
   }

   resp, err := maya.Post(endpoint, nil, payloadData)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()

   var response struct {
      Data MovieStreamData `json:"data"`
   }
   if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
      return nil, err
   }

   return &response.Data, nil
}
