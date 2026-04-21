package kanopy

import (
   "encoding/json"
   "net/url"

   "41.neocities.org/maya"
)

type Video struct {
   VideoId int    `json:"videoId"`
   Title   string `json:"title"`
   Alias   string `json:"alias"`
}

func GetVideo(alias, jwt string) (*Video, error) {
   videoUrl, err := url.Parse("https://www.kanopy.com/kapi/videos/alias/" + alias)
   if err != nil {
      return nil, err
   }

   headers := map[string]string{
      "x-version":     "!/!/!/!",
      "authorization": "Bearer " + jwt,
   }

   resp, err := maya.Get(videoUrl, headers)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   var result struct {
      Type  string `json:"type"`
      Video *Video `json:"video"`
   }
   if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
      return nil, err
   }
   return result.Video, nil
}
