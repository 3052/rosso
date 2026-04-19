package disney

import (
   "41.neocities.org/maya"
   "bytes"
   "encoding/json"
   "net/http"
   "net/url"
)

// request: Account
func (t *Token) FetchStream(mediaId string) (*Stream, error) {
   if err := t.assert("Account"); err != nil {
      return nil, err
   }
   playback_id, err := json.Marshal(map[string]string{
      "mediaId": mediaId,
   })
   if err != nil {
      return nil, err
   }
   body, err := json.Marshal(map[string]any{
      "playback": map[string]any{
         "attributes": map[string]any{
            "assetInsertionStrategy": "SGAI",
            "codecs": map[string]any{
               "supportsMultiCodecMaster": true, // 4K
               "video": []string{
                  "h.264",
                  "h.265",
               },
            },
            "videoRanges": []string{"HDR10"},
         },
      },
      "playbackId": playback_id,
   })
   if err != nil {
      return nil, err
   }
   // /v7/playback/ctr-high
   // /v7/playback/tv-drm-ctr-h265-atmos
   req, err := http.NewRequest(
      "POST", "https://disney.playback.edge.bamgrid.com/v7/playback/ctr-regular",
      bytes.NewReader(body),
   )
   if err != nil {
      return nil, err
   }
   req.Header.Set("authorization", "Bearer "+t.AccessToken)
   req.Header.Set("content-type", "application/json")
   req.Header.Set("x-application-version", "")
   req.Header.Set("x-bamsdk-client-id", "")
   req.Header.Set("x-bamsdk-platform", "")
   req.Header.Set("x-bamsdk-version", "")
   req.Header.Set("x-dss-feature-filtering", "true")
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   var result struct {
      Errors []Error
      Stream Stream
   }
   err = json.NewDecoder(resp.Body).Decode(&result)
   if err != nil {
      return nil, err
   }
   if len(result.Errors) >= 1 {
      return nil, &result.Errors[0]
   }
   return &result.Stream, nil
}

// request: Account
func (t *Token) FetchPage(entity string) (*Page, error) {
   if err := t.assert("Account"); err != nil {
      return nil, err
   }
   resp, err := maya.Get(
      &url.URL{
         Scheme:   "https",
         Host:     "disney.api.edge.bamgrid.com",
         Path:     "/explore/v1.12/page/entity-" + entity,
         RawQuery: "limit=0",
      },
      map[string]string{"authorization": "Bearer " + t.AccessToken},
   )
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   var result struct {
      Data struct {
         Errors []Error // 2026-04-11
         Page   Page
      }
   }
   err = json.NewDecoder(resp.Body).Decode(&result)
   if err != nil {
      return nil, err
   }
   if len(result.Data.Errors) >= 1 {
      return nil, &result.Data.Errors[0]
   }
   return &result.Data.Page, nil
}
