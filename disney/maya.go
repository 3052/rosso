package disney

import (
   "41.neocities.org/maya"
   "bytes"
   "encoding/json"
   "errors"
   "io"
   "net/http"
   "net/url"
)

// expires: 4 hours
// request: Account
func (t *Token) Refresh() error {
   if err := t.assert("Account"); err != nil {
      return err
   }
   body, err := json.Marshal(map[string]any{
      "query": mutation_refresh_token,
      "variables": map[string]any{
         "input": map[string]string{
            "refreshToken": t.RefreshToken,
         },
      },
   })
   if err != nil {
      return err
   }
   req, err := http.NewRequest(
      "POST", "https://disney.api.edge.bamgrid.com/graph/v1/device/graphql",
      bytes.NewReader(body),
   )
   if err != nil {
      return err
   }
   req.Header.Set("authorization", "Bearer "+client_api_key)
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return err
   }
   defer resp.Body.Close()
   var result struct {
      Extensions struct {
         Sdk struct {
            Token Token
         }
      }
   }
   err = json.NewDecoder(resp.Body).Decode(&result)
   if err != nil {
      return err
   }
   *t = result.Extensions.Sdk.Token
   return nil
}

// L3 max: 720p
// request: Account
func (t *Token) FetchWidevine(body []byte) ([]byte, error) {
   if err := t.assert("Account"); err != nil {
      return nil, err
   }
   resp, err := maya.Post(
      &url.URL{
         Scheme: "https",
         Host:   "disney.playback.edge.bamgrid.com",
         Path:   "/widevine/v1/obtain-license",
      },
      map[string]string{"authorization": t.AccessToken},
      body,
   )
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   return io.ReadAll(resp.Body)
}

// SL2000 max: 720p
// SL3000 max: 2160p
// request: Account
func (t *Token) FetchPlayReady(body []byte) ([]byte, error) {
   if err := t.assert("Account"); err != nil {
      return nil, err
   }
   resp, err := maya.Post(
      &url.URL{
         Scheme: "https",
         Host:   "disney.playback.edge.bamgrid.com",
         Path:   "/playready/v1/obtain-license.asmx",
      },
      map[string]string{"authorization": t.AccessToken},
      body,
   )
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   if resp.StatusCode != 200 {
      return nil, errors.New(resp.Status)
   }
   return io.ReadAll(resp.Body)
}

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
   resp, err := maya.Post(
      &url.URL{
         Scheme: "https",
         Host:   "disney.playback.edge.bamgrid.com",
         // /v7/playback/ctr-high
         // /v7/playback/tv-drm-ctr-h265-atmos
         Path: "/v7/playback/ctr-regular",
      },
      map[string]string{
         "authorization":           "Bearer " + t.AccessToken,
         "content-type":            "application/json",
         "x-application-version":   "",
         "x-bamsdk-client-id":      "",
         "x-bamsdk-platform":       "",
         "x-bamsdk-version":        "",
         "x-dss-feature-filtering": "true",
      },
      body,
   )
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
