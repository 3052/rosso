package hboMax

import (
   "41.neocities.org/maya"
   "encoding/json"
   "errors"
   "fmt"
   "io"
   "net/http"
   "net/url"
)

func InitiateRequest(st *http.Cookie, market string) (*Initiate, error) {
   req := http.Request{
      Method: "POST",
      URL: &url.URL{
         Scheme: "https",
         Host:   fmt.Sprintf("default.beam-%v.prd.api.discomax.com", market),
         Path:   "/authentication/linkDevice/initiate",
      },
      Header: http.Header{},
   }
   req.AddCookie(st)
   req.Header.Set("x-device-info", device_info)
   resp, err := http.DefaultClient.Do(&req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   if resp.StatusCode != http.StatusOK {
      return nil, errors.New(resp.Status)
   }
   var result struct {
      Data struct {
         Attributes Initiate
      }
   }
   err = json.NewDecoder(resp.Body).Decode(&result)
   if err != nil {
      return nil, err
   }
   return &result.Data.Attributes, nil
}

// you must
// /authentication/linkDevice/initiate
// first or this will always fail
func LoginRequest(st string) (*Login, error) {
   resp, err := maya.Post(
      &url.URL{
         Scheme: "https",
         Host:   "default.prd.api.hbomax.com",
         Path:   "/authentication/linkDevice/login",
      },
      map[string]string{"cookie": st},
      nil,
   )
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   var result struct {
      Data struct {
         Attributes Login
      }
   }
   err = json.NewDecoder(resp.Body).Decode(&result)
   if err != nil {
      return nil, err
   }
   return &result.Data.Attributes, nil
}

func StRequest() (string, error) {
   resp, err := maya.Get(
      &url.URL{
         Scheme:   "https",
         Host:     "default.prd.api.hbomax.com",
         Path:     "/token",
         RawQuery: "realm=bolt",
      },
      map[string]string{
         "x-device-info":  device_info,
         "x-disco-client": disco_client,
      },
   )
   if err != nil {
      return "", err
   }
   defer resp.Body.Close()
   for _, cookie := range resp.Cookies() {
      if cookie.Name == "st" {
         return cookie.String(), nil
      }
   }
   return "", errors.New("named cookie not present")
}

// SL2000 max 1080p
// SL3000 max 2160p
func (p *Playback) PlayReadyRequest(body []byte) ([]byte, error) {
   target, err := url.Parse(p.Drm.Schemes.PlayReady.LicenseUrl)
   if err != nil {
      return nil, err
   }
   resp, err := maya.Post(
      target, map[string]string{"content-type": "text/xml"}, body,
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
