package kanopy

import (
   "encoding/json"
   "net/url"

   "41.neocities.org/maya"
)

type Video struct {
   VideoId int    `json:"videoId"`
   Title   string `json:"title"`
}

type VideoResponse struct {
   Type  string `json:"type"`
   Video Video  `json:"video"`
   Alias string `json:"alias"`
}

func GetVideo(loginData *Login, alias string) (*VideoResponse, error) {
   endpoint := &url.URL{
      Scheme: "https",
      Host:   "www.kanopy.com",
      Path:   "/kapi/videos/alias/" + alias,
   }

   headers := map[string]string{
      "authorization": "Bearer " + loginData.Jwt,
      "x-version":     "!/!/!/!",
   }

   resp, err := maya.Get(endpoint, headers)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()

   var video VideoResponse
   if err := json.NewDecoder(resp.Body).Decode(&video); err != nil {
      return nil, err
   }

   return &video, nil
}
