package hulu

import (
   "41.neocities.org/maya"
   "encoding/json"
   "errors"
   "io"
   "net/url"
   "path"
)

type DeepLink struct {
   EabId   string `json:"eab_id"`
   Message string
}

func (*Device) CachePath() string {
   return "rosso/hulu/Device"
}

type Device struct {
   DeviceToken string `json:"device_token"`
   Message     string // 2026-05-02
   UserToken   string `json:"user_token"`
}

func (*Playlist) CachePath() string {
   return "rosso/hulu/Playlist"
}

type Playlist struct {
   DashPrServer *Url `json:"dash_pr_server"`
   StreamUrl    *Url `json:"stream_url"` // MPD
   WvServer     *Url `json:"wv_server"`
}

func (p *Playlist) FetchWidevine(body []byte) ([]byte, error) {
   resp, err := maya.Post(
      &p.WvServer.Url,
      map[string]string{"content-type": "application/x-protobuf"}, body,
   )
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   return io.ReadAll(resp.Body)
}

func (p *Playlist) FetchPlayReady(body []byte) ([]byte, error) {
   resp, err := maya.Post(&p.DashPrServer.Url, nil, body)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   body, err = io.ReadAll(resp.Body)
   if err != nil {
      return nil, err
   }
   if resp.StatusCode != 200 {
      var result struct {
         Message string
      }
      err = json.Unmarshal(body, &result)
      if err != nil {
         return nil, err
      }
      return nil, errors.New(result.Message)
   }
   return body, nil
}

type Url struct {
   Url url.URL
}

func (u *Url) UnmarshalText(text []byte) error {
   return u.Url.UnmarshalBinary(text)
}

func (u *Url) MarshalText() ([]byte, error) {
   return u.Url.MarshalBinary()
}

///

func (d *Device) DeepLink(id string) (*DeepLink, error) {
   resp, err := maya.Get(
      &url.URL{
         Scheme: "https",
         Host:   "discover.hulu.com",
         Path:   "/content/v5/deeplink/playback",
         RawQuery: url.Values{
            "id":        {id},
            "namespace": {"entity"},
         }.Encode(),
      },
      map[string]string{"authorization": "Bearer " + d.UserToken},
   )
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()

   var result DeepLink
   err = json.NewDecoder(resp.Body).Decode(&result)
   if err != nil {
      return nil, err
   }
   if result.EabId == "" {
      return nil, errors.New("content is not playable: missing eab_id in response")
   }
   return &result, nil
}

// returns user_token only
func (d *Device) TokenRefresh() error {
   body := url.Values{
      "action":       {"token_refresh"},
      "device_token": {d.DeviceToken},
   }.Encode()
   resp, err := maya.Post(
      &url.URL{
         Scheme: "https",
         Host:   "auth.hulu.com",
         Path:   "/v1/device/device_token/authenticate",
      },
      map[string]string{"content-type": "application/x-www-form-urlencoded"},
      []byte(body),
   )
   if err != nil {
      return err
   }
   defer resp.Body.Close()
   var result Device
   err = json.NewDecoder(resp.Body).Decode(&result)
   if err != nil {
      return err
   }
   if result.Message != "" {
      return errors.New(result.Message)
   }
   return nil
}

type Details struct {
   VodItems struct {
      Focus struct {
         Entity struct {
            Bundle struct {
               EabId string `json:"eab_id"`
            }
         }
      }
   } `json:"vod_items"`
}

func (d *Device) GetDetails(movie string) (*Details, error) {
   resp, err := maya.Get(
      &url.URL{
         Scheme:   "https",
         Host:     "discover.hulu.com",
         Path:     "/content/v5/hubs/movie/" + movie,
         RawQuery: "limit=0",
      },
      map[string]string{"authorization": "Bearer " + d.UserToken},
   )
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   var result struct {
      Details Details
   }
   err = json.NewDecoder(resp.Body).Decode(&result)
   if err != nil {
      return nil, err
   }
   return &result.Details, nil
}

