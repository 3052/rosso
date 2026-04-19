package hulu

import (
   "41.neocities.org/maya"
   "encoding/json"
   "errors"
   "io"
   "net/url"
)

// L3 max 1080p
// SL2000 max 1080p
// SL3000 max 2160p
func (d *Device) Playlist(eabId string) (*Playlist, error) {
   body, err := json.Marshal(map[string]any{
      "deejay_device_id": deejay[0].device_id,
      "version":          deejay[0].key_version,
      "content_eab_id":   eabId,
      "unencrypted":      true,
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
   if result.Message != "" {
      return nil, errors.New(result.Message)
   }
   return &result, nil
}

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
   if result.Message != "" {
      return nil, errors.New(result.Message)
   }
   return &result, nil
}

// returns user_token only
func (d *Device) TokenRefresh() (*Device, error) {
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
      return nil, err
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

func (p *Playlist) FetchWidevine(body []byte) ([]byte, error) {
   target, err := url.Parse(p.WvServer)
   if err != nil {
      return nil, err
   }
   resp, err := maya.Post(
      target, map[string]string{"content-type": "application/x-protobuf"}, body,
   )
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   return io.ReadAll(resp.Body)
}

func (p *Playlist) FetchPlayReady(body []byte) ([]byte, error) {
   target, err := url.Parse(p.DashPrServer)
   if err != nil {
      return nil, err
   }
   resp, err := maya.Post(target, nil, body)
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
