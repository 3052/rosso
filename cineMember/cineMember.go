package cineMember

import (
   "encoding/json"
   "errors"
   "io"
   "net/http"
   "net/url"
   "strconv"
   "strings"
)

// must run login first
func FetchStream(session *http.Cookie, id int) (*Stream, error) {
   var req http.Request
   req.URL = &url.URL{
      Scheme:   "https",
      Host:     "www.cinemember.nl",
      Path:     "/elements/films/stream.php",
      RawQuery: "id=" + strconv.Itoa(id),
   }
   req.Header = http.Header{}
   req.AddCookie(session)
   resp, err := http.DefaultClient.Do(&req)
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

type Stream struct {
   Error    string
   Links    []MediaLink
   NoAccess bool
}

func (s *Stream) Dash() (*MediaLink, error) {
   for i := range s.Links {
      if s.Links[i].MimeType == "application/dash+xml" {
         return &s.Links[i], nil
      }
   }
   // Create and return the error directly.
   return nil, errors.New("DASH link not found")
}

// extracts the numeric ID and converts it to an integer
func FetchId(link string) (int, error) {
   resp, err := http.Get(link)
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

type Dash struct {
   Body []byte
   Url  *url.URL
}

type MediaLink struct {
   MimeType string
   Url      string
}

func (m *MediaLink) Dash() (*Dash, error) {
   resp, err := http.Get(m.Url)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   var result Dash
   result.Body, err = io.ReadAll(resp.Body)
   if err != nil {
      return nil, err
   }
   result.Url = resp.Request.URL
   return &result, nil
}

func FetchSession() (*http.Cookie, error) {
   // We ignore the error here because the method and URL are hardcoded and
   // known to be valid.
   var req http.Request
   req.Method = "HEAD"
   req.URL = &url.URL{
      Scheme: "https",
      Host:   "www.cinemember.nl",
      Path:   "/nl",
   }
   req.Header = http.Header{}
   // THIS IS NEEDED OTHERWISE SUBTITLES ARE MISSING, GOD IS DEAD
   req.Header.Set("User-Agent", "Windows")
   resp, err := http.DefaultClient.Do(&req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   for _, cookie := range resp.Cookies() {
      if cookie.Name == "PHPSESSID" {
         return cookie, nil
      }
   }
   return nil, errors.New("PHPSESSID cookie not found in response")
}

func FetchLogin(session *http.Cookie, email, password string) error {
   data := url.Values{
      "emaillogin": {email},
      "password":   {password},
   }.Encode()
   req, err := http.NewRequest(
      "POST", "https://www.cinemember.nl/elements/overlays/account/login.php",
      strings.NewReader(data),
   )
   if err != nil {
      return err
   }
   req.AddCookie(session)
   req.Header.Set("content-type", "application/x-www-form-urlencoded")
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return err
   }
   defer resp.Body.Close()
   _, err = io.Copy(io.Discard, resp.Body)
   return err
}
