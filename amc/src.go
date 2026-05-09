package amc

import "net/url"

type Source struct {
   Codecs     string
   KeySystems KeySystems `json:"key_systems"`
   Src        Src        // MPD
   Type       string
}

type Src struct {
   Url url.URL
}

func (s *Src) UnmarshalText(text []byte) error {
   return s.Url.UnmarshalBinary(text)
}

func (s *Src) MarshalText() ([]byte, error) {
   return s.Url.MarshalBinary()
}
