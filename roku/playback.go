package roku

import (
   "encoding/json"
   "net/url"

   "41.neocities.org/maya"
)

type PlaybackConfig struct {
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

func GetPlaybackConfig(token *AccountToken, rokuId string) (*PlaybackConfig, error) {
   target := &url.URL{
      Scheme: "https",
      Host:   "googletv.web.roku.com",
      Path:   "/api/v3/playback",
   }
   headers := map[string]string{
      "content-type":         "application/json",
      "user-agent":           "trc-googletv; production; 0",
      "x-roku-content-token": token.AuthToken,
   }

   reqBody, err := json.Marshal(map[string]string{
      "mediaFormat": "DASH",
      "providerId":  "rokuavod",
      "rokuId":      rokuId,
   })
   if err != nil {
      return nil, err
   }

   resp, err := maya.Post(target, headers, reqBody)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()

   var config PlaybackConfig
   if err := json.NewDecoder(resp.Body).Decode(&config); err != nil {
      return nil, err
   }
   return &config, nil
}
