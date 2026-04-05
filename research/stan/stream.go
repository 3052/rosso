package stan

import (
   "encoding/json"
   "errors"
   "net/http"
   "net/url"
   "strconv"
   "strings"
)

func (a AppSession) Stream(id int64) (*ProgramStream, error) {
   req, err := http.NewRequest(
      "GET", "https://api.stan.com.au/concurrency/v1/streams", nil,
   )
   if err != nil {
      return nil, err
   }
   req.Header.Set("x-forwarded-for", "1.128.0.0")
   req.URL.RawQuery = url.Values{
      "drm": {"widevine"}, // need for .Media.DRM
      "format": {"dash"}, // 404 otherwise
      "jwToken": {a.JwToken},
      "programId": {strconv.FormatInt(id, 10)},
      "quality": {"auto"}, // note `high` or `ultra` should work too
   }.Encode()
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   if resp.StatusCode != http.StatusOK {
      var b strings.Builder
      resp.Write(&b)
      return nil, errors.New(b.String())
   }
   stream := new(ProgramStream)
   err = json.NewDecoder(resp.Body).Decode(stream)
   if err != nil {
      return nil, err
   }
   return stream, nil
}

type ProgramStream struct {
   Media struct {
      DRM *struct {
         CustomData string
         KeyId string
      }
      VideoUrl string
   }
}
func (ProgramStream) WrapRequest(b []byte) ([]byte, error) {
   return b, nil
}

func (p ProgramStream) RequestHeader() (http.Header, error) {
   head := make(http.Header)
   head.Set("dt-custom-data", p.Media.DRM.CustomData)
   return head, nil
}

// final slash is needed
func (ProgramStream) RequestUrl() (string, bool) {
   return "https://lic.drmtoday.com/license-proxy-widevine/cenc/", true
}

func (ProgramStream) UnwrapResponse(b []byte) ([]byte, error) {
   var s struct {
      License []byte
   }
   err := json.Unmarshal(b, &s)
   if err != nil {
      return nil, err
   }
   return s.License, nil
}

var BaseUrl = []string{
   "023-stan.akamaized.net",
   "666-stan.akamaized.net", // geo block
   "aws.stan.video",
   "gec.stan.video",
}

func (p ProgramStream) BaseUrl(host string) (*url.URL, error) {
   video, err := url.Parse(p.Media.VideoUrl)
   if err != nil {
      return nil, err
   }
   video.Host = host
   return video, nil
}
