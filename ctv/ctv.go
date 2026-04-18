package ctv

import (
   _ "embed"
   "errors"
   "net/url"
   "strings"
)

func GetManifest(manifest string) (*url.URL, error) {
   return url.Parse(strings.Replace(manifest, "/best/", "/ultimate/", 1))
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
