package oldflix

import (
   "encoding/json"
   "errors"
   "fmt"
   "io"
   "net/http"
   "net/url"
   "strings"
)

func (w *Watch) FetchHls() (*Hls, error) {
   resp, err := http.Get(w.Playlist[0].File)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   body, err := io.ReadAll(resp.Body)
   if err != nil {
      return nil, err
   }
   return &Hls{Body: body, Url: resp.Request.URL}, nil
}

type Watch struct {
   Message  string
   Playlist []struct {
      File string
   }
}

const BaseUrl = "https://oldflix-api.azurewebsites.net"

func (b *Browse) FetchWatch(trackId, token string) (*Watch, error) {
   data := url.Values{
      "id": {b.Id},
      "m":  {b.Movie.Id},
      "tk": {trackId}, // tk is the audio/language track id
   }.Encode()
   req, err := http.NewRequest(
      "POST", BaseUrl+"/api/watch/play", strings.NewReader(data),
   )
   if err != nil {
      return nil, err
   }
   req.Header.Set("Authorization", "Bearer "+token)
   req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
   resp, err := http.DefaultClient.Do(req)
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
         return &track_data, nil
      }
   }
   return nil, errors.New("track with language 'Original' not found")
}

type Browse struct {
   Id    string
   Movie struct {
      Id     string
      Tracks []Track
   }
}

type Hls struct {
   Body []byte
   Url  *url.URL
}

// https://oldflix.com.br/browse/play/5d5d54a4d55dc050f8468513
func (l *Login) FetchBrowse(contentId string) (*Browse, error) {
   data := url.Values{"id": {contentId}}.Encode()
   req, err := http.NewRequest(
      "POST", BaseUrl+"/api/browse/play", strings.NewReader(data),
   )
   if err != nil {
      return nil, err
   }
   req.Header.Set("authorization", "Bearer "+l.Token)
   req.Header.Set("content-type", "application/x-www-form-urlencoded")
   resp, err := http.DefaultClient.Do(req)
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

func FetchLogin(username, password string) (*Login, error) {
   resp, err := http.PostForm(BaseUrl+"/api/token", url.Values{
      "password": {password},
      "username": {username},
   })
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

type Login struct {
   Status int
   Token  string
}

type Track struct {
   Id   string
   Lang string
   Lnk  string
}
