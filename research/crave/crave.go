package crave

import (
   _ "embed"
   "net/url"
)

//go:embed GetShowpage.gql
var get_showpage string

func (s *Stream) GetManifest() (*url.URL, error) {
   return url.Parse(s.Playback)
}
