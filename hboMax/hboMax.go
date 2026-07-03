package hboMax

import (
   "41.neocities.org/maya"
   "encoding/json"
   "errors"
   "fmt"
   "io"
   "net/url"
   "strings"
)

const (
   disco_client = "!:!:beam:!"
   device_info  = "!/!(!/!;!/!;!/!)"
)

const Markets = "amer apac emea latam"

type Cookie struct {
   Name  string
   Value string
}

func StRequest() (*Cookie, error) {
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
         "x-disco-params": "realm=bolt",
      },
   )
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   for _, each := range resp.Cookies() {
      if each.Name == "st" {
         return &Cookie{Name: each.Name, Value: each.Value}, nil
      }
   }
   return nil, errors.New("named cookie not present")
}

func (*Cookie) CachePath() string {
   return "rosso/hboMax/Cookie"
}

func (c *Cookie) String() string {
   return fmt.Sprintf("%v=%v", c.Name, c.Value)
}

type Initiate struct {
   LinkingCode string
   TargetUrl   string
}

func InitiateRequest(st *Cookie, market string) (*Initiate, error) {
   resp, err := maya.Post(
      &url.URL{
         Scheme: "https",
         Host:   fmt.Sprintf("default.beam-%v.prd.api.discomax.com", market),
         Path:   "/authentication/linkDevice/initiate",
      },
      map[string]string{
         "cookie":         st.String(),
         "x-device-info":  device_info,
         "x-disco-client": disco_client,
         "x-disco-params": "realm=bolt",
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

func (i *Initiate) String() string {
   var data strings.Builder
   data.WriteString("target URL: ")
   data.WriteString(i.TargetUrl)
   data.WriteString("\nlinking code: ")
   data.WriteString(i.LinkingCode)
   return data.String()
}

type Login struct {
   Token string
}

// you must
// /authentication/linkDevice/initiate
// first or this will always fail
func LoginRequest(st *Cookie) (*Login, error) {
   resp, err := maya.Post(
      &url.URL{
         Scheme: "https",
         Host:   "default.prd.api.hbomax.com",
         Path:   "/authentication/linkDevice/login",
      },
      map[string]string{
         "cookie":         st.String(),
         "x-device-info":  device_info,
         "x-disco-client": disco_client,
         "x-disco-params": "realm=bolt",
      },
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

func (*Login) CachePath() string {
   return "rosso/hboMax/Login"
}

type Playback struct {
   Drm struct {
      Schemes struct {
         PlayReady *Scheme
         Widevine  *Scheme
      }
   }
   Errors []struct { // 2026-05-27
      Detail string // 2026-05-27
   }
   Fallback struct {
      Manifest struct {
         Url *Url // _fallback.mpd:1080p, .mpd:4K
      }
   }
   Manifest struct {
      Url string // 1080p
   }
}

func PlayReadyRequest(token, editId string) (*Playback, error) {
   return playback_request(token, editId, "playready")
}

func WidevineRequest(token, editId string) (*Playback, error) {
   return playback_request(token, editId, "widevine")
}

func playback_request(token, edit_id, drm string) (*Playback, error) {
   body, err := json.Marshal(map[string]any{
      "editId":               edit_id,
      "consumptionType":      "streaming",
      "appBundle":            "",         // required
      "applicationSessionId": "",         // required
      "firstPlay":            false,      // required
      "gdpr":                 false,      // required
      "playbackSessionId":    "",         // required
      "userPreferences":      struct{}{}, // required
      "capabilities": map[string]any{
         "contentProtection": map[string]any{
            "contentDecryptionModules": []any{
               map[string]string{
                  "drmKeySystem": drm,
               },
            },
         },
         "manifests": map[string]any{
            "formats": map[string]any{
               "dash": struct{}{}, // required
            }, // required
         }, // required
      }, // required
      "deviceInfo": map[string]any{
         "player": map[string]any{
            "mediaEngine": map[string]string{
               "name":    "", // required
               "version": "", // required
            }, // required
            "playerView": map[string]int{
               "height": 0, // required
               "width":  0, // required
            }, // required
            "sdk": map[string]string{
               "name":    "", // required
               "version": "", // required
            }, // required
         }, // required
      }, // required
   })
   if err != nil {
      return nil, err
   }
   resp, err := maya.Post(
      &url.URL{
         Scheme: "https",
         Host:   "default.prd.api.hbomax.com",
         Path:   "/playback-orchestrator/any/playback-orchestrator/v1/playbackInfo",
      },
      map[string]string{
         "authorization":  "Bearer " + token,
         "content-type":   "application/json",
         "x-device-info":  device_info,
         "x-disco-client": disco_client,
         "x-disco-params": "realm=bolt",
      },
      body,
   )
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   var result Playback
   err = json.NewDecoder(resp.Body).Decode(&result)
   if err != nil {
      return nil, err
   }
   if len(result.Errors) >= 1 {
      return nil, errors.New(result.Errors[0].Detail)
   }
   return &result, nil
}

func (*Playback) CachePath() string {
   return "rosso/hboMax/Playback"
}

func (p *Playback) GetManifest() *url.URL {
   manifest := p.Fallback.Manifest.Url.Url
   manifest.Path = strings.Replace(manifest.Path, "_fallback", "", 1)
   return &manifest
}

// SL2000 max 1080p
// SL3000 max 2160p
func (p *Playback) PlayReadyRequest(body []byte) ([]byte, error) {
   resp, err := maya.Post(
      &p.Drm.Schemes.PlayReady.LicenseUrl.Url,
      map[string]string{"content-type": "text/xml"}, body,
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

func (p *Playback) WidevineRequest(body []byte) ([]byte, error) {
   resp, err := maya.Post(
      &p.Drm.Schemes.Widevine.LicenseUrl.Url,
      map[string]string{"content-type": "application/x-protobuf"}, body,
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

type Scheme struct {
   LicenseUrl *Url
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
