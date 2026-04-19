package tubi

import (
   "41.neocities.org/maya"
   "encoding/json"
   "fmt"
   "io"
   "net/url"
   "strconv"
)

func GetContent(contentId int) (*Content, error) {
   query := make(url.Values)
   query.Set("content_id", strconv.Itoa(contentId))
   query.Set("deviceId", "!")
   query.Add("limit_resolutions[]", "h264_1080p")
   query.Add("limit_resolutions[]", "h265_1080p")
   query.Set("platform", "web")
   query.Add("video_resources[]", "dash")
   query.Add("video_resources[]", "dash_widevine")
   resp, err := maya.Get(&url.URL{
      Scheme:   "https",
      Host:     "uapi.adrise.tv",
      Path:     "/cms/content",
      RawQuery: query.Encode(),
   }, nil)
   if err != nil {
      return nil, fmt.Errorf("request failed: %w", err)
   }
   defer resp.Body.Close()

   if resp.StatusCode != 200 {
      return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
   }

   var content Content
   if err := json.NewDecoder(resp.Body).Decode(&content); err != nil {
      return nil, fmt.Errorf("failed to decode JSON response: %w", err)
   }

   return &content, nil
}

func (s *LicenseServer) PostLicense(body []byte) ([]byte, error) {
   target, err := url.Parse(s.Url)
   if err != nil {
      return nil, err
   }
   resp, err := maya.Post(
      target, map[string]string{"content-type": "application/x-protobuf"}, body,
   )
   if err != nil {
      return nil, fmt.Errorf("request failed: %w", err)
   }
   defer resp.Body.Close()

   if resp.StatusCode != 200 {
      return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
   }

   body, err = io.ReadAll(resp.Body)
   if err != nil {
      return nil, fmt.Errorf("failed to read response body: %w", err)
   }

   return body, nil
}
