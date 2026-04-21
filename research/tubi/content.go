// File: content.go
package tubi

import (
   "encoding/json"
   "net/url"

   "41.neocities.org/maya"
)

type ContentResponse struct {
   VideoResources []VideoResource `json:"video_resources"`
}

type VideoResource struct {
   Type          string        `json:"type"`
   Codec         string        `json:"codec"`
   Resolution    string        `json:"resolution"`
   Manifest      Manifest      `json:"manifest"`
   LicenseServer LicenseServer `json:"license_server"`
}

type Manifest struct {
   Url string `json:"url"`
}

type LicenseServer struct {
   Url string `json:"url"`
}

func GetContent(contentId string, deviceId string) (*ContentResponse, error) {
   targetUrl := &url.URL{
      Scheme: "https",
      Host:   "uapi.adrise.tv",
      Path:   "/cms/content",
   }

   q := targetUrl.Query()
   q.Set("content_id", contentId)
   q.Set("deviceId", deviceId)
   q.Add("limit_resolutions[]", "h264_1080p")
   q.Add("limit_resolutions[]", "h265_1080p")
   q.Set("platform", "web")
   q.Add("video_resources[]", "dash")
   q.Add("video_resources[]", "dash_widevine")
   targetUrl.RawQuery = q.Encode()

   resp, err := maya.Get(targetUrl, nil)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()

   var response ContentResponse
   if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
      return nil, err
   }

   return &response, nil
}
