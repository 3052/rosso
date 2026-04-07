package main

import (
   "bytes"
   "encoding/json"
   "fmt"
   "log"
   "net/http"
   "net/url"
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
   u, err := url.Parse("https://m7cp.login.solocoo.tv/login")
   if err != nil {
      return err
   }

   payload := LoginSubmitPayload{
      Ticket: a.Ticket,
      UserInput: UserInput{
         Username: username,
         Password: password,
      },
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
      Label    string `json:"label"`
      Result   string `json:"result"`
      SsoToken string `json:"ssoToken"`
   }
   if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
      return err
   }

   if result.Result == "success" {
      log.Println("Login successful! Acquired final SSO token.")
      a.SsoToken = result.SsoToken
   } else {
      log.Printf("Login response label: %s", result.Label)
   }

   return nil
}
