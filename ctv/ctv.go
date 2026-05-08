package ctv

import (
   "41.neocities.org/maya"
   _ "embed"
   "encoding/json"
   "errors"
   "fmt"
   "io"
   "net/url"
   "strings"
)

func (a *AxisContent) Manifest(play *Playback) (*url.URL, error) {
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
      return nil, err
   }
   defer resp.Body.Close()
   data, err := io.ReadAll(resp.Body)
   if err != nil {
      return nil, err
   }
   if resp.StatusCode != 200 {
      var result struct {
         Message string // 2026-05-07
      }
      err = json.Unmarshal(data, &result)
      if err != nil {
         return nil, err
      }
      return nil, errors.New(result.Message)
   }
   return url.Parse(strings.Replace(string(data), "/best/", "/ultimate/", 1))
}

type Playback struct {
   ContentPackages []struct {
      Id int
   }
}

type ResolvedPath struct {
   LastSegment struct {
      Content struct {
         FirstPlayableContent *struct {
            Id string
         }
         Id string
      }
   }
}

func (r *ResolvedPath) get_id() string {
   if fpc := r.LastSegment.Content.FirstPlayableContent; fpc != nil {
      return fpc.Id
   }
   return r.LastSegment.Content.Id
}

//go:embed resolvePath.gql
var query_resolve_path string

//go:embed axisContent.gql
var query_axis_content string

// https://ctv.ca/shows/friends/the-one-with-the-bullies-s2e21
func GetPath(urlData string) (string, error) {
   urlParse, err := url.Parse(urlData)
   if err != nil {
      return "", err
   }
   if urlParse.Scheme == "" {
      return "", errors.New("invalid URL: scheme is missing")
   }
   return urlParse.Path, nil
}

type AxisContent struct {
   AxisId                int
   AxisPlaybackLanguages []struct {
      DestinationCode string
   }
}

func (a *AxisContent) Playback() (*Playback, error) {
   resp, err := maya.Get(
      &url.URL{
         Scheme: "https",
         Host:   "capi.9c9media.com",
         Path: fmt.Sprintf(
            "/destinations/%v/platforms/desktop/contents/%v",
            a.AxisPlaybackLanguages[0].DestinationCode, a.AxisId,
         ),
         RawQuery: "$include=[ContentPackages]",
      },
      nil,
   )
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   result := &Playback{}
   err = json.NewDecoder(resp.Body).Decode(result)
   if err != nil {
      return nil, err
   }
   return result, nil
}

func FetchWidevine(body []byte) ([]byte, error) {
   resp, err := maya.Post(
      &url.URL{Scheme: "https", Host: "license.9c9media.ca", Path: "/widevine"},
      map[string]string{"content-type": "application/x-protobuf"},
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

func (r *ResolvedPath) AxisContent() (*AxisContent, error) {
   body, err := json.Marshal(map[string]any{
      "query": query_axis_content,
      "variables": map[string]string{
         "id": r.get_id(),
      },
   })
   if err != nil {
      return nil, err
   }
   resp, err := maya.Post(
      &url.URL{
         Scheme: "https",
         Host:   "www.ctv.ca",
         Path:   "/space-graphql/apq/graphql",
      },
      // you need this for the first request, then can omit
      map[string]string{"graphql-client-platform": "entpay_web"},
      body,
   )
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   var result struct {
      Data struct {
         AxisContent AxisContent
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
   return &result.Data.AxisContent, nil
}

func Resolve(path string) (*ResolvedPath, error) {
   body, err := json.Marshal(map[string]any{
      "query": query_resolve_path,
      "variables": map[string]string{
         "path": path,
      },
   })
   if err != nil {
      return nil, err
   }
   resp, err := maya.Post(
      &url.URL{
         Scheme: "https",
         Host:   "www.ctv.ca",
         Path:   "/space-graphql/apq/graphql",
      },
      // you need this for the first request, then can omit
      map[string]string{"graphql-client-platform": "entpay_web"},
      body,
   )
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   body, err = io.ReadAll(resp.Body)
   if err != nil {
      return nil, err
   }
   var result struct {
      Data struct {
         ResolvedPath *ResolvedPath
      }
   }
   err = json.Unmarshal(body, &result)
   if err != nil {
      return nil, err
   }
   if result.Data.ResolvedPath == nil {
      return nil, errors.New(string(body))
   }
   return result.Data.ResolvedPath, nil
}
