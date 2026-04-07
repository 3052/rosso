package main

import (
   "bytes"
   "encoding/json"
   "fmt"
   "net/http"
)

type ProvisionPayload struct {
   OsVersion        string `json:"osVersion"`
   DeviceModel      string `json:"deviceModel"`
   DeviceType       string `json:"deviceType"`
   DeviceSerial     string `json:"deviceSerial"`
   DeviceOem        string `json:"deviceOem"`
   DevicePrettyName string `json:"devicePrettyName"`
   AppVersion       string `json:"appVersion"`
   Language         string `json:"language"`
   Brand            string `json:"brand"`
   FeatureLevel     int    `json:"featureLevel"`
}

func (a *App) Provision() error {
   url := fmt.Sprintf("%s/v1/provision", a.TVApiBaseURL)

   payload := ProvisionPayload{
      OsVersion:        "Windows 10",
      DeviceModel:      "Firefox",
      DeviceType:       "PC",
      DeviceSerial:     a.DeviceSerial,
      DeviceOem:        "Firefox",
      DevicePrettyName: "Firefox 140.0",
      AppVersion:       "12.9",
      Language:         "en_US",
      Brand:            "m7cp",
      FeatureLevel:     7,
   }

   body, _ := json.Marshal(payload)
   req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
   if err != nil {
      return err
   }
   setCommonHeaders(req)
   req.Header.Set("Content-Type", "application/json")

   resp, err := a.Client.Do(req)
   if err != nil {
      return err
   }
   defer resp.Body.Close()

   var result struct {
      Session struct {
         ProvisionData string `json:"provisionData"`
      } `json:"session"`
   }
   if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
      return err
   }

   a.ProvisionData = result.Session.ProvisionData
   return nil
}
