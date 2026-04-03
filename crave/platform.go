package crave

import (
   "bytes"
   "encoding/json"
   "errors"
   "fmt"
   "io"
   "net/http"
   "net/url"
)

// L3 max 720p
func (c *ContentPackage) FetchWidevine(contentId int, accessToken string, payload []byte) ([]byte, error) {
   data, err := marshal(map[string]any{
      "payload": payload,
      "playbackContext": map[string]any{
         "contentId":        contentId,
         "contentpackageId": c.Id, // lower-case 'p' as per their API
         "platformId":       1,    // Hardcoded to 1 for Web
         "destinationId":    c.DestinationId,
         "jwt":              accessToken,
      },
   })
   if err != nil {
      return nil, err
   }
   req, err := http.NewRequest(
      "POST", "https://license.9c9media.com/widevine", bytes.NewBuffer(data),
   )
   if err != nil {
      return nil, err
   }
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   data, err = io.ReadAll(resp.Body)
   if err != nil {
      return nil, err
   }
   if resp.StatusCode != http.StatusOK {
      var result struct {
         Message string
      }
      err = json.Unmarshal(data, &result)
      if err != nil {
         return nil, err
      }
      return nil, errors.New(result.Message)
   }
   return data, nil
}

func (c *ContentPackage) FetchPlayReady(contentId int, accessToken string, payload []byte) ([]byte, error) {
   data, err := marshal(map[string]any{
      "payload": payload,
      "playbackContext": map[string]any{
         "contentId":        contentId,
         "contentpackageId": c.Id, // lower-case 'p' as per their API
         "platformId":       1,    // Hardcoded to 1 for Web
         "destinationId":    c.DestinationId,
         "jwt":              accessToken,
      },
   })
   if err != nil {
      return nil, err
   }
   req, err := http.NewRequest(
      "POST", "https://license.9c9media.com/playready", bytes.NewBuffer(data),
   )
   if err != nil {
      return nil, err
   }
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   data, err = io.ReadAll(resp.Body)
   if err != nil {
      return nil, err
   }
   if resp.StatusCode != http.StatusOK {
      var result struct {
         Message string
      }
      err = json.Unmarshal(data, &result)
      if err != nil {
         return nil, err
      }
      return nil, errors.New(result.Message)
   }
   return data, nil
}

func (c *ContentPackage) FetchManifest(contentId int, accessToken string) (*Manifest, error) {
   req := http.Request{
      URL: &url.URL{
         Scheme: "https",
         Host:   "stream.video.9c9media.com",
         Path: fmt.Sprintf(
            "/meta/content/%v/contentpackage/%v/destination/%v/platform/1",
            contentId, c.Id, c.DestinationId,
         ),
         RawQuery: url.Values{
            "filter": {"ff"}, // 1080p
            "format": {"mpd"},
            "hd":     {"true"}, // 1080p
            "mcv":    {"true"}, // H.264 + HEVC
            "uhd":    {"true"}, // HEVC
         }.Encode(),
      },
      Header: http.Header{},
   }
   req.Header.Set("authorization", "Bearer "+accessToken)
   resp, err := http.DefaultClient.Do(&req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   var result Manifest
   err = json.NewDecoder(resp.Body).Decode(&result)
   if err != nil {
      return nil, err
   }
   if result.Message != "" {
      return nil, errors.New(result.Message)
   }
   return &result, nil
}
