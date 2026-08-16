package plex

import (
   "bytes"
   "encoding/json"
   "errors"
   "io"
   "log"
   "net/http"
   "net/url"
   "strings"
)

func AcquireWidevineLicense(mediaData *Media, userData *User, body []byte) ([]byte, error) {
   log.Print("AcquireWidevineLicense")
   if len(mediaData.Part) == 0 {
      return nil, errors.New("no media parts found")
   }
   if mediaData.Part[0].License == "" {
      return nil, errors.New("no license path found")
   }

   endpoint := &url.URL{
      Scheme: "https",
      Host:   "vod.provider.plex.tv",
      Path:   mediaData.Part[0].License,
   }

   query := url.Values{}
   query.Set("x-plex-drm", "widevine")
   query.Set("x-plex-token", userData.AuthToken)
   endpoint.RawQuery = query.Encode()

   // Keep the body nil when empty so no Content-Length/Content-Type is sent.
   var bodyReader io.Reader
   if body != nil {
      bodyReader = bytes.NewReader(body)
   }

   req, err := http.NewRequest(http.MethodPost, endpoint.String(), bodyReader)
   if err != nil {
      return nil, err
   }

   resp, err := do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()

   return io.ReadAll(resp.Body)
}

// https://watch.plex.tv/embed/movie/memento-2000
// https://watch.plex.tv/movie/memento-2000
// https://watch.plex.tv/watch/movie/memento-2000
func ParsePath(rawUrl string) (string, error) {
   // Find the starting position of the "/movie/" marker.
   startIndex := strings.Index(rawUrl, "/movie/")
   if startIndex == -1 {
      return "", errors.New("no /movie/ segment found in URL")
   }
   // The slug must not be empty. Check if the string ends right after "/movie/".
   if len(rawUrl) == startIndex+len("/movie/") {
      return "", errors.New("movie slug is empty")
   }
   // Return the slice from the start of the marker to the end of the string.
   return rawUrl[startIndex:], nil
}

func CreateUser() (*User, error) {
   endpoint := &url.URL{
      Scheme: "https",
      Host:   "plex.tv",
      Path:   "/api/v2/users/anonymous",
   }

   req, err := http.NewRequest(http.MethodPost, endpoint.String(), nil)
   if err != nil {
      return nil, err
   }
   req.Header.Set("accept", "application/json")
   req.Header.Set("x-plex-client-identifier", "x")
   req.Header.Set("x-plex-product", "Plex Mediaverse")

   resp, err := do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()

   var result User
   if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
      return nil, err
   }
   return &result, nil
}

type MatchContainer struct {
   Metadata []MatchItem `json:"Metadata"`
}

func GetMetadataMatches(urlPath string, userData *User) (*MatchContainer, error) {
   endpoint := &url.URL{
      Scheme: "https",
      Host:   "discover.provider.plex.tv",
      Path:   "/library/metadata/matches",
   }

   query := url.Values{}
   query.Set("url", urlPath)
   query.Set("x-plex-token", userData.AuthToken)
   endpoint.RawQuery = query.Encode()

   req, err := http.NewRequest(http.MethodGet, endpoint.String(), nil)
   if err != nil {
      return nil, err
   }
   req.Header.Set("accept", "application/json")

   resp, err := do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()

   var result struct {
      MediaContainer MatchContainer `json:"MediaContainer"`
   }
   if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
      return nil, err
   }

   return &result.MediaContainer, nil
}

func (*Media) CachePath() string {
   return "rosso/plex/Media"
}

func (m *Media) GetManifest(userData *User) *url.URL {
   log.Print("GetManifest")
   endpoint := &url.URL{
      Scheme: "https",
      Host:   "vod.provider.plex.tv",
      Path:   m.Part[0].Key,
   }
   query := url.Values{}
   query.Set("x-plex-token", userData.AuthToken)
   endpoint.RawQuery = query.Encode()
   return endpoint
}

func (*User) CachePath() string {
   return "rosso/plex/User"
}

func (vod *VodMetadata) GetDash() (*Media, error) {
   for _, item := range vod.Metadata {
      for _, media_data := range item.Media {
         if media_data.Protocol == "dash" {
            return &media_data, nil
         }
      }
   }
   return nil, errors.New("dash media not found")
}
