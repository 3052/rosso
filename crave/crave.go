package crave

import (
   _ "embed"
   "errors"
   "fmt"
   "net/url"
   "strconv"
   "strings"
)

func (s *Subscription) String() string {
   var data strings.Builder
   data.WriteString("display name = ")
   data.WriteString(s.Experience.DisplayName)
   data.WriteString("\nexpiration date = ")
   data.WriteString(s.ExpirationDate)
   return data.String()
}

/*
https://crave.ca/en/movie/anaconda-2025-59881
https://crave.ca/en/play/anaconda-2025-3300246
https://crave.ca/movie/anaconda-2025-59881
https://crave.ca/play/anaconda-2025-3300246
https://crave.ca/play/heated-rivalry/ill-believe-in-anything-s1e5-3233873
*/
func ParseMedia(rawUrl string) (*Media, error) {
   parsedUrl, err := url.Parse(rawUrl)
   if err != nil {
      return nil, err
   }
   // Split the path directly.
   parts := strings.Split(parsedUrl.Path, "/")
   if len(parts) < 3 {
      return nil, errors.New("invalid URL path format")
   }
   // Anchor the URL by looking for the explicit media type
   var typePart string
   for _, part := range parts {
      if part == "movie" || part == "play" {
         typePart = part
         break
      }
   }
   if typePart == "" {
      return nil, errors.New("missing media type (movie/play) in URL")
   }
   // Safely grab the last segment (the slug containing the ID)
   lastPart := parts[len(parts)-1]
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
   // Populate struct based on the anchored type
   media_data := &Media{}
   switch typePart {
   case "movie":
      media_data.Id = id
   case "play":
      media_data.FirstContent.Id = id
   }
   return media_data, nil
}

func (p *Profile) String() string {
   var data strings.Builder
   data.WriteString("nickname = ")
   data.WriteString(p.Nickname)
   if p.HasPin {
      data.WriteString("\nhas pin = true")
   } else {
      data.WriteString("\nhas pin = false")
   }
   if p.Master {
      data.WriteString("\nmaster = true")
   } else {
      data.WriteString("\nmaster = false")
   }
   data.WriteString("\nmaturity = ")
   data.WriteString(p.Maturity)
   data.WriteString("\nid = ")
   data.WriteString(p.Id)
   return data.String()
}

//go:embed GetShowpage.gql
var get_showpage string

func (s *Stream) GetManifest() (*url.URL, error) {
   return url.Parse(s.Playback)
}
