package roku

import (
   "41.neocities.org/maya"
   "io"
   "net/url"
   "strings"
)

func (p *Playback) GetWidevineLicense(challenge []byte) ([]byte, error) {
   target, err := url.Parse(p.Drm.Widevine.LicenseServer)
   if err != nil {
      return nil, err
   }
   headers := map[string]string{
      "content-type": "application/x-protobuf",
      "user-agent":   "Go-http-client/2.0",
   }

   resp, err := maya.Post(target, headers, challenge)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()

   return io.ReadAll(resp.Body)
}

func (p *Playback) GetManifest() (*url.URL, error) {
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
