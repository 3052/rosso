package mubi

import (
   "fmt"
   "net/url"
   "strings"
)

func (s *SecureUrl) GetManifest() (*url.URL, error) {
   s.Url = strings.NewReplacer(
      ".AVC1", "",
      ".ex-eac3", "",
      ".ex-vtt", "",
   ).Replace(s.Url)
   return url.Parse(s.Url)
}

type SecureUrl struct {
   TextTrackUrls []struct {
      Id  string
      Url string
   } `json:"text_track_urls"`
   Url         string // MPD
   UserMessage string `json:"user_message"`
}

type LinkCode struct {
   AuthToken string `json:"auth_token"`
   LinkCode  string `json:"link_code"`
}

func (l *LinkCode) String() string {
   var data strings.Builder
   data.WriteString("TO LOG IN AND START WATCHING\n")
   data.WriteString("Go to\n")
   data.WriteString("mubi.com/android\n")
   data.WriteString("and enter the code below\n")
   data.WriteString(l.LinkCode)
   return data.String()
}

// "android" requires headers:
// client-device-identifier
// client-version
const client = "web"

var ClientCountry = "US"

func (f *Film) String() string {
   data := &strings.Builder{}
   data.WriteString("title = ")
   data.WriteString(f.Title)
   data.WriteString("\nid = ")
   fmt.Fprint(data, f.Id)
   return data.String()
}

type Film struct {
   Title string
   Id    int
}

type Session struct {
   Token string
   User  struct {
      Id int
   }
}
