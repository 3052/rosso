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

func PhpSessId() (string, error) {
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
         return cookie.String(), nil
      }
   }
   return "", errors.New("PHPSESSID cookie not found in response")
}

func (m *MediaLink) GetManifest() (*url.URL, error) {
   return url.Parse(m.Url)
}

type Stream struct {
   Error    string
   Links    []MediaLink
   NoAccess bool
}

func (s *Stream) Dash() (*MediaLink, error) {
   for _, link := range s.Links {
      if link.MimeType == "application/dash+xml" {
         return &link, nil
      }
   }
   return nil, errors.New("DASH link not found")
}

type MediaLink struct {
   MimeType string
   Url      string
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

func FetchLogin(phpSessId, email, password string) error {
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
         "cookie":       phpSessId,
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
func FetchStream(phpSessId string, id int) (*Stream, error) {
   resp, err := maya.Get(
      &url.URL{
         Scheme:   "https",
         Host:     "www.cinemember.nl",
         Path:     "/elements/films/stream.php",
         RawQuery: fmt.Sprint("id=", id),
      },
      map[string]string{"cookie": phpSessId},
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