// L3 max 1080p
// SL2000 max 1080p
// SL3000 max 2160p
func (d *Device) Playlist(eabId string) (*Playlist, error) {
   body, err := json.Marshal(map[string]any{
      "deejay_device_id": deejay[0].device_id,
      "content_eab_id":   eabId,
      "unencrypted":      true,
      "version":          deejay[0].key_version,
      "playback": map[string]any{
         "audio": map[string]any{
            "codecs": map[string]any{
               "selection_mode": "ALL",
               "values": []any{
                  map[string]string{"type": "AAC"},
                  map[string]string{"type": "EC3"},
               },
            },
         },
         "drm": map[string]any{
            "multi_key":      true, // NEED THIS FOR 4K UHD
            "selection_mode": "ALL",
            "values": []any{
               map[string]string{
                  "security_level": "L3",
                  "type":           "WIDEVINE",
                  "version":        "MODULAR",
               },
               map[string]string{
                  "security_level": "SL2000",
                  "type":           "PLAYREADY",
                  "version":        "V2",
               },
            },
         },
         "version": 2, // needs to be exactly 2 for 1080p
         "manifest": map[string]string{
            "type": "DASH",
         },
         "segments": map[string]any{
            "selection_mode": "ALL",
            "values": []any{
               map[string]any{
                  "type": "FMP4",
                  "encryption": map[string]string{
                     "mode": "CENC",
                     "type": "CENC",
                  },
               },
            },
         },
         "video": map[string]any{
            "codecs": map[string]any{
               "selection_mode": "ALL",
               "values": []any{
                  map[string]any{
                     "height":  9999,
                     "level":   "9",
                     "profile": "HIGH",
                     "type":    "H264",
                     "width":   9999,
                  },
                  map[string]any{
                     "height":  9999,
                     "level":   "9",
                     "profile": "MAIN_10",
                     "tier":    "MAIN",
                     "type":    "H265",
                     "width":   9999,
                  },
               },
            },
         },
      },
   })
   if err != nil {
      return nil, err
   }
   resp, err := maya.Post(
      &url.URL{
         Scheme: "https",
         Host:   "play.hulu.com",
         Path:   "/v6/playlist",
      },
      map[string]string{
         "authorization": "Bearer " + d.UserToken,
         "content-type":  "application/json",
      },
      body,
   )
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   var result Playlist
   err = json.NewDecoder(resp.Body).Decode(&result)
   if err != nil {
      return nil, err
   }
   return &result, nil
}

func FetchDevice(email, password string) (*Device, error) {
   body := url.Values{
      "friendly_name": {"!"},
      "password":      {password},
      "serial_number": {"!"},
      "user_email":    {email},
   }.Encode()
   resp, err := maya.Post(
      &url.URL{
         Scheme: "https",
         Host:   "auth.hulu.com",
         Path:   "/v2/livingroom/password/authenticate",
      },
      map[string]string{"content-type": "application/x-www-form-urlencoded"},
      []byte(body),
   )
   if err != nil {
      return nil, err
   }
   if resp.StatusCode != 200 {
      return nil, errors.New(resp.Status)
   }
   defer resp.Body.Close()
   var result struct {
      Data Device
   }
   err = json.NewDecoder(resp.Body).Decode(&result)
   if err != nil {
      return nil, err
   }
   return &result.Data, nil
}

// https://hulu.com/movie/05e76ad8-c3dd-4c3e-bab9-df3cf71c6871
// https://hulu.com/movie/alien-romulus-05e76ad8-c3dd-4c3e-bab9-df3cf71c6871
func ParseId(urlData string) string {
   part := path.Base(urlData)
   len_part := len(part)
   const len_uuid = 36
   if len_part > len_uuid {
      if part[len_part-len_uuid-1] == '-' {
         return part[len_part-len_uuid:]
      }
   }
   return part
}

var deejay = []struct {
   resolution  string
   device_id   int
   key_version int
}{
   {
      resolution:  "2160p",
      device_id:   210,
      key_version: 1,
   },
   {
      resolution:  "2160p",
      device_id:   208,
      key_version: 1,
   },
   {
      resolution:  "2160p",
      device_id:   204,
      key_version: 4,
   },
   {
      resolution:  "2160p",
      device_id:   188,
      key_version: 17,
   },
   {
      resolution:  "720p",
      device_id:   214,
      key_version: 1,
   },
   {
      resolution:  "720p",
      device_id:   191,
      key_version: 1,
   },
   {
      resolution:  "720p",
      device_id:   190,
      key_version: 1,
   },
   {
      resolution:  "720p",
      device_id:   142,
      key_version: 1,
   },
   {
      resolution:  "720p",
      device_id:   109,
      key_version: 1,
   },
}
