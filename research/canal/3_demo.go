package main

import (
   "bytes"
   "encoding/json"
   "fmt"
   "net/http"
   "time"
)

type DeviceInfo struct {
   OsVersion        string `json:"osVersion"`
   DeviceModel      string `json:"deviceModel"`
   DeviceType       string `json:"deviceType"`
   DeviceSerial     string `json:"deviceSerial"`
   DeviceOem        string `json:"deviceOem"`
   DevicePrettyName string `json:"devicePrettyName"`
   AppVersion       string `json:"appVersion"`
   Language         string `json:"language"`
   Brand            string `json:"brand"`
   Country          string `json:"country,omitempty"`
}

type DemoPayload struct {
   ProvisionData string     `json:"provisionData"`
   DeviceInfo    DeviceInfo `json:"deviceInfo"`
}

func (a *App) Demo() error {
   url := "https://m7cp.login.solocoo.tv/demo"

   payload := DemoPayload{
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
   }

   body, _ := json.Marshal(payload)
   req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
   if err != nil {
      return err
   }

   setCommonHeaders(req)
   req.Header.Set("Content-Type", "application/json")

   // Create dynamic HMAC Authorization Header
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
      SsoToken string `json:"ssoToken"`
   }
   if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
      return err
   }

   a.SsoToken = result.SsoToken
   return nil
}
