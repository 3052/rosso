package main

import (
   "bytes"
   "crypto/hmac"
   "crypto/sha256"
   "encoding/base64"
   "encoding/json"
   "fmt"
   "io"
   "log"
   "net/http"
   "net/url"
   "strconv"
   "strings"
   "time"
)

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

type UserInput struct {
   Username string `json:"username"`
   Password string `json:"password"`
}

type LoginSubmitPayload struct {
   Ticket    string    `json:"ticket"`
   UserInput UserInput `json:"userInput"`
}

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

func main() {
   app := &App{
      Client: &http.Client{
         Timeout: 15 * time.Second,
      },
      DeviceSerial: "w76d15b90-3215-11f1-87ca-01f0af932fb7", // Replace with a dynamic UUID if needed
   }
   app.TVApiBaseURL = "https://tvapi-hlm2.solocoo.tv"

   log.Println("5/6 Initializing Login (Fetching Ticket)...")
   if err := app.LoginInit(); err != nil {
      log.Fatalf("LoginInit failed: %v", err)
   }
   log.Println("6/6 Submitting Credentials...")
   // Feed your credentials here
   if err := app.LoginSubmit("27@riseup.net", "***REMOVED***"); err != nil {
      log.Fatalf("LoginSubmit failed: %v", err)
   }
   log.Println("Successfully logged in!")
}

// setCommonHeaders applies headers that are present across the requests
func setCommonHeaders(req *http.Request) {
   req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:140.0) Gecko/20100101 Firefox/140.0")
   req.Header.Set("Accept", "application/json, text/plain, */*")
   req.Header.Set("Accept-Language", "en-US,en;q=0.5")
   req.Header.Set("Origin", "https://play.canalplus.cz")
   req.Header.Set("Referer", "https://play.canalplus.cz/")
}

// Global variables for authentication
const (
   client_key = "web.NhFyz4KsZ54"
   secret_key = "OXh0-pIwu3gEXz1UiJtqLPscZQot3a0q"
)

func get_client(url_data *url.URL, body []byte) (string, error) {
   encoding := base64.RawURLEncoding
   // 1. base64 raw URL decode secret key
   decoded_key, err := encoding.DecodeString(secret_key)
   if err != nil {
      return "", err
   }
   // Prepare timestamp as string immediately
   timestamp := strconv.FormatInt(time.Now().Unix(), 10)
   body_checksum := sha256.Sum256(body)
   encoded_body_hash := encoding.EncodeToString(body_checksum[:])
   // 2. hmac.New(sha256.New, secret key)
   hash := hmac.New(sha256.New, decoded_key)
   // 3, 4, 5. Write components to the hasher
   io.WriteString(hash, url_data.String())
   io.WriteString(hash, encoded_body_hash)
   io.WriteString(hash, timestamp)
   // 6. base64 raw URL encode the hmac sum
   signature := encoding.EncodeToString(hash.Sum(nil))
   // Construct final result string using strings.Builder
   var data strings.Builder
   data.WriteString("Client key=")
   data.WriteString(client_key)
   data.WriteString(",time=")
   data.WriteString(timestamp)
   data.WriteString(",sig=")
   data.WriteString(signature)
   return data.String(), nil
}

type App struct {
   Client        *http.Client
   DeviceSerial  string
   TVApiBaseURL  string
   ProvisionData string
   SsoToken      string
   BearerToken   string
   Ticket        string
}
