// videos.go
package kanopy

import (
   "encoding/json"
   "net/url"

   "41.neocities.org/maya"
)

type Video struct {
   VideoID int `json:"videoId"`
}

type VideoResponse struct {
   Type  string `json:"type"`
   Video Video  `json:"video"`
}

func GetVideo(jwt string, alias string) (*VideoResponse, error) {
   targetUrl, err := url.Parse("https://www.kanopy.com/kapi/videos/alias/" + alias)
   if err != nil {
      return nil, err
   }

   headers := map[string]string{
      "authorization": "Bearer " + jwt,
      "x-version":     "!/!/!/!",
   }

   resp, err := maya.Get(targetUrl, headers)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()

   var videoResp VideoResponse
   err = json.NewDecoder(resp.Body).Decode(&videoResp)
   if err != nil {
      return nil, err
   }

   return &videoResp, nil
}
