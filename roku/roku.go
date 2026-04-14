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

func ParseDash(playbackUrl string) (*url.URL, error) {
   return url.Parse(playbackUrl)
}

type Playback struct {
   Drm struct {
      Widevine struct {
         LicenseServer string
      }
   }
   Url string // MPD
}

const user_agent = "trc-googletv; production; 0"

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

// codeToken can be empty
//
// /api/v1/account/token
func FetchToken(codeToken string) (*Token, error) {
   req := http.Request{
      URL: &url.URL{
         Scheme: "https",
         Host:   "googletv.web.roku.com",
         Path:   "/api/v1/account/token",
      },
      Header: http.Header{},
   }
   req.Header.Set("user-agent", user_agent)
   if codeToken != "" {
      req.Header.Set("x-roku-content-token", codeToken)
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

// /api/v3/playback
func FetchPlayback(authToken, rokuId string) (*Playback, error) {
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
   req.Header.Set("x-roku-content-token", authToken)
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

func FetchWidevine(licenseServer string, data []byte) ([]byte, error) {
   resp, err := http.Post(
      licenseServer, "application/x-protobuf", bytes.NewReader(data),
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

// /api/v1/account/activation
func FetchActivation(authToken string) (*Activation, error) {
   body, err := json.Marshal(map[string]string{"platform": "googletv"})
   if err != nil {
      return nil, err
   }
   req, err := http.NewRequest(
      "POST", "https://googletv.web.roku.com/api/v1/account/activation",
      bytes.NewReader(body),
   )
   if err != nil {
      return nil, err
   }
   req.Header.Set("content-type", "application/json")
   req.Header.Set("user-agent", user_agent)
   req.Header.Set("x-roku-content-token", authToken)
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

type Token struct {
   AuthToken string
}

type Activation struct {
   Code string
}

type Code struct {
   Token string
}

// /api/v1/account/activation/code
func FetchCode(authToken, activationCode string) (*Code, error) {
   req := http.Request{
      URL: &url.URL{
         Scheme: "https",
         Host:   "googletv.web.roku.com",
         Path:   "/api/v1/account/activation/" + activationCode,
      },
      Header: http.Header{},
   }
   req.Header.Set("user-agent", user_agent)
   req.Header.Set("x-roku-content-token", authToken)
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
