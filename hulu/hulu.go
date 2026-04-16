package hulu

import (
   "bytes"
   "encoding/json"
   "errors"
   "io"
   "net/http"
   "net/url"
   "path"
)

func (p *Playlist) GetManifest() (*url.URL, error) {
   return url.Parse(p.StreamUrl)
}

func (p *Playlist) FetchPlayReady(data []byte) ([]byte, error) {
   resp, err := http.Post(
      p.DashPrServer, "", bytes.NewReader(data),
   )
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   data, err = io.ReadAll(resp.Body)
   if err != nil {
      return nil, err
   }
   if resp.StatusCode != http.StatusOK {
      var result struct {
         Message string
      }
      err = json.Unmarshal(data, &result)
      if err != nil {
         return nil, err
      }
      return nil, errors.New(result.Message)
   }
   return data, nil
}

func (p *Playlist) FetchWidevine(data []byte) ([]byte, error) {
   resp, err := http.Post(
      p.WvServer, "application/x-protobuf", bytes.NewReader(data),
   )
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   return io.ReadAll(resp.Body)
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

func FetchDevice(email, password string) (*Device, error) {
   resp, err := http.PostForm(
      "https://auth.hulu.com/v2/livingroom/password/authenticate", url.Values{
         "friendly_name": {"!"},
         "password":      {password},
         "serial_number": {"!"},
         "user_email":    {email},
      },
   )
   if err != nil {
      return nil, err
   }
   if resp.StatusCode != http.StatusOK {
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

// returns user_token only
func (d *Device) TokenRefresh() (*Device, error) {
   resp, err := http.PostForm(
      "https://auth.hulu.com/v1/device/device_token/authenticate", url.Values{
         "action":       {"token_refresh"},
         "device_token": {d.DeviceToken},
      },
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

type Device struct {
   DeviceToken string `json:"device_token"`
   UserToken   string `json:"user_token"`
}

type DeepLink struct {
   EabId   string `json:"eab_id"`
   Message string
}

func (d *Device) DeepLink(id string) (*DeepLink, error) {
   req := http.Request{
      URL: &url.URL{
         Scheme: "https",
         Host:   "discover.hulu.com",
         Path:   "/content/v5/deeplink/playback",
         RawQuery: url.Values{
            "id":        {id},
            "namespace": {"entity"},
         }.Encode(),
      },
      Header: http.Header{},
   }
   req.Header.Set("authorization", "Bearer "+d.UserToken)
   resp, err := http.DefaultClient.Do(&req)
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
   req, err := http.NewRequest(
      "POST", "https://play.hulu.com/v6/playlist", bytes.NewReader(body),
   )
   if err != nil {
      return nil, err
   }
   req.Header.Set("authorization", "Bearer "+d.UserToken)
   req.Header.Set("content-type", "application/json")
   resp, err := http.DefaultClient.Do(req)
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

type Playlist struct {
   DashPrServer string `json:"dash_pr_server"`
   Message      string
   StreamUrl    string `json:"stream_url"` // MPD
   WvServer     string `json:"wv_server"`
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
