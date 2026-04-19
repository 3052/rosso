package hulu

import (
   "41.neocities.org/maya"
   "encoding/json"
   "errors"
   "io"
   "net/url"
)

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
