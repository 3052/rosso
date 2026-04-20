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
   VideoID     int    `json:"videoId"`
   Title       string `json:"title"`
   Description string `json:"descriptionHtml"`
}

type VideoResponse struct {
   Type  string `json:"type"`
   Video Video  `json:"video"`
}

func GetVideo(alias string, authorization string) (*VideoResponse, error) {
   targetURL, err := url.Parse(fmt.Sprintf("https://www.kanopy.com/kapi/videos/alias/%s", alias))
   if err != nil {
      return nil, err
   }

   headers := map[string]string{
      "x-version":     "!/!/!/!",
      "authorization": authorization,
      "user-agent":    "Go-http-client/2.0",
   }

   resp, err := maya.Get(targetURL, headers)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()

   if resp.StatusCode != 200 {
      return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
   }

   respBody, err := io.ReadAll(resp.Body)
   if err != nil {
      return nil, err
   }

   var videoResp VideoResponse
   if err := json.Unmarshal(respBody, &videoResp); err != nil {
      return nil, err
   }

   return &videoResp, nil
}
