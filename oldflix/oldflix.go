package oldflix

import (
   "bytes"
   "encoding/json"
   "errors"
   "fmt"
   "io"
   "log"
   "net/http"
   "net/url"
)

const azure = "oldflix-api.azurewebsites.net"

// do sends the request, logs the method and URL, and returns the response.
func do(req *http.Request) (*http.Response, error) {
   log.Println(req.Method, req.URL)
   return http.DefaultClient.Do(req)
}

type Browse struct {
   Id    string
   Movie *struct {
      Id     string
      Tracks []*Track
   }
}

func (b *Browse) FetchWatch(trackId, token string) (*Watch, error) {
   body := url.Values{
      "id": {b.Id},
      "m":  {b.Movie.Id},
      "tk": {trackId}, // tk is the audio/language track id
   }.Encode()
   req, err := http.NewRequest(
      "POST",
      (&url.URL{
         Scheme: "https",
         Host:   azure,
         Path:   "/api/watch/play",
      }).String(),
      bytes.NewReader([]byte(body)),
   )
   if err != nil {
      return nil, err
   }
   req.Header.Set("authorization", "Bearer "+token)
   req.Header.Set("content-type", "application/x-www-form-urlencoded")
   resp, err := do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   if resp.StatusCode != http.StatusOK {
      return nil, errors.New(resp.Status)
   }
   var result Watch
   if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
      return nil, fmt.Errorf("failed to decode watch play response: %w", err)
   }
   if result.Message != "" {
      return nil, errors.New(result.Message)
   }
   return &result, nil
}

func (b *Browse) GetOriginal() (*Track, error) {
   for _, track_data := range b.Movie.Tracks {
      if track_data.Lang == "Original" {
         return track_data, nil
      }
   }
   return nil, errors.New("track with language 'Original' not found")
}

type Login struct {
   Status int
   Token  string
}

func FetchLogin(username, password string) (*Login, error) {
   body := url.Values{
      "password": {password},
      "username": {username},
   }.Encode()
   req, err := http.NewRequest(
      "POST",
      (&url.URL{
         Scheme: "https",
         Host:   azure,
         Path:   "/api/token",
      }).String(),
      bytes.NewReader([]byte(body)),
   )
   if err != nil {
      return nil, err
   }
   req.Header.Set("content-type", "application/x-www-form-urlencoded")
   resp, err := do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   data, err := io.ReadAll(resp.Body)
   if err != nil {
      return nil, err
   }
   if resp.StatusCode != http.StatusOK {
      return nil, errors.New(string(data))
   }
   result := &Login{}
   err = json.Unmarshal(data, result)
   if err != nil {
      return nil, fmt.Errorf("failed to decode login response: %w", err)
   }
   return result, nil
}

func (*Login) CachePath() string {
   return "rosso/oldflix/Login"
}

// https://oldflix.com.br/browse/play/5d5d54a4d55dc050f8468513
func (l *Login) FetchBrowse(id string) (*Browse, error) {
   body := url.Values{"id": {id}}.Encode()
   req, err := http.NewRequest(
      "POST",
      (&url.URL{
         Scheme: "https",
         Host:   azure,
         Path:   "/api/browse/play",
      }).String(),
      bytes.NewReader([]byte(body)),
   )
   if err != nil {
      return nil, err
   }
   req.Header.Set("authorization", "Bearer "+l.Token)
   req.Header.Set("content-type", "application/x-www-form-urlencoded")
   resp, err := do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   if resp.StatusCode != http.StatusOK {
      return nil, errors.New(resp.Status)
   }
   result := &Browse{}
   err = json.NewDecoder(resp.Body).Decode(result)
   if err != nil {
      return nil, fmt.Errorf("failed to decode browse play response: %w", err)
   }
   return result, nil
}

type Track struct {
   Id   string
   Lang string
   Lnk  string
}

type Watch struct {
   Message  string
   Playlist []*struct {
      File string
   }
}

// oldflix.go
