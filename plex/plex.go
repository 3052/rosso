package plex

import (
   "errors"
   "net/url"
   "strings"
)

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

func (vm *VodMetadata) GetDashMedia() (*VodMedia, error) {
   for _, media := range vm.Media {
      if media.Protocol == "dash" {
         return &media, nil
      }
   }
   
   return nil, errors.New("dash media not found")
}

func BuildManifestUrl(part *VodPart, authToken string) *url.URL {
   endpoint := &url.URL{
      Scheme: "https",
      Host:   "vod.provider.plex.tv",
      Path:   part.Key,
   }

   query := url.Values{}
   query.Set("x-plex-token", authToken)
   endpoint.RawQuery = query.Encode()

   return endpoint
}
