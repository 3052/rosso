package roku

import (
   "41.neocities.org/maya"
   "encoding/json"
   "io"
   "net/url"
   "strings"
)

type AccountToken struct {
   AuthToken  string `json:"authToken"`
   IsLoggedIn bool   `json:"isLoggedIn"`
   Ip         string `json:"ip"`
   Rida       string `json:"rida"`
}

// status can be nil
func GetAccountToken(status *ActivationStatus) (*AccountToken, error) {
   target := &url.URL{
      Scheme: "https",
      Host:   "googletv.web.roku.com",
      Path:   "/api/v1/account/token",
   }
   headers := map[string]string{
      "user-agent": "trc-googletv; production; 0",
   }
   if status != nil {
      headers["x-roku-content-token"] = status.Token
   }

   resp, err := maya.Get(target, headers)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()

   var token AccountToken
   if err := json.NewDecoder(resp.Body).Decode(&token); err != nil {
      return nil, err
   }
   return &token, nil
}

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
