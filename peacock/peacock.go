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
   "net/http"
   "net/url"
   "slices"
   "strings"
   "time"
)

func (e *Endpoint) GetManifest() (*url.URL, error) {
   return url.Parse(e.Url)
}

func (p *Playout) GetFastly() (*Endpoint, error) {
   for _, endpoint_data := range p.Asset.Endpoints {
      if endpoint_data.Cdn == "FASTLY" {
         return &endpoint_data, nil
      }
   }
   return nil, errors.New("FASTLY endpoint not found")
}

type Playout struct {
   Asset struct {
      Endpoints []Endpoint
   }
   Description string
   Protection  struct {
      LicenceAcquisitionUrl string
   }
}

type Endpoint struct {
   Cdn string
   Url string
}

const (
   sky_client  = "NBCU-ANDROID-v3"
   sky_key     = "JuLQgyFz9n89D9pxcN6ZWZXKWfgj2PNBUb32zybj"
   sky_version = "1.0"
)

func generate_sky_ott(method, path string, headers http.Header, body []byte) string {
   // Sort headers by key.
   headerKeys := make([]string, 0, len(headers))
   for key := range headers {
      headerKeys = append(headerKeys, key)
   }
   slices.Sort(headerKeys)
   // Build the special headers string.
   var headersBuilder bytes.Buffer
   for _, key := range headerKeys {
      lowerKey := strings.ToLower(key)
      if strings.HasPrefix(lowerKey, "x-skyott-") {
         value := headers.Get(key)
         headersBuilder.WriteString(lowerKey)
         headersBuilder.WriteString(": ")
         headersBuilder.WriteString(value)
         headersBuilder.WriteByte('\n')
      }
   }
   // MD5 the headers string and request body.
   headersHash := md5.Sum(headersBuilder.Bytes())
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

// userToken is good for one day
type Token struct {
   Description string
   UserToken   string
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
   req, err := http.NewRequest(
      "POST", "https://ovp.peacocktv.com/video/playouts/vod",
      bytes.NewReader(body),
   )
   if err != nil {
      return nil, err
   }
   // `application/json` fails
   req.Header.Set("content-type", "application/vnd.playvod.v1+json")
   req.Header.Set("x-skyott-usertoken", t.UserToken)
   req.Header.Set(
      "x-sky-signature",
      generate_sky_ott(req.Method, req.URL.Path, req.Header, body),
   )
   resp, err := http.DefaultClient.Do(req)
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

var Territory = "US"
