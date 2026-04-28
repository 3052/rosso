package roku

import (
   "net/url"
   "strings"
)

func (p *PlaybackConfig) GetManifest() (*url.URL, error) {
   return url.Parse(p.Url)
}

func (a *AccountActivation) String() string {
   var data strings.Builder
   data.WriteString("1 Visit the URL\n")
   data.WriteString("\ttherokuchannel.com/link\n")
   data.WriteString("2 Enter the activation code\n")
   data.WriteByte('\t')
   data.WriteString(a.Code)
   return data.String()
}
