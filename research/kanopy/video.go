// File: get_video.go
package kanopy

import (
   "encoding/json"
   "net/url"

   "41.neocities.org/maya"
)

type VideoResponse struct {
   Type  string `json:"type"`
   Video struct {
      VideoID     int    `json:"videoId"`
      Title       string `json:"title"`
      Description string `json:"descriptionHtml"`
      IsFree      bool   `json:"isFree"`
      Alias       string `json:"alias"`
   } `json:"video"`
}

func GetVideo(alias string, token string) (*VideoResponse, error) {
   reqURL, err := url.Parse("https://www.kanopy.com/kapi/videos/alias/" + alias)
   if err != nil {
      return nil, err
   }

   headers := map[string]string{
      "authorization": "Bearer " + token,
      "x-version":     "!/!/!/!",
      "user-agent":    "Go-http-client/2.0",
   }

   resp, err := maya.Get(reqURL, headers)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()

   var result VideoResponse
   if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
      return nil, err
   }

   return &result, nil
}
