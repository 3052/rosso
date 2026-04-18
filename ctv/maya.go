package ctv

import (
   "41.neocities.org/maya"
   "bytes"
   "encoding/json"
   "errors"
   "fmt"
   "io"
   "net/http"
   "net/url"
   "strings"
)

func Resolve(path string) (*ResolvedPath, error) {
   data, err := json.Marshal(map[string]any{
      "query": query_resolve_path,
      "variables": map[string]string{
         "path": path,
      },
   })
   if err != nil {
      return nil, err
   }
   req, err := http.NewRequest(
      "POST", "https://www.ctv.ca/space-graphql/apq/graphql",
      bytes.NewReader(data),
   )
   if err != nil {
      return nil, err
   }
   // you need this for the first request, then can omit
   req.Header.Set("graphql-client-platform", "entpay_web")
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   data, err = io.ReadAll(resp.Body)
   if err != nil {
      return nil, err
   }
   var result struct {
      Data struct {
         ResolvedPath *ResolvedPath
      }
   }
   err = json.Unmarshal(data, &result)
   if err != nil {
      return nil, err
   }
   if result.Data.ResolvedPath == nil {
      return nil, errors.New(string(data))
   }
   return result.Data.ResolvedPath, nil
}

func (a *AxisContent) Manifest(play *Playback) (string, error) {
   resp, err := maya.Get(
      &url.URL{
         Scheme: "https",
         Host:   "capi.9c9media.com",
         Path: fmt.Sprint(
            "/destinations/", a.AxisPlaybackLanguages[0].DestinationCode,
            "/platforms/desktop/playback/contents/", a.AxisId,
            "/contentPackages/", play.ContentPackages[0].Id,
            "/manifest.mpd",
         ),
         RawQuery: "action=reference",
      },
      nil,
   )
   if err != nil {
      return "", err
   }
   defer resp.Body.Close()
   if resp.StatusCode != 200 {
      return "", errors.New(resp.Status)
   }
   var data strings.Builder
   _, err = io.Copy(&data, resp.Body)
   if err != nil {
      return "", err
   }
   return data.String(), nil
}
