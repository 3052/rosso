package roku

import (
   "encoding/json"
   "net/url"

   "41.neocities.org/maya"
)

type Playback struct {
   Url           string `json:"url"`
   MediaFormat   string `json:"mediaFormat"`
   AdPolicyName  string `json:"adPolicyName"`
   KidsDirected  bool   `json:"kidsDirected"`
   RokuNielsenId string `json:"rokuNielsenId"`
   TraceId       string `json:"traceId"`
   Drm           Drm    `json:"drm"`
}

type Drm struct {
   Widevine Widevine `json:"widevine"`
}

type Widevine struct {
   LicenseServer string `json:"licenseServer"`
}

type PlaybackPayload struct {
   MediaFormat string `json:"mediaFormat"`
   ProviderId  string `json:"providerId"`
   RokuId      string `json:"rokuId"`
}

func CreatePlayback(userToken ContentToken, providerId string, rokuId string) (*Playback, error) {
   endpoint := &url.URL{
      Scheme: "https",
      Host:   "googletv.web.roku.com",
      Path:   "/api/v3/playback",
   }

   headers := map[string]string{
      "user-agent":           "trc-googletv; production; 0",
      "x-roku-content-token": string(userToken),
      "content-type":         "application/json",
   }

   payload := PlaybackPayload{
      MediaFormat: "DASH",
      ProviderId:  providerId,
      RokuId:      rokuId,
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

   var playback Playback
   if err := json.NewDecoder(resp.Body).Decode(&playback); err != nil {
      return nil, err
   }

   return &playback, nil
}
