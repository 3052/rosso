// file: video.go
package kanopy

import (
   "encoding/json"
   "io"
   "net/url"

   "41.neocities.org/maya"
)

type Video struct {
   VideoID int    `json:"videoId"`
   Title   string `json:"title"`
}

type VideoResponse struct {
   Type  string `json:"type"`
   Video Video  `json:"video"`
}

func GetVideo(jwt string, alias string) (*VideoResponse, error) {
   targetUrl := &url.URL{
      Scheme: "https",
      Host:   "www.kanopy.com",
      Path:   "/kapi/videos/alias/" + alias,
   }

   headers := map[string]string{
      "x-version":     "!/!/!/!",
      "authorization": "Bearer " + jwt,
   }

   resp, err := maya.Get(targetUrl, headers)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()

   bodyBytes, err := io.ReadAll(resp.Body)
   if err != nil {
      return nil, err
   }

   var videoResp VideoResponse
   if err := json.Unmarshal(bodyBytes, &videoResp); err != nil {
      return nil, err
   }

   return &videoResp, nil
}
