package roku

import (
   "bytes"
   "encoding/json"
   "errors"
   "io"
   "net/http"
   "net/url"
   "strings"
)

func (p *Playback) Dash() (*Dash, error) {
   resp, err := http.Get(p.Url)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   body, err := io.ReadAll(resp.Body)
   if err != nil {
      return nil, err
   }
   return &Dash{Body: body, Url: resp.Request.URL}, nil
}

// /api/v3/playback
func (t *Token) Playback(rokuId string) (*Playback, error) {
   data, err := json.Marshal(map[string]string{
      "mediaFormat": "DASH",
      "providerId":  "rokuavod",
      "rokuId":      rokuId,
   })
   if err != nil {
      return nil, err
   }
   req, err := http.NewRequest(
      "POST", "https://googletv.web.roku.com/api/v3/playback",
      bytes.NewReader(data),
   )
   if err != nil {
      return nil, err
   }
   req.Header.Set("content-type", "application/json")
   req.Header.Set("user-agent", user_agent)
   req.Header.Set("x-roku-content-token", t.AuthToken)
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   if resp.StatusCode != http.StatusOK {
      return nil, errors.New(resp.Status)
   }
   result := &Playback{}
   err = json.NewDecoder(resp.Body).Decode(result)
   if err != nil {
      return nil, err
   }
   return result, nil
}

// /api/v1/account/activation
func (t *Token) Activation() (*Activation, error) {
   data, err := json.Marshal(map[string]string{"platform": "googletv"})
   if err != nil {
      return nil, err
   }
   req, err := http.NewRequest(
      "POST", "https://googletv.web.roku.com/api/v1/account/activation",
      bytes.NewReader(data),
   )
   if err != nil {
      return nil, err
   }
   req.Header.Set("content-type", "application/json")
   req.Header.Set("user-agent", user_agent)
   req.Header.Set("x-roku-content-token", t.AuthToken)
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   result := &Activation{}
   err = json.NewDecoder(resp.Body).Decode(result)
   if err != nil {
      return nil, err
   }
   return result, nil
}

// input can be nil
//
// /api/v1/account/token
func FetchToken(codeData *Code) (*Token, error) {
   req := http.Request{
      URL: &url.URL{
         Scheme: "https",
         Host:   "googletv.web.roku.com",
         Path:   "/api/v1/account/token",
      },
      Header: http.Header{},
   }
   req.Header.Set("user-agent", user_agent)
   if codeData != nil {
      req.Header.Set("x-roku-content-token", codeData.Token)
   }
   resp, err := http.DefaultClient.Do(&req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   result := &Token{}
   err = json.NewDecoder(resp.Body).Decode(result)
   if err != nil {
      return nil, err
   }
   return result, nil
}
// /api/v1/account/activation/code
func (t *Token) Code(activationData *Activation) (*Code, error) {
   req := http.Request{
      URL: &url.URL{
         Scheme: "https",
         Host:   "googletv.web.roku.com",
         Path:   "/api/v1/account/activation/" + activationData.Code,
      },
      Header: http.Header{},
   }
   req.Header.Set("user-agent", user_agent)
   req.Header.Set("x-roku-content-token", t.AuthToken)
   resp, err := http.DefaultClient.Do(&req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   result := &Code{}
   err = json.NewDecoder(resp.Body).Decode(result)
   if err != nil {
      return nil, err
   }
   return result, nil
}

func (p *Playback) Widevine(data []byte) ([]byte, error) {
   resp, err := http.Post(
      p.Drm.Widevine.LicenseServer, "application/x-protobuf",
      bytes.NewReader(data),
   )
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   if resp.StatusCode != http.StatusOK {
      return nil, errors.New(resp.Status)
   }
   return io.ReadAll(resp.Body)
}

type Dash struct {
   Body []byte
   Url  *url.URL
}

const user_agent = "trc-googletv; production; 0"

type Token struct {
   AuthToken string
}

type Activation struct {
   Code string
}

type Code struct {
   Token string
}

func (a *Activation) String() string {
   var data strings.Builder
   data.WriteString("1 Visit the URL\n")
   data.WriteString("  therokuchannel.com/link\n")
   data.WriteString("\n")
   data.WriteString("2 Enter the activation code\n")
   data.WriteString("  ")
   data.WriteString(a.Code)
   return data.String()
}

type Playback struct {
   Drm struct {
      Widevine struct {
         LicenseServer string
      }
   }
   Url string // MPD
}
