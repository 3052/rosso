package roku

import (
   "net/url"
   "strings"
)

func (p *Playback) GetManifest() (*url.URL, error) {
   return url.Parse(p.Url)
}

type Playback struct {
   Drm struct {
      Widevine struct {
         LicenseServer string
      }
   }
   Url string // MPD
}

func FormatActivation(code string) string {
   var data strings.Builder
   data.WriteString("1 Visit the URL\n")
   data.WriteString("  therokuchannel.com/link\n")
   data.WriteString("\n")
   data.WriteString("2 Enter the activation code\n")
   data.WriteString("  ")
   data.WriteString(code)
   return data.String()
}

type Token struct {
   AuthToken string
}

type Activation struct {
   Code string
}

type Code struct {
   Token string
}

const user_agent = "trc-googletv; production; 0"
