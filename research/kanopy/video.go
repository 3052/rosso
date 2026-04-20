package kanopy

import (
   "encoding/json"
   "net/url"

   "41.neocities.org/maya"
)

type Video struct {
   VideoID int    `json:"videoId"`
   Title   string `json:"title"`
   Alias   string `json:"alias"`
}

type VideoResponse struct {
   Type  string `json:"type"`
   Video *Video `json:"video"`
}

func GetVideo(alias, jwt string) (*VideoResponse, error) {
   videoURL, err := url.Parse("https://www.kanopy.com/kapi/videos/alias/" + alias)
   if err != nil {
      return nil, err
   }

   headers := map[string]string{
      "x-version":     "!/!/!/!",
      "authorization": "Bearer " + jwt,
   }

   resp, err := maya.Get(videoURL, headers)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()

   var videoResp VideoResponse
   if err := json.NewDecoder(resp.Body).Decode(&videoResp); err != nil {
      return nil, err
   }

   return &videoResp, nil
}
