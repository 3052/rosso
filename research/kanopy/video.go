// video.go
package kanopy

import (
   "encoding/json"
   "fmt"
   "io"
   "net/url"

   "41.neocities.org/maya"
)

type Video struct {
   VideoId int `json:"videoId"`
}

type GetVideoResponse struct {
   Type  string `json:"type"`
   Video Video  `json:"video"`
}

func GetVideo(alias string, authorization string) (*GetVideoResponse, error) {
   targetUrl, err := url.Parse(fmt.Sprintf("https://www.kanopy.com/kapi/videos/alias/%s", url.PathEscape(alias)))
   if err != nil {
      return nil, err
   }

   headers := map[string]string{
      "authorization": authorization,
   }

   resp, err := maya.Get(targetUrl, headers)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()

   respBody, err := io.ReadAll(resp.Body)
   if err != nil {
      return nil, err
   }

   var videoResp GetVideoResponse
   if err := json.Unmarshal(respBody, &videoResp); err != nil {
      return nil, err
   }

   return &videoResp, nil
}
