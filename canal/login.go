package canal

import (
   "crypto/hmac"
   "crypto/sha256"
   "encoding/base64"
   "encoding/json"
   "errors"
   "fmt"
   "io"
   "net/http"
   "strconv"
   "strings"
   "time"
)

// Global variables for authentication
const (
   client_key = "web.NhFyz4KsZ54"
   secret_key = "OXh0-pIwu3gEXz1UiJtqLPscZQot3a0q"
)

const device_serial = "!!!!"

const user_agent = "Mozilla/5.0 Windows"

func get_client(target string, body []byte) (string, error) {
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
   io.WriteString(hash, target)
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

type Asset struct {
   Title  string
   Images []struct {
      Size string
      Type string
      Url  string
   }
   Id string
}

func (a *Asset) String() string {
   var data strings.Builder
   data.WriteString("title: ")
   data.WriteString(a.Title)
   for _, image := range a.Images {
      if image.Size == "lg" {
         if image.Type == "la" {
            data.WriteString("\nimage: ")
            data.WriteString(image.Url)
         }
      }
   }
   data.WriteString("\nid: ")
   data.WriteString(a.Id)
   return data.String()
}

type Episode struct {
   Desc   string
   Id     string
   Params struct {
      SeriesEpisode int
   }
   Title string
}

func (e *Episode) String() string {
   data := &strings.Builder{}
   fmt.Fprintln(data, "episode:", e.Params.SeriesEpisode)
   fmt.Fprintln(data, "title:", e.Title)
   fmt.Fprintln(data, "desc:", e.Desc)
   fmt.Fprint(data, "tracking: ", e.Id)
   return data.String()
}

type Login struct {
   SsoToken string // this last one day
}

type Player struct {
   Drm struct {
      LicenseUrl string
   }
   Subtitles []struct {
      Url string
   }
   Url string // MPD
}

func (*Player) CachePath() string {
   return "rosso/canal/Player"
}

func (p *Player) FetchWidevine(body []byte) ([]byte, error) {
   resp, err := doPost(p.Drm.LicenseUrl, nil, body)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   return io.ReadAll(resp.Body)
}

type Ticket struct {
   Ticket string
}

func FetchTicket() (*Ticket, error) {
   body, err := json.Marshal(map[string]any{
      "deviceInfo": map[string]string{
         "brand":        "m7cp", // sg.ui.sso.fatal.internal_error
         "deviceModel":  "Firefox",
         "deviceOem":    "Firefox",
         "deviceSerial": device_serial,
         "deviceType":   "PC",
         "osVersion":    "Windows 10",
      },
   })
   if err != nil {
      return nil, err
   }
   target := "https://m7cp.login.solocoo.tv/login"
   client, err := get_client(target, body)
   if err != nil {
      return nil, err
   }
   resp, err := doPost(
      target,
      map[string]string{
         "authorization": client,
         "user-agent":    user_agent,
      },
      body,
   )
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   var result Ticket
   err = json.NewDecoder(resp.Body).Decode(&result)
   if err != nil {
      return nil, err
   }
   return &result, nil
}

func (t *Ticket) Login(username, password string) (*Login, error) {
   body, err := json.Marshal(map[string]any{
      "ticket": t.Ticket,
      "userInput": map[string]string{
         "username": username,
         "password": password,
      },
   })
   if err != nil {
      return nil, err
   }
   target := "https://m7cp.login.solocoo.tv/login"
   client, err := get_client(target, body)
   if err != nil {
      return nil, err
   }
   resp, err := doPost(
      target,
      map[string]string{
         "authorization": client,
         "user-agent":    user_agent,
      },
      body,
   )
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   var result Login
   err = json.NewDecoder(resp.Body).Decode(&result)
   if err != nil {
      return nil, err
   }
   if resp.StatusCode != http.StatusOK {
      return nil, errors.New(resp.Status)
   }
   return &result, nil
}

// login.go
