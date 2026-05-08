package cineMember

import (
   "41.neocities.org/maya"
   "encoding/json"
   "errors"
   "fmt"
   "io"
   "net/url"
   "strconv"
   "strings"
)

type Cookie string

func PhpSessId() (Cookie, error) {
   resp, err := maya.Head(
      &url.URL{Scheme: "https", Host: "www.cinemember.nl", Path: "/nl"},
      // THIS IS NEEDED OTHERWISE SUBTITLES ARE MISSING, GOD IS DEAD
      map[string]string{"user-agent": "Windows"},
   )
   if err != nil {
      return "", err
   }
   defer resp.Body.Close()
   for _, cookie := range resp.Cookies() {
      if cookie.Name == "PHPSESSID" {
         value := fmt.Sprintf("%v=%v", cookie.Name, cookie.Value)
         return Cookie(value), nil
      }
   }
   return "", errors.New("PHPSESSID cookie not found in response")
}

func FetchLogin(phpSessId Cookie, email, password string) error {
   body := url.Values{
      "emaillogin": {email},
      "password":   {password},
   }.Encode()
   resp, err := maya.Post(
      &url.URL{
         Scheme: "https",
         Host:   "www.cinemember.nl",
         Path:   "/elements/overlays/account/login.php",
      },
      map[string]string{
         "content-type": "application/x-www-form-urlencoded",
         "cookie":       string(phpSessId),
      },
      []byte(body),
   )
   if err != nil {
      return err
   }
   defer resp.Body.Close()
   _, err = io.Copy(io.Discard, resp.Body)
   return err
}

// must run login first
func FetchStream(phpSessId Cookie, id int) (*Stream, error) {
   resp, err := maya.Get(
      &url.URL{
         Scheme:   "https",
         Host:     "www.cinemember.nl",
         Path:     "/elements/films/stream.php",
         RawQuery: fmt.Sprint("id=", id),
      },
      map[string]string{"cookie": string(phpSessId)},
   )
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

// extracts the numeric ID and converts it to an integer
func FetchId(urlData string) (int, error) {
   target, err := url.Parse(urlData)
   if err != nil {
      return 0, err
   }
   resp, err := maya.Get(target, nil)
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

type Stream struct {
   Error string
   Links []struct {
      MimeType string
      Url      string
   }
   NoAccess bool
}

func (s *Stream) Dash() (*url.URL, error) {
   for _, link := range s.Links {
      if link.MimeType == "application/dash+xml" {
         return url.Parse(link.Url)
      }
   }
   return nil, errors.New("DASH link not found")
}
