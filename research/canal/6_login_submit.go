package main

import (
   "bytes"
   "encoding/json"
   "fmt"
   "log"
   "net/http"
   "time"
)

type UserInput struct {
   Username string `json:"username"`
   Password string `json:"password"`
}

type LoginSubmitPayload struct {
   Ticket    string    `json:"ticket"`
   UserInput UserInput `json:"userInput"`
}

func (a *App) LoginSubmit(username, password string) error {
   url := "https://m7cp.login.solocoo.tv/login"

   payload := LoginSubmitPayload{
      Ticket: a.Ticket,
      UserInput: UserInput{
         Username: username,
         Password: password,
      },
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
      Label    string `json:"label"`
      Result   string `json:"result"`
      SsoToken string `json:"ssoToken"`
   }
   if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
      return err
   }

   if result.Result == "success" {
      log.Println("Login successful! New Authed SSO Token obtained.")
      a.SsoToken = result.SsoToken // You are now fully logged in
   } else {
      log.Printf("Login response: %s", result.Result)
   }

   return nil
}
