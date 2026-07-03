package ctv

import (
   "bytes"
   _ "embed"
   "encoding/json"
   "errors"
   "fmt"
   "io"
   "log"
   "net/http"
   "net/url"
   "strings"
)

//go:embed axisContent.gql
var query_axis_content string

//go:embed resolvePath.gql
var query_resolve_path string

func FetchWidevine(body []byte) ([]byte, error) {
   req, err := http.NewRequest("POST", "https://license.9c9media.ca/widevine", bytes.NewReader(body))
   if err != nil {
      return nil, err
   }
   req.Header.Set("Content-Type", "application/x-protobuf")

   resp, err := doRequest(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()

   if resp.StatusCode != http.StatusOK {
      return nil, errors.New(resp.Status)
   }
   return io.ReadAll(resp.Body)
}

// https://ctv.ca/shows/friends/the-one-with-the-bullies-s2e21
func GetPath(urlData string) (string, error) {
   parse, err := url.Parse(urlData)
   if err != nil {
      return "", err
   }
   if parse.Scheme == "" {
      return "", errors.New("invalid URL: scheme is missing")
   }
   return parse.Path, nil
}

// Helper function to perform the request and log Method + URL
func doRequest(req *http.Request) (*http.Response, error) {
   log.Println(req.Method, req.URL)
   return http.DefaultClient.Do(req)
}

type AxisContent struct {
   AxisId                int
   AxisPlaybackLanguages []struct {
      DestinationCode string
   }
}

func (a *AxisContent) Manifest(play *Playback) (*url.URL, error) {
   endpoint := fmt.Sprintf(
      "https://capi.9c9media.com/destinations/%v/platforms/desktop/playback/contents/%v/contentPackages/%v/manifest.mpd?action=reference",
      a.AxisPlaybackLanguages[0].DestinationCode, a.AxisId, play.ContentPackages[0].Id,
   )

   req, err := http.NewRequest("GET", endpoint, nil)
   if err != nil {
      return nil, err
   }

   resp, err := doRequest(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()

   data, err := io.ReadAll(resp.Body)
   if err != nil {
      return nil, err
   }

   if resp.StatusCode != http.StatusOK {
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

func (a *AxisContent) Playback() (*Playback, error) {
   endpoint := fmt.Sprintf(
      "https://capi.9c9media.com/destinations/%v/platforms/desktop/contents/%v?$include=[ContentPackages]",
      a.AxisPlaybackLanguages[0].DestinationCode, a.AxisId,
   )

   req, err := http.NewRequest("GET", endpoint, nil)
   if err != nil {
      return nil, err
   }

   resp, err := doRequest(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()

   if resp.StatusCode != http.StatusOK {
      return nil, errors.New(resp.Status)
   }

   result := &Playback{}
   err = json.NewDecoder(resp.Body).Decode(result)
   if err != nil {
      return nil, err
   }
   return result, nil
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

   req, err := http.NewRequest("POST", "https://www.ctv.ca/space-graphql/apq/graphql", bytes.NewReader(body))
   if err != nil {
      return nil, err
   }
   // you need this for the first request, then can omit
   req.Header.Set("graphql-client-platform", "entpay_web")
   req.Header.Set("Content-Type", "application/json")

   resp, err := doRequest(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()

   if resp.StatusCode != http.StatusOK {
      return nil, errors.New(resp.Status)
   }

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

   req, err := http.NewRequest("POST", "https://www.ctv.ca/space-graphql/apq/graphql", bytes.NewReader(body))
   if err != nil {
      return nil, err
   }
   // you need this for the first request, then can omit
   req.Header.Set("graphql-client-platform", "entpay_web")
   req.Header.Set("Content-Type", "application/json")

   resp, err := doRequest(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()

   if resp.StatusCode != http.StatusOK {
      return nil, errors.New(resp.Status)
   }

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

func (r *ResolvedPath) get_id() string {
   if fpc := r.LastSegment.Content.FirstPlayableContent; fpc != nil {
      return fpc.Id
   }
   return r.LastSegment.Content.Id
}
