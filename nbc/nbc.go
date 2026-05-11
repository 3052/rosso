package nbc

import (
   "41.neocities.org/maya"
   "crypto/hmac"
   "crypto/sha256"
   _ "embed"
   "encoding/hex"
   "encoding/json"
   "errors"
   "fmt"
   "io"
   "net/url"
   "strconv"
   "strings"
   "time"
)

type Url struct {
   Url url.URL
}

func (u *Url) UnmarshalText(text []byte) error {
   return u.Url.UnmarshalBinary(text)
}

func (u *Url) MarshalText() ([]byte, error) {
   return u.Url.MarshalBinary()
}

type Stream struct {
   PlaybackUrl *Url // MPD
}

func (s Stream) GetManifest() *url.URL {
   manifest := s.PlaybackUrl.Url
   manifest.Path = strings.Replace(manifest.Path, "_2sec", "", 1)
   return &manifest
}

func (m *Metadata) FetchStream() (*Stream, error) {
   resp, err := maya.Get(
      &url.URL{
         Scheme: "https",
         Host:   "lemonade.nbc.com",
         Path:   fmt.Sprintf("/v1/vod/%v/%v", m.MpxAccountId, m.MpxGuid),
         RawQuery: url.Values{
            "platform":        {"web"},
            "programmingType": {m.ProgrammingType},
         }.Encode(),
      },
      nil,
   )
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   if resp.StatusCode != 200 {
      return nil, errors.New(resp.Status)
   }
   result := &Stream{}
   err = json.NewDecoder(resp.Body).Decode(result)
   if err != nil {
      return nil, err
   }
   return result, nil
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
   resp, err := maya.Post(
      &url.URL{
         Scheme: "https",
         Host:   "friendship.nbc.com",
         Path:   "/v3/graphql",
      },
      map[string]string{"content-type": "application/json"},
      body,
   )
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   if resp.StatusCode != 200 {
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

func FetchWidevine(body []byte) ([]byte, error) {
   resp, err := maya.Post(
      &url.URL{
         Scheme:   "https",
         Host:     "drmproxy.digitalsvc.apps.nbcuni.com",
         Path:     "/drm-proxy/license/widevine",
         RawQuery: build_query("widevine"),
      },
      map[string]string{"content-type": "application/octet-stream"},
      body,
   )
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   return io.ReadAll(resp.Body)
}

// https://nbc.com/saturday-night-live/video/november-15-glen-powell/9000454161
func GetName(urlData string) (string, error) {
   url_parse, err := url.Parse(urlData)
   if err != nil {
      return "", err
   }
   return strings.TrimPrefix(url_parse.Path, "/"), nil
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

const drmProxySecret = "Whn8QFuLFM7Heiz6fYCYga7cYPM8ARe6"

func playReady() *url.URL {
   return &url.URL{
      Scheme:   "https",
      Host:     "drmproxy.digitalsvc.apps.nbcuni.com",
      Path:     "/drm-proxy/license/playready",
      RawQuery: build_query("playready"),
   }
}
