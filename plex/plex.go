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

// do logs and executes the request.
func do(req *http.Request) (*http.Response, error) {
   log.Println(req.Method, req.URL)
   return http.DefaultClient.Do(req)
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

type MatchItem struct {
   Guid      string `json:"guid"`
   Key       string `json:"key"`
   RatingKey string `json:"ratingKey"`
   Title     string `json:"title"`
   Type      string `json:"type"`
}

type Media struct {
   Id       string    `json:"id"`
   Protocol string    `json:"protocol"`
   Part     []VodPart `json:"Part"`
}

func (*Media) CachePath() string {
   return "rosso/plex/Media"
}

func (m *Media) GetManifest(userData *User) *url.URL {
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

type MetadataItem struct {
   Guid  string  `json:"guid"`
   Title string  `json:"title"`
   Media []Media `json:"Media"`
}

type User struct {
   Id        int    `json:"id"`
   Uuid      string `json:"uuid"`
   AuthToken string `json:"authToken"`
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
   req.Header.Set("X-Plex-Client-Identifier", "!")
   req.Header.Set("X-Plex-Product", "Plex Mediaverse")

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

func (*User) CachePath() string {
   return "rosso/plex/User"
}

type VodMetadata struct {
   Metadata []MetadataItem `json:"Metadata"`
}

func GetVodMetadata(match *MatchItem, userData *User) (*VodMetadata, error) {
   endpoint := &url.URL{
      Scheme: "https",
      Host:   "vod.provider.plex.tv",
      Path:   match.Key,
   }

   req, err := http.NewRequest(http.MethodGet, endpoint.String(), nil)
   if err != nil {
      return nil, err
   }
   req.Header.Set("x-plex-token", userData.AuthToken)

   resp, err := do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()

   var result struct {
      MediaContainer VodMetadata `json:"MediaContainer"`
   }
   if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
      return nil, err
   }

   return &result.MediaContainer, nil
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

type VodPart struct {
   Id      string `json:"id"`
   Key     string `json:"key"`
   License string `json:"license"`
}
