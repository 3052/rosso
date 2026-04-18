package plex

import (
   "41.neocities.org/maya"
   "encoding/json"
   "errors"
   "io"
   "net/url"
)

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
