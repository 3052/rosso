package cineMember

import (
   "errors"
   "net/url"
)

func (m *MediaLink) GetManifest() (*url.URL, error) {
   return url.Parse(m.Url)
}

type Stream struct {
   Error    string
   Links    []MediaLink
   NoAccess bool
}

func (s *Stream) Dash() (*MediaLink, error) {
   for _, link := range s.Links {
      if link.MimeType == "application/dash+xml" {
         return &link, nil
      }
   }
   return nil, errors.New("DASH link not found")
}

type MediaLink struct {
   MimeType string
   Url      string
}
