package main

import (
   "bytes"
   "encoding/json"
   "fmt"
   "log"
   "net/http"
)

func (a *App) LoginSubmit(username, password string) error {
   url := "https://m7cp.login.solocoo.tv/login"

   payload := map[string]interface{}{
      "ticket": a.Ticket,
      "userInput": map[string]string{
         "username": username,
         "password": password,
      },
   }

   body, _ := json.Marshal(payload)
   req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
   if err != nil {
      return err
   }
   setCommonHeaders(req)
   req.Header.Set("Content-Type", "application/json")
   // Using hardcoded Client Auth signature from HAR
   req.Header.Set("Authorization", "Client key=web.NhFyz4KsZ54,time=1775520341,sig=GxHEdRx_fczydn8Y3dTiwUYZxjs62EhtMyL4jXzd6LE")

   resp, err := a.Client.Do(req)
   if err != nil {
      return err
   }
   defer resp.Body.Close()

   if resp.StatusCode != 200 {
      return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
   }

   var result struct {
      Label    string `json:"label"`
      Result   string `json:"result"`
      SsoToken string `json:"ssoToken"`
   }
   if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
      return err
   }

   if result.Result == "success" {
      log.Println("Login successful! New SSO Token obtained.")
      a.SsoToken = result.SsoToken
   } else {
      log.Printf("Login response: %s", result.Result)
   }

   return nil
}
