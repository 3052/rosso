package oldflix

import (
   "41.neocities.org/maya"
   "encoding/json"
   "errors"
   "fmt"
   "io"
   "net/url"
)

func FetchLogin(username, password string) (*Login, error) {
   body := url.Values{
      "password": {password},
      "username": {username},
   }.Encode()
   resp, err := maya.Post(
      &url.URL{
         Scheme: "https",
         Host:   azure,
         Path:   "/api/token",
      },
      map[string]string{"content-type": "application/x-www-form-urlencoded"},
      []byte(body),
   )
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   data, err := io.ReadAll(resp.Body)
   if err != nil {
      return nil, err
   }
   if resp.StatusCode != 200 {
      return nil, errors.New(string(data))
   }
   result := &Login{}
   err = json.Unmarshal(data, result)
   if err != nil {
      return nil, fmt.Errorf("failed to decode login response: %w", err)
   }
   return result, nil
}

// https://oldflix.com.br/browse/play/5d5d54a4d55dc050f8468513
func (l *Login) FetchBrowse(contentId string) (*Browse, error) {
   body := url.Values{"id": {contentId}}.Encode()
   resp, err := maya.Post(
      &url.URL{
         Scheme: "https",
         Host:   azure,
         Path:   "/api/browse/play",
      },
      map[string]string{
         "authorization": "Bearer " + l.Token,
         "content-type":  "application/x-www-form-urlencoded",
      },
      []byte(body),
   )
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   if resp.StatusCode != 200 {
      return nil, errors.New(resp.Status)
   }
   result := &Browse{}
   err = json.NewDecoder(resp.Body).Decode(result)
   if err != nil {
      return nil, fmt.Errorf("failed to decode browse play response: %w", err)
   }
   return result, nil
}

func (b *Browse) FetchWatch(trackId, token string) (*Watch, error) {
   body := url.Values{
      "id": {b.Id},
      "m":  {b.Movie.Id},
      "tk": {trackId}, // tk is the audio/language track id
   }.Encode()
   resp, err := maya.Post(
      &url.URL{
         Scheme: "https",
         Host:   azure,
         Path:   "/api/watch/play",
      },
      map[string]string{
         "authorization": "Bearer " + token,
         "content-type":  "application/x-www-form-urlencoded",
      },
      []byte(body),
   )
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   if resp.StatusCode != 200 {
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

func (w *Watch) GetManifest() (*url.URL, error) {
   return url.Parse(w.Playlist[0].File)
}

type Watch struct {
   Message  string
   Playlist []struct {
      File string
   }
}

const azure = "oldflix-api.azurewebsites.net"

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

type Login struct {
   Status int
   Token  string
}

type Track struct {
   Id   string
   Lang string
   Lnk  string
}
