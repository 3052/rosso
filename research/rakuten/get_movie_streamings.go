package rakuten

import (
   "encoding/json"
   "net/url"

   "41.neocities.org/maya"
)

type MovieLicenseUuid string

type MovieStreamings struct {
   StreamInfos []MovieStreamInfo `json:"stream_infos"`
}

type MovieStreamInfo struct {
   Wrid MovieLicenseUuid `json:"wrid"`
   Url  string           `json:"url"`
}

func GetMovieStreamings(movie *CinemaMovie, audioLanguage string, classificationId int) (*MovieStreamings, error) {
   link := &url.URL{
      Scheme: "https",
      Host:   "gizmo.rakuten.tv",
      Path:   "/v3/avod/streamings",
   }

   payload := map[string]any{
      "audio_language":              audioLanguage,
      "audio_quality":               "2.0",
      "classification_id":           classificationId,
      "content_id":                  movie.Id,
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

   resp, err := maya.Post(link, headers, body)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()

   var respWrapper struct {
      Data MovieStreamings `json:"data"`
   }
   if err := json.NewDecoder(resp.Body).Decode(&respWrapper); err != nil {
      return nil, err
   }
   return &respWrapper.Data, nil
}
