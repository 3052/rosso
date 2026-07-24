// peacock.go
package peacock

import (
   "bytes"
   "crypto/hmac"
   "crypto/md5"
   "crypto/sha1"
   "encoding/base64"
   "encoding/json"
   "errors"
   "fmt"
   "io"
   "log"
   "maps"
   "net/http"
   "net/url"
   "slices"
   "strings"
   "time"
)

const (
   sky_client  = "NBCU-ANDROID-v3"
   sky_key     = "JuLQgyFz9n89D9pxcN6ZWZXKWfgj2PNBUb32zybj"
   sky_version = "1.0"
)

func doRequest(req *http.Request) (*http.Response, error) {
   log.Println(req.Method, req.URL)
   return http.DefaultClient.Do(req)
}

func generate_sky_ott(method, path string, header map[string]string, body []byte) string {
   // Sort headers by key
   header_keys := slices.Sorted(maps.Keys(header))
   // Build the special headers string
   var headers bytes.Buffer
   for _, key := range header_keys {
      lowerKey := strings.ToLower(key)
      if strings.HasPrefix(lowerKey, "x-skyott-") {
         headers.WriteString(lowerKey)
         headers.WriteString(": ")
         headers.WriteString(header[key])
         headers.WriteByte('\n')
      }
   }
   // MD5 the headers string and request body.
   headersHash := md5.Sum(headers.Bytes())
   headersMD5 := fmt.Sprintf("%x", headersHash)
   bodyHash := md5.Sum(body)
   bodyMD5 := fmt.Sprintf("%x", bodyHash)
   // Get current timestamp string directly.
   timestampStr := fmt.Sprint(time.Now().Unix())
   // Construct the payload to be signed for the HMAC.
   var payload bytes.Buffer
   payload.WriteString(method)
   payload.WriteByte('\n')
   payload.WriteString(path)
   payload.WriteByte('\n')
   payload.WriteByte('\n')
   payload.WriteString(sky_client)
   payload.WriteByte('\n')
   payload.WriteString(sky_version)
   payload.WriteByte('\n')
   payload.WriteString(headersMD5)
   payload.WriteByte('\n')
   payload.WriteString(timestampStr)
   payload.WriteByte('\n')
   payload.WriteString(bodyMD5)
   payload.WriteByte('\n')
   // Calculate the HMAC signature.
   mac := hmac.New(sha1.New, []byte(sky_key))
   payload.WriteTo(mac)
   signature := base64.StdEncoding.EncodeToString(mac.Sum(nil))
   // Format the final output string.
   return fmt.Sprintf(
      "SkyOTT client=%q,signature=%q,timestamp=%q,version=%q",
      sky_client,
      signature,
      timestampStr,
      sky_version,
   )
}

type Endpoint struct {
   Cdn string
   Url string
}

type Playout struct {
   Asset struct {
      Endpoints []Endpoint
   }
   Description string
   Protection  struct {
      LicenceAcquisitionUrl *Url
   }
}

func (*Playout) CachePath() string {
   return "rosso/peacock/Playout"
}

