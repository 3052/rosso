package rakuten

import (
   "41.neocities.org/maya"
   "encoding/json"
   "errors"
   "io"
   "net/http"
   "net/url"
   "strconv"
)

// Season fetches episodes for a specific season (GET).
func (c *Content) Season(seasonId string) (*Season, error) {
   urlData := url.URL{
      Scheme: "https",
      Host:   "gizmo.rakuten.tv",
      Path:   "/v3/seasons/" + seasonId,
      RawQuery: url.Values{
         "classification_id": {strconv.Itoa(c.ClassificationId)},
         "device_identifier": {DeviceId},
         "market_code":       {c.MarketCode},
      }.Encode(),
   }

   resp, err := http.Get(urlData.String())
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()

   if resp.StatusCode != http.StatusOK {
      return nil, errors.New(resp.Status)
   }

   var result struct {
      Data Season
   }
   if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
      return nil, err
   }
   return &result.Data, nil
}

// For TV Shows, 'id' should be the Episode ID.
// For Movies, 'id' is ignored (uses c.Id).
func (c *Content) FetchStreamInfo(id, audioLanguage string, playerData Player, quality VideoQuality) (*StreamInfo, error) {
   data := map[string]string{
      "audio_language":              audioLanguage,
      "audio_quality":               "2.0",
      "classification_id":           strconv.Itoa(c.ClassificationId),
      "device_identifier":           DeviceId,
      "device_serial":               "not implemented",
      "device_stream_video_quality": string(quality),
      "player":                      string(playerData),
      "subtitle_language":           "MIS",
      "video_type":                  "stream",
   }
   switch c.Type {
   case "tv_shows":
      data["content_id"] = id
      data["content_type"] = "episodes"
   case "movies":
      data["content_id"] = c.Id
      data["content_type"] = "movies"
   }
   body, err := json.Marshal(data)
   if err != nil {
      return nil, err
   }
   resp, err := maya.Post(
      &url.URL{
         Scheme: "https",
         Host:   "gizmo.rakuten.tv",
         Path:   "/v3/avod/streamings",
      },
      map[string]string{"content-type": "application/json"},
      body,
   )
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   var result struct {
      Data struct {
         StreamInfos []StreamInfo `json:"stream_infos"`
      }
   }
   err = json.NewDecoder(resp.Body).Decode(&result)
   if err != nil {
      return nil, err
   }
   return &result.Data.StreamInfos[0], nil
}

func (s *StreamInfo) FetchWidevine(body []byte) ([]byte, error) {
   target, err := url.Parse(s.LicenseUrl)
   if err != nil {
      return nil, err
   }
   resp, err := maya.Post(
      target, map[string]string{"content-type": "application/x-protobuf"}, body,
   )
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   return io.ReadAll(resp.Body)
}
