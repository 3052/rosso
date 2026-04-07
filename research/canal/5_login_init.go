package main

import (
   "bytes"
   "encoding/json"
   "fmt"
   "net/http"
   "time"
)

type LoginInitPayload struct {
   ProvisionData string     `json:"provisionData"`
   DeviceInfo    DeviceInfo `json:"deviceInfo"`
   OldSsoToken   string     `json:"oldSsoToken"`
}

func (a *App) LoginInit() error {
   url := "https://m7cp.login.solocoo.tv/login"

   payload := LoginInitPayload{
      ProvisionData: a.ProvisionData,
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

   body, _ := json.Marshal(payload)
   req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
   if err != nil {
      return err
   }

   setCommonHeaders(req)
   req.Header.Set("Content-Type", "application/json")

   timestamp := time.Now().Unix()
   req.Header.Set("Authorization", GenerateAuthorizationHeader(url, body, timestamp))

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
