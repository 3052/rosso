package rakuten

import (
   "encoding/json"
   "net/url"

   "41.neocities.org/maya"
)

func FetchEpisodeStreaming(chapter *Episode, rating *Classification, audio *Language, deviceIdentifier string) (*Streamings, error) {
   target := &url.URL{
      Scheme: "https",
      Host:   "gizmo.rakuten.tv",
      Path:   "/v3/avod/streamings",
   }

   payload := StreamingRequest{
      AudioLanguage:            audio.Id,
      AudioQuality:             "2.0",
      ClassificationId:         rating.NumericalId,
      ContentId:                chapter.Id,
      ContentType:              "episodes",
      DeviceIdentifier:         deviceIdentifier,
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
