package stan

import (
   "bytes"
   "encoding/json"
   "log"
   "net/http"
   "net/url"
   "strconv"
   "strings"
)

var BaseUrl = []string{
   "023-stan.akamaized.net",
   "666-stan.akamaized.net", // geo block
   "aws.stan.video",
   "gec.stan.video",
}

func do(req *http.Request) (*http.Response, error) {
   log.Println(req.Method, req.URL)
   return http.DefaultClient.Do(req)
}

type ActivationCode struct {
   Code string
   Url  string
}

func FetchActivationCode() (*ActivationCode, error) {
   req, err := http.NewRequest(
      "POST", "https://api.stan.com.au/login/v1/activation-codes/",
      strings.NewReader(url.Values{
         "generate": {"true"},
      }.Encode()),
   )
   if err != nil {
      return nil, err
   }
   req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
   resp, err := do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   var code ActivationCode
   err = json.NewDecoder(resp.Body).Decode(&code)
   if err != nil {
      return nil, err
   }
   return &code, nil
}

func (*ActivationCode) CachePath() string {
   return "rosso/stan/ActivationCode"
}

func (a *ActivationCode) String() string {
   var data strings.Builder
   data.WriteString("Stan.\n")
   data.WriteString("Log in with code\n")
   data.WriteString("1. Visit stan.com.au/activate\n")
   data.WriteString("2. Enter the code:\n")
   data.WriteString(a.Code)
   return data.String()
}

func (a *ActivationCode) Token() (*WebToken, error) {
   req, err := http.NewRequest("GET", a.Url, nil)
   if err != nil {
      return nil, err
   }
   resp, err := do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   var web WebToken
   err = json.NewDecoder(resp.Body).Decode(&web)
   if err != nil {
      return nil, err
   }
   return &web, nil
}

type AppSession struct {
   JwToken string
}

func (a *AppSession) FetchMedia(id int) (*Media, error) {
   req, err := http.NewRequest(
      "GET", "https://api.stan.com.au/concurrency/v1/streams", nil,
   )
   if err != nil {
      return nil, err
   }
   req.Header.Set("x-forwarded-for", "1.128.0.0")
   req.URL.RawQuery = url.Values{
      "drm":       {"widevine"}, // need for .Media.Drm
      "format":    {"dash"},     // 404 otherwise
      "jwToken":   {a.JwToken},
      "programId": {strconv.Itoa(id)},
      "quality":   {"auto"}, // note `high` or `ultra` should work too
   }.Encode()
   resp, err := do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   var result struct {
      Media Media
   }
   err = json.NewDecoder(resp.Body).Decode(&result)
   if err != nil {
      return nil, err
   }
   return &result.Media, nil
}

type Media struct {
   Drm *struct {
      CustomData string
      KeyId      string
   }
   VideoUrl string
}

func (m *Media) BaseUrl(host string) (*url.URL, error) {
   video, err := url.Parse(m.VideoUrl)
   if err != nil {
      return nil, err
   }
   video.Host = host
   return video, nil
}

func (*Media) CachePath() string {
   return "rosso/stan/Media"
}

func (m *Media) License(data []byte) ([]byte, error) {
   req, err := http.NewRequest(
      // final slash is needed
      "POST", "https://lic.drmtoday.com/license-proxy-widevine/cenc/",
      bytes.NewReader(data),
   )
   if err != nil {
      return nil, err
   }
   req.Header.Set("dt-custom-data", m.Drm.CustomData)
   resp, err := do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   var result struct {
      License []byte
   }
   err = json.NewDecoder(resp.Body).Decode(&result)
   if err != nil {
      return nil, err
   }
   return result.License, nil
}

type WebToken struct {
   JwToken   string
   ProfileId string
}

func (*WebToken) CachePath() string {
   return "rosso/stan/WebToken"
}

func (w *WebToken) FetchSession() (*AppSession, error) {
   req, err := http.NewRequest(
      "POST", "https://api.stan.com.au/login/v1/sessions/mobile/app",
      strings.NewReader(url.Values{
         "jwToken": {w.JwToken},
      }.Encode()),
   )
   if err != nil {
      return nil, err
   }
   req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
   resp, err := do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   result := &AppSession{}
   err = json.NewDecoder(resp.Body).Decode(result)
   if err != nil {
      return nil, err
   }
   return result, nil
}

// stan.go
