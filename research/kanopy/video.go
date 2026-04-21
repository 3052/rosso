// video.go
package kanopy

import (
   "encoding/json"
   "io"
   "net/url"

   "41.neocities.org/maya"
)

type VideoResponse struct {
   Type  string `json:"type"`
   Video Video  `json:"video"`
}

type Video struct {
   VideoId         int    `json:"videoId"`
   Title           string `json:"title"`
   DescriptionHtml string `json:"descriptionHtml"`
   ProductionYear  int    `json:"productionYear"`
   IsKids          bool   `json:"isKids"`
   DurationSeconds int    `json:"durationSeconds"`
   IsFree          bool   `json:"isFree"`
}

func GetVideo(alias string, jwt string) (*VideoResponse, error) {
   target := &url.URL{
      Scheme: "https",
      Host:   "www.kanopy.com",
      Path:   "/kapi/videos/alias/" + alias,
   }

   headers := map[string]string{
      "x-version":     "!/!/!/!",
      "authorization": "Bearer " + jwt,
   }

   resp, err := maya.Get(target, headers)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()

   respBytes, err := io.ReadAll(resp.Body)
   if err != nil {
      return nil, err
   }

   var videoResp VideoResponse
   if err := json.Unmarshal(respBytes, &videoResp); err != nil {
      return nil, err
   }

   return &videoResp, nil
}