// L3 max 1080p
func (p *Playout) FetchWidevine(body []byte) ([]byte, error) {
   target := p.Protection.LicenceAcquisitionUrl.Url
   req, err := http.NewRequest("POST", target.String(), bytes.NewReader(body))
   if err != nil {
      return nil, err
   }
   req.Header.Set("x-sky-signature", generate_sky_ott("POST", target.Path, nil, body))
   resp, err := doRequest(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   if resp.StatusCode != 200 {
      return nil, errors.New(resp.Status)
   }
   return io.ReadAll(resp.Body)
}

func (p *Playout) GetFastly() (*url.URL, error) {
   for _, endpoint_data := range p.Asset.Endpoints {
      if endpoint_data.Cdn == "FASTLY" {
         return url.Parse(endpoint_data.Url)
      }
   }
   return nil, errors.New("FASTLY endpoint not found")
}

// userToken is good for one day
type Token struct {
   Description string
   UserToken   string
}

func FetchToken(idSession *Cookie) (*Token, error) {
   body, err := json.Marshal(map[string]any{
      "auth": map[string]string{
         "authScheme":        "MESSO",
         "proposition":       "NBCUOTT",
         "provider":          "NBCU",
         "providerTerritory": Territory,
      },
      "device": map[string]string{
         // if empty /drm/widevine/acquirelicense will fail with
         // {
         //    "errorCode": "OVP_00306",
         //    "description": "Security failure"
         // }
         "drmDeviceId": "UNKNOWN",
         // if incorrect /video/playouts/vod will fail with
         // {
         //    "errorCode": "OVP_00311",
         //    "description": "Unknown deviceId"
         // }
         // changing this too often will result in a four hour block
         // {
         //    "errorCode": "OVP_00014",
         //    "description": "Maximum number of streaming devices exceeded"
         // }
         "id":       "PC",
         "platform": "ANDROIDTV",
         "type":     "TV",
      },
   })
   if err != nil {
      return nil, err
   }
   target := url.URL{
      Scheme: "https",
      Host:   "ovp.peacocktv.com",
      Path:   "/auth/tokens",
   }
   req, err := http.NewRequest("POST", target.String(), bytes.NewReader(body))
   if err != nil {
      return nil, err
   }
   req.Header.Set("content-type", "application/vnd.tokens.v1+json")
   req.Header.Set("cookie", idSession.String())
   req.Header.Set("x-sky-signature", generate_sky_ott("POST", target.Path, nil, body))
   resp, err := doRequest(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   var result Token
   err = json.NewDecoder(resp.Body).Decode(&result)
   if err != nil {
      return nil, err
   }
   if result.Description != "" {
      return nil, errors.New(result.Description)
   }
   return &result, nil
}

func (t *Token) FetchPlayout(variantId string) (*Playout, error) {
   body, err := json.Marshal(map[string]any{
      "device": map[string]any{
         "capabilities": []any{
            map[string]string{
               "acodec":     "AAC",
               "container":  "ISOBMFF",
               "protection": "WIDEVINE",
               "transport":  "DASH",
               "vcodec":     "H264",
            },
         },
         "maxVideoFormat": "HD",
      },
      "personaParentalControlRating": 9,
      // "contentId": "GMO_00000000261361_02_HDSDR",
      "providerVariantId": variantId,
   })
   if err != nil {
      return nil, err
   }
   target := url.URL{
      Scheme: "https",
      Host:   "ovp.peacocktv.com",
      Path:   "/video/playouts/vod",
   }
   req, err := http.NewRequest("POST", target.String(), bytes.NewReader(body))
   if err != nil {
      return nil, err
   }
   // `application/json` fails
   req.Header.Set("content-type", "application/vnd.playvod.v1+json")
   req.Header.Set("x-skyott-usertoken", t.UserToken)
   header := map[string]string{
      "content-type":       "application/vnd.playvod.v1+json",
      "x-skyott-usertoken": t.UserToken,
   }
   req.Header.Set("x-sky-signature", generate_sky_ott(
      "POST", target.Path, header, body,
   ))
   resp, err := doRequest(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   var result Playout
   err = json.NewDecoder(resp.Body).Decode(&result)
   if err != nil {
      return nil, err
   }
   if result.Description != "" {
      return nil, errors.New(result.Description)
   }
   return &result, nil
}

type Url struct {
   Url url.URL
}

func (u *Url) MarshalText() ([]byte, error) {
   return u.Url.MarshalBinary()
}

func (u *Url) UnmarshalText(text []byte) error {
   return u.Url.UnmarshalBinary(text)
}

func (*Cookie) CachePath() string {
   return "rosso/peacock/Cookie"
}

func (c *Cookie) String() string {
   return fmt.Sprintf("%v=%v", c.Name, c.Value)
}
