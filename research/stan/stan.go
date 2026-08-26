package stan

import (
   "bytes"
   "encoding/json"
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

type ActivationCode struct {
   Code string
   Url  string
}

func FetchActivationCode() (*ActivationCode, error) {
   resp, err := http.PostForm(
      "https://api.stan.com.au/login/v1/activation-codes/", url.Values{
         "generate": {"true"},
      },
   )
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
   resp, err := http.Get(a.Url)
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

func (a *AppSession) Stream(id int64) (*ProgramStream, error) {
   req, err := http.NewRequest(
      "GET", "https://api.stan.com.au/concurrency/v1/streams", nil,
   )
   if err != nil {
      return nil, err
   }
   req.Header.Set("x-forwarded-for", "1.128.0.0")
   req.URL.RawQuery = url.Values{
      "drm":       {"widevine"}, // need for .Media.DRM
      "format":    {"dash"},     // 404 otherwise
      "jwToken":   {a.JwToken},
      "programId": {strconv.FormatInt(id, 10)},
      "quality":   {"auto"}, // note `high` or `ultra` should work too
   }.Encode()
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   result := &ProgramStream{}
   err = json.NewDecoder(resp.Body).Decode(result)
   if err != nil {
      return nil, err
   }
   return result, nil
}

type LegacyProgram struct {
   ReleaseYear           int
   SeriesTitle           string
   Title                 string
   TvSeasonEpisodeNumber int
   TvSeasonNumber        int
}

func (p *LegacyProgram) New(id int64) error {
   address := func() string {
      b := []byte("https://api.stan.com.au/programs/v1/legacy/programs/")
      b = strconv.AppendInt(b, id, 10)
      b = append(b, ".json"...)
      return string(b)
   }()
   resp, err := http.Get(address)
   if err != nil {
      return err
   }
   defer resp.Body.Close()
   return json.NewDecoder(resp.Body).Decode(p)
}

type ProgramStream struct {
   Media struct {
      DRM *struct {
         CustomData string
         KeyId      string
      }
      VideoUrl string
   }
}

func (p ProgramStream) BaseUrl(host string) (*url.URL, error) {
   video, err := url.Parse(p.Media.VideoUrl)
   if err != nil {
      return nil, err
   }
   video.Host = host
   return video, nil
}

func (p *ProgramStream) License(data []byte) ([]byte, error) {
   req, err := http.NewRequest(
      // final slash is needed
      "POST", "https://lic.drmtoday.com/license-proxy-widevine/cenc/",
      bytes.NewReader(data),
   )
   if err != nil {
      return nil, err
   }
   req.Header.Set("dt-custom-data", p.Media.DRM.CustomData)
   resp, err := http.DefaultClient.Do(req)
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

func (w *WebToken) Session() (*AppSession, error) {
   resp, err := http.PostForm(
      "https://api.stan.com.au/login/v1/sessions/mobile/app", url.Values{
         "jwToken": {w.JwToken},
      },
   )
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
