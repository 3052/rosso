package webDriver

import (
   "bytes"
   "encoding/json"
   "io"
   "net/http"
   "strings"
)

const address = "http://127.0.0.1:4444/session"

func (s *Session) New() error {
   data, err := json.Marshal(map[string]any{
      "capabilities": map[string]any{
         "alwaysMatch": map[string]any{
            "proxy": map[string]string{
               "proxyType": "manual",
               "sslProxy":  "res.proxy-seller.com:10000",
            },
         },
      },
   })
   if err != nil {
      return err
   }
   resp, err := http.Post(address, "application/json", bytes.NewReader(data))
   if err != nil {
      return err
   }
   defer resp.Body.Close()
   return json.NewDecoder(resp.Body).Decode(s)
}

// w3c.github.io/webdriver#navigate-to
func (s Session) Navigate(url string) error {
   data, err := json.Marshal(map[string]string{
      "url": url,
   })
   if err != nil {
      return err
   }
   req, err := http.NewRequest("POST", address, bytes.NewReader(data))
   if err != nil {
      return err
   }
   req.URL.Path += func() string {
      var b strings.Builder
      b.WriteByte('/')
      b.WriteString(s.Value.SessionId)
      b.WriteString("/url")
      return b.String()
   }()
   req.Header.Set("content-type", "application/json")
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return err
   }
   defer resp.Body.Close()
   _, err = io.Copy(io.Discard, resp.Body)
   if err != nil {
      return err
   }
   return nil
}

// w3c.github.io/webdriver#cookies
func (s Session) Cookie() (*Cookie, error) {
   req, _ := http.NewRequest("", address, nil)
   req.URL.Path += func() string {
      var b strings.Builder
      b.WriteByte('/')
      b.WriteString(s.Value.SessionId)
      b.WriteString("/cookie")
      return b.String()
   }()
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   cookie1 := &Cookie{}
   err = json.NewDecoder(resp.Body).Decode(cookie1)
   if err != nil {
      return nil, err
   }
   return cookie1, nil
}

// w3c.github.io/webdriver#sessions
type Session struct {
   Value struct {
      SessionId string
   }
}

type Cookie struct {
   Value []struct {
      Name  string
      Value string
   }
}
