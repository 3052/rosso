package peacock

import (
   "41.neocities.org/maya"
   "bytes"
   "encoding/json"
   "errors"
   "io"
   "net/http"
   "net/url"
)

// L3 max 1080p
func (p *Playout) FetchWidevine(body []byte) ([]byte, error) {
   req, err := http.NewRequest(
      "POST", p.Protection.LicenceAcquisitionUrl, bytes.NewReader(body),
   )
   if err != nil {
      return nil, err
   }
   req.Header.Set(
      "x-sky-signature",
      generate_sky_ott(req.Method, req.URL.Path, req.Header, body),
   )
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   if resp.StatusCode != http.StatusOK {
      return nil, errors.New(resp.Status)
   }
   return io.ReadAll(resp.Body)
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
   return "", errors.New("http: named cookie not present")
}
