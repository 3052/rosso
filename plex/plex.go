package plex

import (
   "41.neocities.org/maya"
   "encoding/json"
   "errors"
   "io"
   "net/url"
   "strings"
)

// /movie/memento-2000
func FetchMatch(token, path string) (*Metadata, error) {
   resp, err := maya.Get(
      &url.URL{
         Scheme: "https",
         Host:   "discover.provider.plex.tv",
         Path:   "/library/metadata/matches",
         RawQuery: url.Values{
            "url":          {path},
            "x-plex-token": {token},
         }.Encode(),
      },
      map[string]string{"accept": "application/json"},
   )
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   var result struct {
      Error struct {
         Message string
      }
      MediaContainer struct {
         Metadata []Metadata
      }
   }
   err = json.NewDecoder(resp.Body).Decode(&result)
   if err != nil {
      return nil, err
   }
   if result.Error.Message != "" {
      return nil, errors.New(result.Error.Message)
   }
   return &result.MediaContainer.Metadata[0], nil
}

func (m *Metadata) Fetch(token string) (*Metadata, error) {
   resp, err := maya.Get(
      &url.URL{
         Scheme: "https",
         Host:   "vod.provider.plex.tv",
         Path:   "/library/metadata/" + m.RatingKey,
      },
      map[string]string{
         "accept":       "application/json",
         "x-plex-token": token,
      },
   )
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   if resp.StatusCode != 200 {
      return nil, errors.New(resp.Status)
   }
   var result struct {
      MediaContainer struct {
         Metadata []Metadata
      }
   }
   err = json.NewDecoder(resp.Body).Decode(&result)
   if err != nil {
      return nil, err
   }
   return &result.MediaContainer.Metadata[0], nil
}

func (p *Part) FetchWidevine(token string, body []byte) ([]byte, error) {
   target, err := url.Parse(p.License)
   if err != nil {
      return nil, err
   }
   target.Scheme = "https"
   target.Host = "vod.provider.plex.tv"
   target.RawQuery = url.Values{
      "x-plex-drm":   {"widevine"},
      "x-plex-token": {token},
   }.Encode()
   resp, err := maya.Post(target, nil, body)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   return io.ReadAll(resp.Body)
}

func FetchUser() (*User, error) {
   resp, err := maya.Post(
      &url.URL{
         Scheme: "https", Host: "plex.tv", Path: "/api/v2/users/anonymous",
      },
      map[string]string{
         "accept":                   "application/json",
         "x-plex-product":           "Plex Mediaverse",
         "x-plex-client-identifier": "!",
      },
      nil,
   )
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   result := &User{}
   err = json.NewDecoder(resp.Body).Decode(result)
   if err != nil {
      return nil, err
   }
   return result, nil
}

type User struct {
   AuthToken string
}

type Metadata struct {
   Media []struct {
      Part     []Part
      Protocol string
   }
   RatingKey string
}

type Part struct {
   Key     string
   License string
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

///

func (m *Metadata) GetDash() (*Part, error) {
   for _, media := range m.Media {
      if media.Protocol == "dash" {
         // Success: Return the part and a nil error.
         // This will panic if media.Part is empty, matching the
         // behavior of your original function.
         return &media.Part[0], nil
      }
   }
   // Failure: No "dash" protocol was found.
   return nil, errors.New("DASH media part not found")
}

func (p *Part) GetManifest(token string) *url.URL {
   return &url.URL{
      Scheme:   "https",
      Host:     "vod.provider.plex.tv",
      Path:     p.Key, // /library/parts/6730016e43b96c02321d7860-dash.mpd
      RawQuery: url.Values{"x-plex-token": {token}}.Encode(),
   }
}
