package canal

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

func LoginSubmit(ticket, username, password string) (string, error) {
   u, err := url.Parse("https://m7cp.login.solocoo.tv/login")
   if err != nil {
      return "", err
   }
   payload := LoginSubmitPayload{
      Ticket: ticket,
      UserInput: UserInput{
         Username: username,
         Password: password,
      },
   }
   body, err := json.Marshal(payload)
   if err != nil {
      return "", err
   }
   authHeader, err := get_client(u, body)
   if err != nil {
      return "", err
   }
   req, err := http.NewRequest("POST", u.String(), bytes.NewBuffer(body))
   if err != nil {
      return "", err
   }
   req.Header.Set("Authorization", authHeader)
   req.Header.Set("User-Agent", user_agent)
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return "", err
   }
   defer resp.Body.Close()
   if resp.StatusCode != 200 {
      return "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
   }
   var result struct {
      Label    string `json:"label"`
      Result   string `json:"result"`
      SsoToken string `json:"ssoToken"`
   }
   if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
      return "", err
   }
   if result.Result == "success" {
      log.Println("Login successful! Acquired final SSO token.")
      return result.SsoToken, nil
   }
   return "", fmt.Errorf("login response label: %s", result.Label)
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

type DeviceInfo struct {
   Brand        string `json:"brand"`
   DeviceModel  string `json:"deviceModel"`
   DeviceOem    string `json:"deviceOem"`
   DeviceSerial string `json:"deviceSerial"`
   DeviceType   string `json:"deviceType"`
   OsVersion    string `json:"osVersion"`
}

type LoginSubmitPayload struct {
   Ticket    string    `json:"ticket"`
   UserInput UserInput `json:"userInput"`
}

type UserInput struct {
   Username string `json:"username"`
   Password string `json:"password"`
}

type LoginInitPayload struct {
   DeviceInfo DeviceInfo `json:"deviceInfo"`
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

const device_serial = "!!!!"

func LoginInit() (string, error) {
   u, err := url.Parse("https://m7cp.login.solocoo.tv/login")
   if err != nil {
      return "", err
   }
   payload := LoginInitPayload{
      DeviceInfo: DeviceInfo{
         Brand:        "m7cp",
         DeviceModel:  "Firefox",
         DeviceOem:    "Firefox",
         DeviceSerial: device_serial,
         DeviceType:   "PC",
         OsVersion:    "Windows 10",
      },
   }
   body, err := json.Marshal(payload)
   if err != nil {
      return "", err
   }
   authHeader, err := get_client(u, body)
   if err != nil {
      return "", err
   }
   req, err := http.NewRequest("POST", u.String(), bytes.NewBuffer(body))
   if err != nil {
      return "", err
   }
   req.Header.Set("Authorization", authHeader)
   req.Header.Set("User-Agent", user_agent)
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return "", err
   }
   defer resp.Body.Close()

   if resp.StatusCode != 200 {
      return "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
   }
   var result struct {
      Ticket string `json:"ticket"`
   }
   if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
      return "", err
   }
   return result.Ticket, nil
}

const user_agent = "Mozilla/5.0 Windows"
