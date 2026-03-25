package ctv

import (
   "bytes"
   _ "embed"
   "encoding/json"
   "errors"
   "fmt"
   "io"
   "net/http"
   "net/url"
   "strings"
)

func FetchDash(manifest string) (*Dash, error) {
   resp, err := http.Get(strings.Replace(manifest, "/best/", "/ultimate/", 1))
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

func (r *ResolvedPath) get_id() string {
   if fpc := r.LastSegment.Content.FirstPlayableContent; fpc != nil {
      return fpc.Id
   }
   return r.LastSegment.Content.Id
}

func (r *ResolvedPath) AxisContent() (*AxisContent, error) {
   data, err := json.Marshal(map[string]any{
      "query": query_axis_content,
      "variables": map[string]string{
         "id": r.get_id(),
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

func Widevine(data []byte) ([]byte, error) {
   resp, err := http.Post(
      "https://license.9c9media.ca/widevine", "application/x-protobuf",
      bytes.NewReader(data),
   )
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   if resp.StatusCode != http.StatusOK {
      return nil, errors.New(resp.Status)
   }
   return io.ReadAll(resp.Body)
}

func (a *AxisContent) Playback() (*Playback, error) {
   var req http.Request
   req.Header = http.Header{}
   req.URL = &url.URL{
      Scheme: "https",
      Host:   "capi.9c9media.com",
      Path: fmt.Sprintf(
         "/destinations/%v/platforms/desktop/contents/%v",
         a.AxisPlaybackLanguages[0].DestinationCode, a.AxisId,
      ),
      RawQuery: "$include=[ContentPackages]",
   }
   resp, err := http.DefaultClient.Do(&req)
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

type Dash struct {
   Body []byte
   Url  *url.URL
}

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

type Playback struct {
   ContentPackages []struct {
      Id int
   }
}

type AxisContent struct {
   AxisId                int
   AxisPlaybackLanguages []struct {
      DestinationCode string
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

///

func (a *AxisContent) Manifest(play *Playback) (string, error) {
   var req http.Request
   req.URL = &url.URL{
      Scheme: "https",
      Host:   "capi.9c9media.com",
      Path: fmt.Sprint(
         "/destinations/", a.AxisPlaybackLanguages[0].DestinationCode,
         "/platforms/desktop/playback/contents/", a.AxisId,
         "/contentPackages/", play.ContentPackages[0].Id,
         "/manifest.mpd",
      ),
      RawQuery: "action=reference",
   }
   req.Header = http.Header{}
   resp, err := http.DefaultClient.Do(&req)
   if err != nil {
      return "", err
   }
   defer resp.Body.Close()
   if resp.StatusCode != http.StatusOK {
      return "", errors.New(resp.Status)
   }
   var data strings.Builder
   _, err = io.Copy(&data, resp.Body)
   if err != nil {
      return "", err
   }
   return data.String(), nil
}
