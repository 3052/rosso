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
   "net/http"
   "net/url"
   "strconv"
   "strings"
   "time"
)

func (s Stream) GetManifest() (*url.URL, error) {
   return url.Parse(strings.Replace(s.PlaybackUrl, "_2sec", "", 1))
}

// https://nbc.com/saturday-night-live/video/november-15-glen-powell/9000454161
func GetName(urlData string) (string, error) {
   url_parse, err := url.Parse(urlData)
   if err != nil {
      return "", err
   }
   return strings.TrimPrefix(url_parse.Path, "/"), nil
}

func FetchMetadata(name string) (*Metadata, error) {
   data, err := json.Marshal(map[string]any{
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
   resp, err := http.Post(
      "https://friendship.nbc.com/v3/graphql", "application/json",
      bytes.NewReader(data),
   )
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   if resp.StatusCode != http.StatusOK {
      return nil, errors.New(resp.Status)
   }
   var result struct {
      Data struct {
         Page struct {
            Metadata Metadata
         }
      }
      Errors []struct {
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
   return &result.Data.Page.Metadata, nil
}

type Metadata struct {
   MpxAccountId    int `json:",string"`
   MpxGuid         int `json:",string"`
   ProgrammingType string
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

//go:embed page.gql
var query_page string

func (m *Metadata) Stream() (*Stream, error) {
   req := http.Request{
      URL: &url.URL{
         Scheme: "https",
         Host:   "lemonade.nbc.com",
         Path:   fmt.Sprintf("/v1/vod/%v/%v", m.MpxAccountId, m.MpxGuid),
         RawQuery: url.Values{
            "platform":        {"web"},
            "programmingType": {m.ProgrammingType},
         }.Encode(),
      },
      Header: http.Header{},
   }
   resp, err := http.DefaultClient.Do(&req)
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

const drmProxySecret = "Whn8QFuLFM7Heiz6fYCYga7cYPM8ARe6"

func playReady() *url.URL {
   return &url.URL{
      Scheme:   "https",
      Host:     "drmproxy.digitalsvc.apps.nbcuni.com",
      Path:     "/drm-proxy/license/playready",
      RawQuery: build_query("playready"),
   }
}

type Stream struct {
   PlaybackUrl string // MPD
}
