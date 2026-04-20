package kanopy

import (
   "encoding/json"
   "io"
   "net/url"

   "41.neocities.org/maya"
)

type VideoDetails struct {
   Type  string `json:"type"`
   Video Video  `json:"video"`
}

type Video struct {
   VideoID            int      `json:"videoId"`
   Title              string   `json:"title"`
   DescriptionHTML    string   `json:"descriptionHtml"`
   Images             Images   `json:"images"`
   HasBurntInCaptions bool     `json:"hasBurntInCaptions"`
   HasCaptions        bool     `json:"hasCaptions"`
   CaptionLanguages   []string `json:"captionLanguages"`
   ProductionYear     int      `json:"productionYear"`
   IsKids             bool     `json:"isKids"`
   DurationSeconds    int      `json:"durationSeconds"`
   IsFree             bool     `json:"isFree"`
   IsRequestable      bool     `json:"isRequestable"`
   Alias              string   `json:"alias"`
   FeedID             int      `json:"feedId"`
}

type Images struct {
   Landscapes ImageSet `json:"landscapes"`
   Posters    ImageSet `json:"posters"`
}

type ImageSet struct {
   Small  string `json:"small"`
   Medium string `json:"medium"`
   Large  string `json:"large"`
}

func GetVideoDetails(alias string, jwt string) (*VideoDetails, error) {
   target := &url.URL{
      Scheme: "https",
      Host:   "www.kanopy.com",
      Path:   "/kapi/videos/alias/" + alias,
   }

   headers := map[string]string{
      "x-version":     "!/!/!/!",
      "authorization": "Bearer " + jwt,
      "user-agent":    "Go-http-client/2.0",
   }

   resp, err := maya.Get(target, headers)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()

   respBytes, err := io.ReadAll(resp.Body)
   if err != nil {
      return nil, err
   }

   var videoDetails VideoDetails
   if err := json.Unmarshal(respBytes, &videoDetails); err != nil {
      return nil, err
   }

   return &videoDetails, nil
}
