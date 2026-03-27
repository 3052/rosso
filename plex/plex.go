package plex

import (
   "bytes"
   "encoding/json"
   "errors"
   "io"
   "net/http"
   "net/url"
   "strings"
)

func (u User) RatingKey(rawUrl string) (*ItemMetadata, error) {
   req := http.Request{
      URL: &url.URL{
         Scheme: "https",
         Host:   "discover.provider.plex.tv",
         Path:   "/library/metadata/matches",
         RawQuery: url.Values{
            "url":          {rawUrl},
            "x-plex-token": {u.AuthToken},
         }.Encode(),
      },
      Header: http.Header{},
   }
   req.Header.Set("accept", "application/json")
   resp, err := http.DefaultClient.Do(&req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   var result struct {
      Error struct {
         Message string
      }
      MediaContainer struct {
         Metadata []ItemMetadata
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

func (u User) Media(item *ItemMetadata, forwardedFor string) (*ItemMetadata, error) {
   var req http.Request
   req.Header = http.Header{}
   req.Header.Set("accept", "application/json")
   req.Header.Set("x-plex-token", u.AuthToken)
   if forwardedFor != "" {
      req.Header.Set("X-Forwarded-For", forwardedFor)
   }
   req.URL = &url.URL{
      Scheme: "https",
      Host:   "vod.provider.plex.tv",
      Path:   "/library/metadata/" + item.RatingKey,
   }
   resp, err := http.DefaultClient.Do(&req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   if resp.StatusCode != http.StatusOK {
      return nil, errors.New(resp.Status)
   }
   var result struct {
      MediaContainer struct {
         Metadata []ItemMetadata
      }
   }
   err = json.NewDecoder(resp.Body).Decode(&result)
   if err != nil {
      return nil, err
   }
   return &result.MediaContainer.Metadata[0], nil
}

type Dash struct {
   Body []byte
   Url  *url.URL
}

// https://watch.plex.tv/embed/movie/memento-2000
// https://watch.plex.tv/movie/memento-2000
// https://watch.plex.tv/watch/movie/memento-2000
func GetPath(rawUrl string) (string, error) {
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

type ItemMetadata struct {
   Media []struct {
      Part     []MediaPart
      Protocol string
   }
   RatingKey string
}

type MediaPart struct {
   Key     string
   License string
}

func (i *ItemMetadata) Dash() (*MediaPart, error) {
   for _, rosso := range i.Media {
      if rosso.Protocol == "dash" {
         // Success: Return the part and a nil error.
         // This will panic if media.Part is empty, matching the
         // behavior of your original function.
         return &rosso.Part[0], nil
      }
   }
   // Failure: No "dash" protocol was found.
   return nil, errors.New("DASH media part not found")
}
func FetchUser() (*User, error) {
   req := http.Request{
      Method: "POST",
      URL: &url.URL{
         Scheme: "https",
         Host:   "plex.tv",
         Path:   "/api/v2/users/anonymous",
      },
      Header: http.Header{},
   }
   req.Header.Set("accept", "application/json")
   req.Header.Set("x-plex-product", "Plex Mediaverse")
   req.Header.Set("x-plex-client-identifier", "!")
   resp, err := http.DefaultClient.Do(&req)
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
func (u User) Dash(part *MediaPart, forwardedFor string) (*Dash, error) {
   req := http.Request{
      URL: &url.URL{
         Scheme:   "https",
         Host:     "vod.provider.plex.tv",
         Path:     part.Key, // /library/parts/6730016e43b96c02321d7860-dash.mpd
         RawQuery: url.Values{"x-plex-token": {u.AuthToken}}.Encode(),
      },
      Header: http.Header{},
   }
   if forwardedFor != "" {
      req.Header.Set("X-Forwarded-For", forwardedFor)
   }
   resp, err := http.DefaultClient.Do(&req)
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

type User struct {
   AuthToken string
}

func (u User) Widevine(part *MediaPart, data []byte) ([]byte, error) {
   req, err := http.NewRequest("POST", part.License, bytes.NewReader(data))
   if err != nil {
      return nil, err
   }
   req.URL.Scheme = "https"
   req.URL.Host = "vod.provider.plex.tv"
   req.URL.RawQuery = url.Values{
      "x-plex-drm":   {"widevine"},
      "x-plex-token": {u.AuthToken},
   }.Encode()
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   return io.ReadAll(resp.Body)
}
