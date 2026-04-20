// video.go
package kanopy

import (
   "encoding/json"
   "net/url"

   "41.neocities.org/maya"
)

type VideoResponse struct {
   Type  string `json:"type"`
   Video struct {
      VideoId int    `json:"videoId"`
      Title   string `json:"title"`
   } `json:"video"`
}

func GetVideo(Authorization string, Alias string) (*VideoResponse, error) {
   targetUrl, parseError := url.Parse("https://www.kanopy.com/kapi/videos/alias/" + Alias)
   if parseError != nil {
      return nil, parseError
   }

   headers := map[string]string{
      "authorization": "Bearer " + Authorization,
      "x-version":     "!/!/!/!",
      "user-agent":    "Go-http-client/2.0",
   }

   resp, requestError := maya.Get(targetUrl, headers)
   if requestError != nil {
      return nil, requestError
   }
   defer resp.Body.Close()

   var videoResp VideoResponse
   decodeError := json.NewDecoder(resp.Body).Decode(&videoResp)
   if decodeError != nil {
      return nil, decodeError
   }
   return &videoResp, nil
}
