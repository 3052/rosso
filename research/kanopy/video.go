package kanopy

import (
   "41.neocities.org/maya"
   "encoding/json"
   "errors"
   "net/url"
   "path"
   "strconv"
   "strings"
)

// Supports URLs such as:
// - https://kanopy.com/video/6440418
// - https://kanopy.com/video/genius-party
// - https://kanopy.com/en/video/genius-party
// - https://kanopy.com/en/product/genius-party
func ParseVideo(urlData string) (*Video, error) {
   url_parse, err := url.Parse(urlData)
   if err != nil {
      return nil, err
   }
   if !strings.Contains(url_parse.Host, "kanopy.com") {
      return nil, errors.New("invalid domain")
   }
   // Get the directory of the path (removes the final identifier).
   // e.g., "/en/product/genius-party" -> "/en/product"
   dir := path.Dir(url_parse.Path)
   // Check if the directory ends with "/video" OR "/product".
   // This supports:
   // - /video/{id}
   // - /en/video/{id}
   // - /en/product/{id}
   if !strings.HasSuffix(dir, "/video") && !strings.HasSuffix(dir, "/product") {
      return nil, errors.New("invalid path structure")
   }
   v := &Video{}
   identifier := path.Base(url_parse.Path)
   numeric_id, err := strconv.Atoi(identifier)
   if err != nil {
      v.Alias = identifier
   } else {
      v.VideoId = numeric_id
   }
   return v, nil
}

type Video struct {
   VideoId int    `json:"videoId"`
   Title   string `json:"title"`
   Alias   string `json:"alias"`
}

type VideoResponse struct {
   Type  string `json:"type"`
   Video *Video `json:"video"`
}

func GetVideo(alias, jwt string) (*VideoResponse, error) {
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

   var videoResp VideoResponse
   if err := json.NewDecoder(resp.Body).Decode(&videoResp); err != nil {
      return nil, err
   }

   return &videoResp, nil
}
