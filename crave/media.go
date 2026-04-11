/*
https://crave.ca/en/movie/anaconda-2025-59881
https://crave.ca/en/play/anaconda-2025-3300246
https://crave.ca/movie/anaconda-2025-59881
https://crave.ca/play/anaconda-2025-3300246
*/
package crave

import (
   "errors"
   "fmt"
   "net/url"
   "strconv"
   "strings"
)

type Media struct {
   FirstContent struct {
      Id int `json:"id,string"`
   }
   Id int `json:"id,string"`
}

func ParseMedia(rawURL string) (*Media, error) {
   parsedURL, err := url.Parse(rawURL)
   if err != nil {
      return nil, err
   }
   // Split the path directly.
   // e.g., "/en/movie/anaconda-2025-59881" -> ["", "en", "movie", "anaconda-2025-59881"]
   parts := strings.Split(parsedURL.Path, "/")
   // We need at least 3 parts: the empty string (before the first "/"), the type, and the slug
   if len(parts) < 3 {
      return nil, errors.New("invalid URL path format")
   }
   // Safely grab the last two segments
   lastPart := parts[len(parts)-1] // e.g., "anaconda-2025-59881"
   typePart := parts[len(parts)-2] // e.g., "movie" or "play"
   // Find the last dash to extract the ID
   dashIdx := strings.LastIndex(lastPart, "-")
   if dashIdx == -1 || dashIdx == len(lastPart)-1 {
      return nil, errors.New("no ID found at the end of the URL")
   }
   idStr := lastPart[dashIdx+1:]
   // Convert extracted string to integer
   id, err := strconv.Atoi(idStr)
   if err != nil {
      return nil, fmt.Errorf("invalid ID format: %w", err)
   }
   // Populate struct based on the type
   media := &Media{}
   switch typePart {
   case "movie":
      media.Id = id
   case "play":
      media.FirstContent.Id = id
   default:
      return nil, fmt.Errorf("unknown media type: %s", typePart)
   }
   return media, nil
}
