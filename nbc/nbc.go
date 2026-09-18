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
var queryPage string

func FetchWidevine(body []byte) ([]byte, error) {
   req, err := http.NewRequest("POST", licenseURL("widevine"), bytes.NewReader(body))
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
   parsedURL, err := url.Parse(urlData)
   if err != nil {
      return "", err
   }
   return strings.TrimPrefix(parsedURL.Path, "/"), nil
}

func PlayReady() string {
   return licenseURL("playready")
}

func buildQuery(drmType string) string {
   timestamp := strconv.FormatInt(time.Now().UnixMilli(), 10)

   mac := hmac.New(sha256.New, []byte(drmProxySecret))
   fmt.Fprint(mac, timestamp, drmType)
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

// licenseURL returns the DRM-proxy licence URL for the given DRM type.
func licenseURL(drmType string) string {
   return "https://drmproxy.digitalsvc.apps.nbcuni.com/drm-proxy/license/" +
      drmType + "?" + buildQuery(drmType)
}

type Metadata struct {
   MpxAccountId    int `json:",string"`
   MpxGuid         int `json:",string"`
   ProgrammingType string
}

func FetchMetadata(name string) (*Metadata, error) {
   body, err := json.Marshal(map[string]any{
      "query": queryPage,
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

   req, err := http.NewRequest("POST", "https://friendship.nbc.com/v3/graphql", bytes.NewReader(body))
   if err != nil {
      return nil, err
   }
   req.Header.Set("content-type", "application/json")
   req.Header.Add("x-tve-platform", "web")

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
   if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
      return nil, err
   }
   if len(result.Errors) >= 1 {
      return nil, errors.New(result.Errors[0].Message)
   }

   return result.Data.Page.Metadata, nil
}

func (m *Metadata) FetchStream() (*Stream, error) {
   params := url.Values{
      "platform":        {"web"},
      "programmingType": {m.ProgrammingType},
   }
   target := fmt.Sprintf(
      "https://lemonade.nbc.com/v1/vod/%d/%d?%s",
      m.MpxAccountId, m.MpxGuid, params.Encode(),
   )

   req, err := http.NewRequest("GET", target, nil)
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
   if err := json.NewDecoder(resp.Body).Decode(result); err != nil {
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
