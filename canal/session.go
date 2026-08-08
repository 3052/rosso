package canal

import (
   "bytes"
   "encoding/json"
   "errors"
   "fmt"
   "log"
   "net/http"
   "net/url"
)

func Do(req *http.Request) (*http.Response, error) {
   log.Println(req.Method, req.URL)
   return http.DefaultClient.Do(req)
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

type Session struct {
   SsoToken string
   Token    string // this last one hour
}

func FetchSession(ssoToken string) (*Session, error) {
   body, err := json.Marshal(map[string]string{
      "brand":        "m7cp",
      "deviceSerial": device_serial,
      "deviceType":   "PC",
      "ssoToken":     ssoToken,
   })
   if err != nil {
      return nil, err
   }
   resp, err := doPost(
      "https://tvapi-hlm2.solocoo.tv/v1/session",
      nil,
      body,
   )
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   var result Session
   if err = json.NewDecoder(resp.Body).Decode(&result); err != nil {
      return nil, err
   }
   return &result, nil
}

func (*Session) CachePath() string {
   return "rosso/canal/Session"
}

func (s *Session) Episodes(tracking string, season int) ([]Episode, error) {
   location := &url.URL{
      Scheme: "https",
      Host:   "tvapi-hlm2.solocoo.tv",
      Path:   "/v1/assets",
      RawQuery: url.Values{
         "limit": {"99"},
         "query": {fmt.Sprintf("episodes,%v,season,%v", tracking, season)},
      }.Encode(),
   }
   resp, err := doGet(location.String(), map[string]string{"authorization": "Bearer " + s.Token})
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   var result struct {
      Assets []Episode
   }
   err = json.NewDecoder(resp.Body).Decode(&result)
   if err != nil {
      return nil, err
   }
   return result.Assets, nil
}

func (s *Session) Player(tracking string) (*Player, error) {
   body, err := json.Marshal(map[string]any{
      "player": map[string]any{
         "capabilities": map[string]any{
            "drmSystems": []string{"Widevine"},
            "mediaTypes": []string{"DASH"},
         },
      },
   })
   if err != nil {
      return nil, err
   }
   resp, err := doPost(
      fmt.Sprintf("https://tvapi-hlm2.solocoo.tv/v1/assets/%v/play", tracking),
      map[string]string{
         "authorization": "Bearer " + s.Token,
         "content-type":  "application/json",
      },
      body,
   )
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   var result Player
   err = json.NewDecoder(resp.Body).Decode(&result)
   if err != nil {
      return nil, err
   }
   return &result, nil
}

func (s *Session) Search(query string) ([]Asset, error) {
   location := &url.URL{
      Scheme:   "https",
      Host:     "tvapi-hlm2.solocoo.tv",
      Path:     "/v1/search",
      RawQuery: url.Values{"query": {query}}.Encode(),
   }
   resp, err := doGet(location.String(), map[string]string{"authorization": "Bearer " + s.Token})
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()

   var result struct {
      Collection []struct {
         Assets []Asset
         Label  string
      }
      Message string // 2026-05-30
   }
   err = json.NewDecoder(resp.Body).Decode(&result)
   if err != nil {
      return nil, err
   }
   if result.Message != "" {
      return nil, errors.New(result.Message)
   }
   for _, collection := range result.Collection {
      if collection.Label == "sg.ui.search.vod" {
         return collection.Assets, nil
      }
   }
   return nil, errors.New("no vod found")
}

// session.go
