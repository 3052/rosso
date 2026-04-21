package tubi

import (
   "encoding/json"
   "net/url"
   "strconv"

   "41.neocities.org/maya"
)

type ContentResponse struct {
   VideoResources []VideoResource `json:"video_resources"`
}

type VideoResource struct {
   Manifest      Manifest      `json:"manifest"`
   LicenseServer LicenseServer `json:"license_server"`
}

type Manifest struct {
   Url      string `json:"url"`
   Duration int    `json:"duration"`
}

type LicenseServer struct {
   Url string `json:"url"`
}

func GetContent(contentId int) (*ContentResponse, error) {
   query := url.Values{}
   query.Set("content_id", strconv.Itoa(contentId))
   query.Set("deviceId", "!")
   query.Add("limit_resolutions[]", "h264_1080p")
   query.Add("limit_resolutions[]", "h265_1080p")
   query.Set("platform", "web")
   query.Add("video_resources[]", "dash")
   query.Add("video_resources[]", "dash_widevine")

   target := &url.URL{
      Scheme:   "https",
      Host:     "uapi.adrise.tv",
      Path:     "/cms/content",
      RawQuery: query.Encode(),
   }

   response, err := maya.Get(target, nil)
   if err != nil {
      return nil, err
   }
   defer response.Body.Close()

   var content ContentResponse
   if err := json.NewDecoder(response.Body).Decode(&content); err != nil {
      return nil, err
   }

   return &content, nil
}
