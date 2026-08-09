package cineMember

import (
   "encoding/json"
   "errors"
   "fmt"
   "io"
   "log"
   "net/http"
   "net/url"
   "strconv"
   "strings"
)

// extracts the numeric ID and converts it to an integer
func FetchId(address string) (int, error) {
   req, err := http.NewRequest("GET", address, nil)
   if err != nil {
      return 0, err
   }
   resp, err := do(req)
   if err != nil {
      return 0, err
   }
   defer resp.Body.Close()
   var data strings.Builder
   _, err = io.Copy(&data, resp.Body)
   if err != nil {
      return 0, err
   }
   // 1. Cut text after "app.play('"
   _, after, found := strings.Cut(data.String(), "app.play('")
   if !found {
      return 0, errors.New("start marker not found")
   }
   // 2. Cut text at the next single quote to isolate the ID string
   idStr, _, found := strings.Cut(after, "'")
   if !found {
      return 0, errors.New("closing quote not found")
   }
   // 3. Convert string to integer
   return strconv.Atoi(idStr)
}

func FetchLogin(phpSessId *Cookie, email, password string) error {
   body := url.Values{
      "emaillogin": {email},
      "password":   {password},
   }.Encode()
   req, err := http.NewRequest(
      "POST",
      "https://www.cinemember.nl/elements/overlays/account/login.php",
      strings.NewReader(body),
   )
   if err != nil {
      return err
   }
   req.Header.Set("content-type", "application/x-www-form-urlencoded")
   req.Header.Set("cookie", phpSessId.String())
   resp, err := do(req)
   if err != nil {
      return err
   }
   defer resp.Body.Close()
   _, err = io.Copy(io.Discard, resp.Body)
   return err
}

func do(req *http.Request) (*http.Response, error) {
   log.Println(req.Method, req.URL)
   return http.DefaultClient.Do(req)
}

type Cookie struct {
   Name  string
   Value string
}

func GetPhpSessId() (*Cookie, error) {
   req, err := http.NewRequest("HEAD", "https://www.cinemember.nl/nl", nil)
   if err != nil {
      return nil, err
   }
   // THIS IS NEEDED OTHERWISE SUBTITLES ARE MISSING, GOD IS DEAD
   req.Header.Set("user-agent", "Windows")
   resp, err := do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   for _, each := range resp.Cookies() {
      if each.Name == "PHPSESSID" {
         return &Cookie{Name: each.Name, Value: each.Value}, nil
      }
   }
   return nil, errors.New("PHPSESSID cookie not found in response")
}

func (*Cookie) CachePath() string {
   return "rosso/cineMember/Cookie"
}

func (c *Cookie) String() string {
   return fmt.Sprintf("%v=%v", c.Name, c.Value)
}

type Stream struct {
   Error string
   Links []struct {
      MimeType string
      Url      string
   }
   NoAccess bool
}

// must run login first
func FetchStream(phpSessId *Cookie, id int) (*Stream, error) {
   req, err := http.NewRequest(
      "GET",
      "https://www.cinemember.nl/elements/films/stream.php",
      nil,
   )
   if err != nil {
      return nil, err
   }
   req.URL.RawQuery = fmt.Sprint("id=", id)
   req.Header.Set("cookie", phpSessId.String())
   resp, err := do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   var result Stream
   err = json.NewDecoder(resp.Body).Decode(&result)
   if err != nil {
      return nil, err
   }
   if result.Error != "" {
      return nil, errors.New(result.Error)
   }
   if result.NoAccess {
      return nil, errors.New("no access")
   }
   return &result, nil
}

func (s *Stream) GetDash() (*url.URL, error) {
   for _, link := range s.Links {
      if link.MimeType == "application/dash+xml" {
         return url.Parse(link.Url)
      }
   }
   return nil, errors.New("DASH link not found")
}

// cineMember.go
