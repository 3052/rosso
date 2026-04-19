package hulu

import (
   "net/url"
   "path"
)

func (p *Playlist) GetManifest() (*url.URL, error) {
   return url.Parse(p.StreamUrl)
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

type Device struct {
   DeviceToken string `json:"device_token"`
   UserToken   string `json:"user_token"`
}

type DeepLink struct {
   EabId   string `json:"eab_id"`
   Message string
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
