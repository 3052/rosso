package hulu

import (
   "bytes"
   "encoding/json"
   "errors"
   "net/http"
   "net/url"
   "path"
)

var Do = func(req *http.Request) (*http.Response, error) {
   return http.DefaultClient.Do(req)
}

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

func doGet(rawUrl string, headers map[string]string) (*http.Response, error) {
   req, err := http.NewRequest("GET", rawUrl, nil)
   if err != nil {
      return nil, err
   }
   for key, val := range headers {
      req.Header.Set(key, val)
   }
   return Do(req)
}

func doPost(rawUrl string, headers map[string]string, body []byte) (*http.Response, error) {
   req, err := http.NewRequest("POST", rawUrl, bytes.NewReader(body))
   if err != nil {
      return nil, err
   }
   for key, val := range headers {
      req.Header.Set(key, val)
   }
   return Do(req)
}

type DeepLink struct {
   EabId   string `json:"eab_id"`
   Message string
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

type Device struct {
   DeviceToken string `json:"device_token"`
   Message     string // 2026-05-02
   UserToken   string `json:"user_token"`
}

func FetchDevice(email, password string) (*Device, error) {
   body := url.Values{
      "friendly_name": {"!"},
      "password":      {password},
      "serial_number": {"!"},
      "user_email":    {email},
   }.Encode()
   resp, err := doPost(
      "https://auth.hulu.com/v2/livingroom/password/authenticate",
      map[string]string{"content-type": "application/x-www-form-urlencoded"},
      []byte(body),
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

func (*Device) CachePath() string {
   return "rosso/hulu/Device"
}

func (d *Device) DeepLink(id string) (*DeepLink, error) {
   location := &url.URL{
      Scheme: "https",
      Host:   "discover.hulu.com",
      Path:   "/content/v5/deeplink/playback",
      RawQuery: url.Values{
         "id":        {id},
         "namespace": {"entity"},
      }.Encode(),
   }
   resp, err := doGet(location.String(), map[string]string{"authorization": "Bearer " + d.UserToken})
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

func (d *Device) GetDetails(movie string) (*Details, error) {
   location := &url.URL{
      Scheme:   "https",
      Host:     "discover.hulu.com",
      Path:     "/content/v5/hubs/movie/" + movie,
      RawQuery: "limit=0",
   }
   resp, err := doGet(location.String(), map[string]string{"authorization": "Bearer " + d.UserToken})
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

// returns user_token only
func (d *Device) TokenRefresh() error {
   body := url.Values{
      "action":       {"token_refresh"},
      "device_token": {d.DeviceToken},
   }.Encode()
   resp, err := doPost(
      "https://auth.hulu.com/v1/device/device_token/authenticate",
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

// device.go
