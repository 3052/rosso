package canal

import (
   "41.neocities.org/maya"
   "encoding/json"
   "errors"
   "fmt"
   "io"
   "net/url"
)

func (t *Ticket) Login(username, password string) (*Login, error) {
   body, err := json.Marshal(map[string]any{
      "ticket": t.Ticket,
      "userInput": map[string]string{
         "username": username,
         "password": password,
      },
   })
   if err != nil {
      return nil, err
   }
   target := &url.URL{
      Scheme: "https", Host: "m7cp.login.solocoo.tv", Path: "/login",
   }
   client, err := get_client(target, body)
   if err != nil {
      return nil, err
   }
   resp, err := maya.Post(
      target,
      map[string]string{
         "authorization": client,
         "user-agent":    user_agent,
      },
      body,
   )
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   var result Login
   err = json.NewDecoder(resp.Body).Decode(&result)
   if err != nil {
      return nil, err
   }
   if resp.StatusCode != 200 {
      return nil, &result
   }
   return &result, nil
}

func FetchTicket() (*Ticket, error) {
   body, err := json.Marshal(map[string]any{
      "deviceInfo": map[string]string{
         "brand":        "m7cp", // sg.ui.sso.fatal.internal_error
         "deviceModel":  "Firefox",
         "deviceOem":    "Firefox",
         "deviceSerial": device_serial,
         "deviceType":   "PC",
         "osVersion":    "Windows 10",
      },
   })
   if err != nil {
      return nil, err
   }
   target := &url.URL{
      Scheme: "https", Host: "m7cp.login.solocoo.tv", Path: "/login",
   }
   client, err := get_client(target, body)
   if err != nil {
      return nil, err
   }
   resp, err := maya.Post(
      target,
      map[string]string{
         "authorization": client,
         "user-agent":    user_agent,
      },
      body,
   )
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   var result Ticket
   err = json.NewDecoder(resp.Body).Decode(&result)
   if err != nil {
      return nil, err
   }
   if result.Message != "" {
      return nil, errors.New(result.Message)
   }
   return &result, nil
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
   resp, err := maya.Post(
      &url.URL{
         Scheme: "https",
         Host:   "tvapi-hlm2.solocoo.tv",
         Path:   fmt.Sprintf("/v1/assets/%v/play", tracking),
      },
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
   if result.Message != "" {
      return nil, errors.New(result.Message)
   }
   return &result, nil
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
   resp, err := maya.Post(
      &url.URL{
         Scheme: "https", Host: "tvapi-hlm2.solocoo.tv", Path: "/v1/session",
      },
      nil,
      body,
   )
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   var result Session
   err = json.NewDecoder(resp.Body).Decode(&result)
   if err != nil {
      return nil, err
   }
   if result.Message != "" {
      return nil, errors.New(result.Message)
   }
   return &result, nil
}

func (s *Session) Episodes(tracking string, season int) ([]Episode, error) {
   resp, err := maya.Get(
      &url.URL{
         Scheme: "https",
         Host:   "tvapi-hlm2.solocoo.tv",
         Path:   "/v1/assets",
         RawQuery: url.Values{
            "limit": {"99"},
            "query": {fmt.Sprintf("episodes,%v,season,%v", tracking, season)},
         }.Encode(),
      },
      map[string]string{"authorization": "Bearer " + s.Token},
   )
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   var result struct {
      Assets  []Episode
      Message string
   }
   err = json.NewDecoder(resp.Body).Decode(&result)
   if err != nil {
      return nil, err
   }
   if result.Message != "" {
      return nil, errors.New(result.Message)
   }
   return result.Assets, nil
}

func (s *Session) Search(query string) ([]Collection, error) {
   resp, err := maya.Get(
      &url.URL{
         Scheme:   "https",
         Host:     "tvapi-hlm2.solocoo.tv",
         Path:     "/v1/search",
         RawQuery: url.Values{"query": {query}}.Encode(),
      },
      map[string]string{"authorization": "Bearer " + s.Token},
   )
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   var result struct {
      Collection []Collection
      Message    string
   }
   err = json.NewDecoder(resp.Body).Decode(&result)
   if err != nil {
      return nil, err
   }
   if result.Message != "" {
      return nil, errors.New(result.Message)
   }
   return result.Collection, nil
}

func (p *Player) FetchWidevine(body []byte) ([]byte, error) {
   target, err := url.Parse(p.Drm.LicenseUrl)
   if err != nil {
      return nil, err
   }
   resp, err := maya.Post(target, nil, body)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   return io.ReadAll(resp.Body)
}
