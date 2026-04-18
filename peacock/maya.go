package peacock

import (
   "41.neocities.org/maya"
   "bytes"
   "crypto/hmac"
   "crypto/md5"
   "crypto/sha1"
   "encoding/base64"
   "encoding/json"
   "errors"
   "fmt"
   "io"
   "maps"
   "net/url"
   "slices"
   "strings"
   "time"
)

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
   header := map[string]string{
      // `application/json` fails
      "content-type":       "application/vnd.playvod.v1+json",
      "x-skyott-usertoken": t.UserToken,
   }
   header["x-sky-signature"] = generate_sky_ott(
      "POST", target.Path, header, body,
   )
   resp, err := maya.Post(&target, header, body)
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

// L3 max 1080p
func (p *Playout) FetchWidevine(body []byte) ([]byte, error) {
   target, err := url.Parse(p.Protection.LicenceAcquisitionUrl)
   if err != nil {
      return nil, err
   }
   resp, err := maya.Post(
      target,
      map[string]string{
         "x-sky-signature": generate_sky_ott("POST", target.Path, nil, body),
      },
      body,
   )
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   if resp.StatusCode != 200 {
      return nil, errors.New(resp.Status)
   }
   return io.ReadAll(resp.Body)
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

func FetchToken(idSession string) (*Token, error) {
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
   resp, err := maya.Post(
      &target,
      map[string]string{
         "content-type":    "application/vnd.tokens.v1+json",
         "cookie":          idSession,
         "x-sky-signature": generate_sky_ott("POST", target.Path, nil, body),
      },
      body,
   )
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

func FetchIdSession(user, password string) (string, error) {
   body := url.Values{
      "userIdentifier": {user},
      "password":       {password},
   }.Encode()
   resp, err := maya.Post(
      &url.URL{
         Scheme: "https",
         Host:   "rango.id.peacocktv.com",
         Path:   "/signin/service/international",
      },
      map[string]string{
         "content-type":         "application/x-www-form-urlencoded",
         "x-skyott-proposition": "NBCUOTT",
         "x-skyott-provider":    "NBCU",
         "x-skyott-territory":   Territory,
      },
      []byte(body),
   )
   if err != nil {
      return "", err
   }
   defer resp.Body.Close()
   var result struct {
      Properties struct {
         Errors struct {
            CategoryErrors []struct {
               Code string
            }
         }
      }
   }
   err = json.NewDecoder(resp.Body).Decode(&result)
   if err != nil {
      return "", err
   }
   if resp.StatusCode != 201 {
      return "", errors.New(result.Properties.Errors.CategoryErrors[0].Code)
   }
   for _, cookie := range resp.Cookies() {
      if cookie.Name == "idsession" {
         return cookie.String(), nil
      }
   }
   return "", errors.New("named cookie not present")
}
