package roku

import (
   "encoding/json"
   "net/url"

   "41.neocities.org/maya"
)

type MediaPlayback struct {
   Url         string `json:"url"`
   Drm         Drm    `json:"drm"`
   MediaFormat string `json:"mediaFormat"`
   TraceId     string `json:"traceId"`
}

type Drm struct {
   Widevine Widevine `json:"widevine"`
}

type Widevine struct {
   LicenseServer string `json:"licenseServer"`
}

type PlaybackRequest struct {
   MediaFormat string `json:"mediaFormat"`
   ProviderId  string `json:"providerId"`
   RokuId      string `json:"rokuId"`
}

func FetchMediaPlayback(token ContentToken, targetId string) (*MediaPlayback, error) {
   endpoint := &url.URL{
      Scheme: "https",
      Host:   "googletv.web.roku.com",
      Path:   "/api/v3/playback",
   }

   headers := map[string]string{
      "content-type":         "application/json",
      "user-agent":           "trc-googletv; production; 0",
      "x-roku-content-token": string(token),
   }

   payload := PlaybackRequest{
      MediaFormat: "DASH",
      ProviderId:  "rokuavod",
      RokuId:      targetId,
   }

   body, err := json.Marshal(payload)
   if err != nil {
      return nil, err
   }

   resp, err := maya.Post(endpoint, headers, body)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()

   var playback MediaPlayback
   if err := json.NewDecoder(resp.Body).Decode(&playback); err != nil {
      return nil, err
   }

   return &playback, nil
}
