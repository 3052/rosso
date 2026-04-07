package main

import (
   "bytes"
   "encoding/json"
   "fmt"
   "net/http"
)

func (a *App) LoginInit() error {
   url := "https://m7cp.login.solocoo.tv/login"

   payload := map[string]interface{}{
      "provisionData": a.ProvisionData,
      "deviceInfo": map[string]interface{}{
         "osVersion":        "Windows 10",
         "deviceModel":      "Firefox",
         "deviceType":       "PC",
         "deviceSerial":     a.DeviceSerial,
         "deviceOem":        "Firefox",
         "devicePrettyName": "Firefox 140.0",
         "appVersion":       "12.9",
         "language":         "en_US",
         "brand":            "m7cp",
         "country":          "CZ",
      },
      "oldSsoToken": a.SsoToken,
   }

   body, _ := json.Marshal(payload)
   req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
   if err != nil {
      return err
   }
   setCommonHeaders(req)
   req.Header.Set("Content-Type", "application/json")
   // Using hardcoded Client Auth signature from HAR
   req.Header.Set("Authorization", "Client key=web.NhFyz4KsZ54,time=1775520334,sig=v62D8nVsb1jR5PkV6b_ou0MqC1TJFmapzeT1Lb6Dniw")

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
