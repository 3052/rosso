package canal

import (
   "bytes"
   "crypto/hmac"
   "crypto/sha256"
   "encoding/base64"
   "encoding/json"
   "errors"
   "fmt"
   "io"
   "net/http"
   "net/url"
   "strconv"
   "strings"
   "time"
)

func (a *Asset) String() string {
   var data strings.Builder
   data.WriteString("title = ")
   data.WriteString(a.Title)
   data.WriteString("\ntype = ")
   data.WriteString(a.Type)
   data.WriteString("\nid = ")
   data.WriteString(a.Id)
   return data.String()
}

type Asset struct {
   Title string
   Type  string
   Id    string
}

type Collection struct {
   Assets []Asset
}

type Session struct {
   Message  string
   SsoToken string
   Token    string // this last one hour
}

func (s *Session) Player(tracking string) (*Player, error) {
   data, err := json.Marshal(map[string]any{
      "player": map[string]any{
         "capabilities": map[string]any{
            "drmSystems": []string{"Widevine"},
            "mediaTypes": []string{"DASH"},
         },
      },
   })
   if err != nil {
      return nil, err
   }
   req, err := http.NewRequest(
      "POST",
      fmt.Sprintf("https://tvapi-hlm2.solocoo.tv/v1/assets/%v/play", tracking),
      bytes.NewReader(data),
   )
   if err != nil {
      return nil, err
   }
   req.Header.Set("authorization", "Bearer "+s.Token)
   req.Header.Set("content-type", "application/json")
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   var result Player
   err = json.NewDecoder(resp.Body).Decode(&result)
   if err != nil {
      return nil, err
   }
   if result.Message != "" {
      return nil, errors.New(result.Message)
   }
   return &result, nil
}

func FetchTicket() (*Ticket, error) {
   data, err := json.Marshal(map[string]any{
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
   req, err := http.NewRequest(
      "POST", "https://m7cp.login.solocoo.tv/login", bytes.NewReader(data),
   )
   if err != nil {
      return nil, err
   }
   client, err := get_client(req.URL, data)
   if err != nil {
      return nil, err
   }
   req.Header.Set("authorization", client)
   req.Header.Set("user-agent", user_agent)
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   var result Ticket
   err = json.NewDecoder(resp.Body).Decode(&result)
   if err != nil {
      return nil, err
   }
   if result.Message != "" {
      return nil, errors.New(result.Message)
   }
   return &result, nil
}

const user_agent = "Mozilla/5.0 Windows"

type Player struct {
   Drm struct {
      LicenseUrl string
   }
   Message   string
   Subtitles []struct {
      Url string
   }
   Url string // MPD
}

func (t *Ticket) Login(username, password string) (*Login, error) {
   data, err := json.Marshal(map[string]any{
      "ticket": t.Ticket,
      "userInput": map[string]string{
         "username": username,
         "password": password,
      },
   })
   if err != nil {
      return nil, err
   }
   req, err := http.NewRequest(
      "POST", "https://m7cp.login.solocoo.tv/login", bytes.NewReader(data),
   )
   if err != nil {
      return nil, err
   }
   client, err := get_client(req.URL, data)
   if err != nil {
      return nil, err
   }
   req.Header.Set("authorization", client)
   req.Header.Set("user-agent", user_agent)
   resp, err := http.DefaultClient.Do(req)
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
      return nil, &result
   }
   return &result, nil
}

type Ticket struct {
   Message string
   Ticket  string
}

func (p *Player) GetManifest() (*url.URL, error) {
   return url.Parse(p.Url)
}

type Episode struct {
   Desc   string
   Id     string
   Params struct {
      SeriesEpisode int
   }
   Title string
}

type Login struct {
   Label    string
   Message  string
   SsoToken string // this last one day
}

func (l *Login) Error() string {
   var data strings.Builder
   data.WriteString("label = ")
   data.WriteString(l.Label)
   data.WriteString("\nmessage = ")
   data.WriteString(l.Message)
   return data.String()
}

const device_serial = "!!!!"

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

func (e *Episode) String() string {
   data := &strings.Builder{}
   fmt.Fprintln(data, "episode =", e.Params.SeriesEpisode)
   fmt.Fprintln(data, "title =", e.Title)
   fmt.Fprintln(data, "desc =", e.Desc)
   fmt.Fprint(data, "tracking = ", e.Id)
   return data.String()
}
