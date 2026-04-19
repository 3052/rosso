package nbc

import (
   "41.neocities.org/maya"
   "encoding/json"
   "errors"
   "fmt"
   "io"
   "net/url"
)

func (m *Metadata) Stream() (*Stream, error) {
   resp, err := maya.Get(
      &url.URL{
         Scheme: "https",
         Host:   "lemonade.nbc.com",
         Path:   fmt.Sprintf("/v1/vod/%v/%v", m.MpxAccountId, m.MpxGuid),
         RawQuery: url.Values{
            "platform":        {"web"},
            "programmingType": {m.ProgrammingType},
         }.Encode(),
      },
      nil,
   )
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   if resp.StatusCode != 200 {
      return nil, errors.New(resp.Status)
   }
   result := &Stream{}
   err = json.NewDecoder(resp.Body).Decode(result)
   if err != nil {
      return nil, err
   }
   return result, nil
}

func FetchMetadata(name string) (*Metadata, error) {
   body, err := json.Marshal(map[string]any{
      "query": query_page,
      "variables": map[string]string{
         "app":      "nbc",
         "name":     name,
         "platform": "web",
         "type":     "VIDEO",
         "userId":   "",
      },
   })
   if err != nil {
      return nil, err
   }
   resp, err := maya.Post(
      &url.URL{
         Scheme: "https",
         Host:   "friendship.nbc.com",
         Path:   "/v3/graphql",
      },
      map[string]string{"content-type": "application/json"},
      body,
   )
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   if resp.StatusCode != 200 {
      return nil, errors.New(resp.Status)
   }
   var result struct {
      Data struct {
         Page struct {
            Metadata Metadata
         }
      }
      Errors []struct {
         Message string
      }
   }
   err = json.NewDecoder(resp.Body).Decode(&result)
   if err != nil {
      return nil, err
   }
   if len(result.Errors) >= 1 {
      return nil, errors.New(result.Errors[0].Message)
   }
   return &result.Data.Page.Metadata, nil
}

func FetchWidevine(body []byte) ([]byte, error) {
   resp, err := maya.Post(
      &url.URL{
         Scheme:   "https",
         Host:     "drmproxy.digitalsvc.apps.nbcuni.com",
         Path:     "/drm-proxy/license/widevine",
         RawQuery: build_query("widevine"),
      },
      map[string]string{"content-type": "application/octet-stream"},
      body,
   )
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   return io.ReadAll(resp.Body)
}
