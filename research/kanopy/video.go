package kanopy

import (
   "encoding/json"
   "io"
   "net/url"

   "41.neocities.org/maya"
)

type Video struct {
   VideoId         int    `json:"videoId"`
   Title           string `json:"title"`
   DescriptionHtml string `json:"descriptionHtml"`
   DurationSeconds int    `json:"durationSeconds"`
}

type VideoResponse struct {
   Type  string `json:"type"`
   Video Video  `json:"video"`
}

func GetVideo(login *LoginResponse, alias string) (*VideoResponse, error) {
   endpoint := &url.URL{
      Scheme: "https",
      Host:   "www.kanopy.com",
      Path:   "/kapi/videos/alias/" + alias,
   }

   headers := map[string]string{
      "authorization": "Bearer " + login.Jwt,
   }

   resp, err := maya.Get(endpoint, headers)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()

   respBody, err := io.ReadAll(resp.Body)
   if err != nil {
      return nil, err
   }

   var video VideoResponse
   if err := json.Unmarshal(respBody, &video); err != nil {
      return nil, err
   }

   return &video, nil
}
