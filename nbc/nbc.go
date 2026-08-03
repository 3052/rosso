package nbc

import (
   "bytes"
   "crypto/hmac"
   "crypto/sha256"
   _ "embed"
   "encoding/hex"
   "encoding/json"
   "errors"
   "fmt"
   "io"
   "log"
   "net/http"
   "net/url"
   "strconv"
   "strings"
   "time"
)

const drmProxySecret = "Whn8QFuLFM7Heiz6fYCYga7cYPM8ARe6"

//go:embed page.gql
var query_page string

func FetchWidevine(body []byte) ([]byte, error) {
   req, err := http.NewRequest(
      "POST",
      (&url.URL{
         Scheme:   "https",
         Host:     "drmproxy.digitalsvc.apps.nbcuni.com",
         Path:     "/drm-proxy/license/widevine",
         RawQuery: build_query("widevine"),
      }).String(),
      bytes.NewReader(body),
   )
   if err != nil {
      return nil, err
   }
   req.Header.Set("content-type", "application/octet-stream")
   resp, err := do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   return io.ReadAll(resp.Body)
}

// https://nbc.com/saturday-night-live/video/november-15-glen-powell/9000454161
func GetName(urlData string) (string, error) {
   parse, err := url.Parse(urlData)
   if err != nil {
      return "", err
   }
   return strings.TrimPrefix(parse.Path, "/"), nil
}

func build_query(drmType string) string {
   timestamp := strconv.FormatInt(time.Now().UnixMilli(), 10)
   mac := hmac.New(sha256.New, []byte(drmProxySecret))
   // Use io.WriteString to write string data directly to the Writer
   io.WriteString(mac, timestamp)
   io.WriteString(mac, drmType)
   hash := hex.EncodeToString(mac.Sum(nil))
   return url.Values{
      "device": {"web"},
      "hash":   {hash},
      "time":   {timestamp},
   }.Encode()
}

// do sends the request, logs the method and URL, and returns the response.
func do(req *http.Request) (*http.Response, error) {
   log.Println(req.Method, req.URL)
   return http.DefaultClient.Do(req)
}

func playReady() *url.URL {
   return &url.URL{
      Scheme:   "https",
      Host:     "drmproxy.digitalsvc.apps.nbcuni.com",
      Path:     "/drm-proxy/license/playready",
      RawQuery: build_query("playready"),
   }
}

type Metadata struct {
   MpxAccountId    int `json:",string"`
   MpxGuid         int `json:",string"`
   ProgrammingType string
}

func FetchMetadata(name string) (*Metadata, error) {
   body, err := json.Marshal(map[string]any{
      "query": query_page,
      "variables": map[string]string{
         "app":      "nbc",
         "name":     name,
         "platform": "web",
         "type":     "VIDEO",
         "userId":   "",
      },
   })
   if err != nil {
      return nil, err
   }
   req, err := http.NewRequest(
      "POST",
      (&url.URL{
         Scheme: "https",
         Host:   "friendship.nbc.com",
         Path:   "/v3/graphql",
      }).String(),
      bytes.NewReader(body),
   )
   if err != nil {
      return nil, err
   }
   req.Header.Set("content-type", "application/json")
   resp, err := do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   if resp.StatusCode != http.StatusOK {
      return nil, errors.New(resp.Status)
   }
   var result struct {
      Data *struct {
         Page *struct {
            Metadata *Metadata
         }
      }
      Errors []*struct {
         Message string
      }
   }
   err = json.NewDecoder(resp.Body).Decode(&result)
   if err != nil {
      return nil, err
   }
   if len(result.Errors) >= 1 {
      return nil, errors.New(result.Errors[0].Message)
   }
   return result.Data.Page.Metadata, nil
}

func (m *Metadata) FetchStream() (*Stream, error) {
   req, err := http.NewRequest(
      "GET",
      (&url.URL{
         Scheme: "https",
         Host:   "lemonade.nbc.com",
         Path:   fmt.Sprintf("/v1/vod/%v/%v", m.MpxAccountId, m.MpxGuid),
         RawQuery: url.Values{
            "platform":        {"web"},
            "programmingType": {m.ProgrammingType},
         }.Encode(),
      }).String(),
      nil,
   )
   if err != nil {
      return nil, err
   }
   resp, err := do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   if resp.StatusCode != http.StatusOK {
      return nil, errors.New(resp.Status)
   }
   result := &Stream{}
   err = json.NewDecoder(resp.Body).Decode(result)
   if err != nil {
      return nil, err
   }
   return result, nil
}

type Stream struct {
   PlaybackUrl string // MPD
}

func (s *Stream) GetManifest() (*url.URL, error) {
   manifest, err := url.Parse(s.PlaybackUrl)
   if err != nil {
      return nil, err
   }
   manifest.Path = strings.Replace(manifest.Path, "_2sec", "", 1)
   return manifest, nil
}

// nbc.go
