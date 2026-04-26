package plex

import (
   "41.neocities.org/maya"
   "encoding/json"
   "errors"
   "io"
   "net/url"
   "strings"
)

func AcquireWidevineLicense(media *VodMedia, anonymous *AnonymousUser, body []byte) ([]byte, error) {
   if len(media.Part) == 0 {
      return nil, errors.New("no media parts found")
   }
   if media.Part[0].License == "" {
      return nil, errors.New("no license path found")
   }

   endpoint := &url.URL{
      Scheme: "https",
      Host:   "vod.provider.plex.tv",
      Path:   media.Part[0].License,
   }

   query := url.Values{}
   query.Set("x-plex-drm", "widevine")
   query.Set("x-plex-token", anonymous.AuthToken)
   endpoint.RawQuery = query.Encode()

   resp, err := maya.Post(endpoint, nil, body)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()

   return io.ReadAll(resp.Body)
}

type AnonymousUser struct {
   Id        int    `json:"id"`
   Uuid      string `json:"uuid"`
   AuthToken string `json:"authToken"`
}

func CreateAnonymousUser() (*AnonymousUser, error) {
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

   var anonymous AnonymousUser
   if err := json.NewDecoder(resp.Body).Decode(&anonymous); err != nil {
      return nil, err
   }

   return &anonymous, nil
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
