package tubi

import (
   "bytes"
   "encoding/json"
   "errors"
   "io"
   "net/http"
   "net/url"
   "strconv"
)

func FetchContent(id int) (*Content, error) {
   req := http.Request{
      URL: &url.URL{
         Scheme: "https",
         Host:   "uapi.adrise.tv",
         Path:   "/cms/content",
         RawQuery: url.Values{
            "content_id":          {strconv.Itoa(id)},
            "deviceId":            {"!"},
            "limit_resolutions[]": {
               "h264_1080p",
               "h265_1080p",
            },
            "platform":            {"web"},
            "video_resources[]": {
               "dash",
               "dash_widevine",
            },
         }.Encode(),
      },
      Header: http.Header{},
   }
   resp, err := http.DefaultClient.Do(&req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   var result Content
   err = json.NewDecoder(resp.Body).Decode(&result)
   if err != nil {
      return nil, err
   }
   if len(result.VideoResources) == 0 {
      return nil, errors.New("no video resources found")
   }
   return &result, nil
}

type Content struct {
   Children     []*Content
   DetailedType string `json:"detailed_type"`
   Id           int    `json:",string"`
   SeriesId     int    `json:"series_id,string"`
   // these should already be in reverse order by resolution
   VideoResources []VideoResource `json:"video_resources"`
}

func (v *VideoResource) Widevine(data []byte) ([]byte, error) {
   resp, err := http.Post(
      v.LicenseServer.Url, "application/x-protobuf", bytes.NewReader(data),
   )
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   return io.ReadAll(resp.Body)
}

func (v *VideoResource) Dash() (*Dash, error) {
   resp, err := http.Get(v.Manifest.Url)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   body, err := io.ReadAll(resp.Body)
   if err != nil {
      return nil, err
   }
   return &Dash{Body: body, Url: resp.Request.URL}, nil
}

type VideoResource struct {
   LicenseServer *struct {
      Url string
   } `json:"license_server"`
   Manifest struct {
      Url string // MPD
   }
   Type string
}

type Dash struct {
   Body []byte
   Url  *url.URL
}
