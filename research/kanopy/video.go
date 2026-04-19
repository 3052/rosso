package kanopy

import (
   "encoding/json"
   "fmt"
   "io"
   "net/url"

   "41.neocities.org/maya"
)

type VideoData struct {
   VideoId int    `json:"videoId"`
   Title   string `json:"title"`
}

type VideoResponse struct {
   Type  string    `json:"type"`
   Video VideoData `json:"video"`
}

func GetVideo(alias string, jwt string) (*VideoResponse, error) {
   targetUrl, err := url.Parse(fmt.Sprintf("https://www.kanopy.com/kapi/videos/alias/%s", alias))
   if err != nil {
      return nil, err
   }

   requestHeaders := map[string]string{
      "authorization": "Bearer " + jwt,
   }

   response, err := maya.Get(targetUrl, requestHeaders)
   if err != nil {
      return nil, err
   }
   defer response.Body.Close()

   responseBytes, err := io.ReadAll(response.Body)
   if err != nil {
      return nil, err
   }

   var videoResponse VideoResponse
   err = json.Unmarshal(responseBytes, &videoResponse)
   if err != nil {
      return nil, err
   }

   return &videoResponse, nil
}
