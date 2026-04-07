package main

import (
   "bytes"
   "encoding/json"
   "fmt"
   "net/http"
   "net/url"
)

func (a *App) LoginInit() error {
   u, err := url.Parse("https://m7cp.login.solocoo.tv/login")
   if err != nil {
      return err
   }
   payload := LoginInitPayload{
      DeviceInfo: DeviceInfo{
         OsVersion:        "Windows 10",
         DeviceModel:      "Firefox",
         DeviceType:       "PC",
         DeviceSerial:     a.DeviceSerial,
         DeviceOem:        "Firefox",
         DevicePrettyName: "Firefox 140.0",
         AppVersion:       "12.9",
         Language:         "en_US",
         Brand:            "m7cp",
         Country:          "CZ",
      },
      OldSsoToken: a.SsoToken,
   }

   body, err := json.Marshal(payload)
   if err != nil {
      return err
   }

   authHeader, err := get_client(u, body)
   if err != nil {
      return err
   }

   req, err := http.NewRequest("POST", u.String(), bytes.NewBuffer(body))
   if err != nil {
      return err
   }
   setCommonHeaders(req)
   req.Header.Set("Content-Type", "application/json")
   req.Header.Set("Authorization", authHeader)

   resp, err := a.Client.Do(req)
   if err != nil {
      return err
   }
   defer resp.Body.Close()

   if resp.StatusCode != 200 {
      return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
   }

   var result struct {
      Ticket string `json:"ticket"`
   }
   if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
      return err
   }

   a.Ticket = result.Ticket
   return nil
}

type LoginInitPayload struct {
   ProvisionData string     `json:"provisionData"`
   DeviceInfo    DeviceInfo `json:"deviceInfo"`
   OldSsoToken   string     `json:"oldSsoToken"`
}
