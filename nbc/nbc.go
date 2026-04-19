package nbc

import (
   "crypto/hmac"
   "crypto/sha256"
   _ "embed"
   "encoding/hex"
   "io"
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

type Stream struct {
   PlaybackUrl string // MPD
}
