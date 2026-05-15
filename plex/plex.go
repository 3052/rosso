package plex

import (
   "41.neocities.org/maya"
   "encoding/json"
   "errors"
   "io"
   "net/url"
   "strings"
)

// https://watch.plex.tv/embed/movie/memento-2000
// https://watch.plex.tv/movie/memento-2000
// https://watch.plex.tv/watch/movie/memento-2000
func ParsePath(input *url.URL) string {
   input.Path = strings.TrimPrefix(input.Path, "/embed")
   return strings.TrimPrefix(input.Path, "/watch")
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

type VodMetadata struct {
   Metadata []MetadataItem `json:"Metadata"`
}

type MetadataItem struct {
   Guid  string  `json:"guid"`
   Title string  `json:"title"`
   Media []Media `json:"Media"`
}

type Media struct {
   Id       string    `json:"id"`
   Protocol string    `json:"protocol"`
   Part     []VodPart `json:"Part"`
}

type VodPart struct {
   Id      string `json:"id"`
   Key     string `json:"key"`
   License string `json:"license"`
}

func GetVodMetadata(match *MatchItem, userData *User) (*VodMetadata, error) {
   endpoint := &url.URL{
      Scheme: "https",
      Host:   "vod.provider.plex.tv",
      Path:   match.Key,
   }

   headers := map[string]string{
      "x-plex-token": userData.AuthToken,
   }

   resp, err := maya.Get(endpoint, headers)
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

type MatchContainer struct {
   Metadata []MatchItem `json:"Metadata"`
}

type MatchItem struct {
   Guid      string `json:"guid"`
   Key       string `json:"key"`
   RatingKey string `json:"ratingKey"`
   Title     string `json:"title"`
   Type      string `json:"type"`
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

   resp, err := maya.Get(endpoint, nil)
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

   resp, err := maya.Post(endpoint, nil, body)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()

   return io.ReadAll(resp.Body)
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

   headers := map[string]string{
      "X-Plex-Client-Identifier": "!",
      "X-Plex-Product":           "Plex Mediaverse",
   }

   resp, err := maya.Post(endpoint, headers, nil)
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
