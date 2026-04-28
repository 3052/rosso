package rakuten

import (
   "41.neocities.org/maya"
   "encoding/json"
   "errors"
   "io"
   "net/url"
   "strconv"
   "strings"
)

type TvShow struct {
   Id      string   `json:"id"`
   Title   string   `json:"title"`
   Seasons []Season `json:"seasons"`
}

func FetchTvShow(tvShowId string, userClassification Classification, targetMarket Market) (*TvShow, error) {
   target := &url.URL{
      Scheme: "https",
      Host:   "gizmo.rakuten.tv",
      Path:   "/v3/tv_shows/" + tvShowId,
   }

   query := url.Values{}
   query.Set("classification_id", strconv.Itoa(userClassification.NumericalId))
   query.Set("device_identifier", "atvui40")
   query.Set("market_code", targetMarket.Code)
   target.RawQuery = query.Encode()

   resp, err := maya.Get(target, nil)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()

   var apiResp struct {
      Data TvShow `json:"data"`
   }

   if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
      return nil, err
   }

   return &apiResp.Data, nil
}

func (show *TvShow) String() string {
   var builder strings.Builder
   for index, currentSeason := range show.Seasons {
      if index >= 1 {
         builder.WriteByte('\n')
      }
      builder.WriteString("season id = ")
      builder.WriteString(currentSeason.Id)
   }
   return builder.String()
}

func (info *StreamInfo) FetchLicense(challenge []byte) ([]byte, error) {
   target, err := url.Parse(info.LicenseUrl)
   if err != nil {
      return nil, err
   }

   resp, err := maya.Post(target, nil, challenge)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()

   return io.ReadAll(resp.Body)
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

func FetchMovieStreaming(contentId string, userClassification Classification, audioLanguage string) (*StreamInfo, error) {
   return fetchStreaming(contentId, "movies", userClassification, audioLanguage)
}

func FetchEpisodeStreaming(contentId string, userClassification Classification, audioLanguage string) (*StreamInfo, error) {
   return fetchStreaming(contentId, "episodes", userClassification, audioLanguage)
}

func fetchStreaming(contentId string, contentType string, userClassification Classification, audioLanguage string) (*StreamInfo, error) {
   target := &url.URL{
      Scheme: "https",
      Host:   "gizmo.rakuten.tv",
      Path:   "/v3/avod/streamings",
   }

   payload := StreamingRequest{
      AudioLanguage:            audioLanguage,
      AudioQuality:             "2.0",
      ClassificationId:         userClassification.NumericalId,
      ContentId:                contentId,
      ContentType:              contentType,
      DeviceIdentifier:         "atvui40",
      DeviceSerial:             "not implemented",
      DeviceStreamVideoQuality: "UHD",
      Player:                   "atvui40:DASH-CENC:PR",
      SubtitleLanguage:         "MIS",
      VideoType:                "stream",
   }

   body, err := json.Marshal(payload)
   if err != nil {
      return nil, err
   }

   headers := map[string]string{
      "content-type": "application/json",
   }

   resp, err := maya.Post(target, headers, body)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()

   var apiResp struct {
      Data struct {
         StreamInfos []StreamInfo `json:"stream_infos"`
      } `json:"data"`
   }

   if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
      return nil, err
   }

   for _, info := range apiResp.Data.StreamInfos {
      return &info, nil
   }

   return nil, errors.New("no stream infos found")
}

func (info *StreamInfo) GetManifest() (*url.URL, error) {
   return url.Parse(info.Url)
}
