package mubi

import (
   "bytes"
   "encoding/base64"
   "encoding/json"
   "errors"
   "fmt"
   "log"
   "net/http"
   "net/url"
   "strings"
)

// "android" requires headers:
// client-device-identifier
// client-version
const client = "web"

var ClientCountry = "US"

func do(req *http.Request) (*http.Response, error) {
   log.Println(req.Method, req.URL)
   return http.DefaultClient.Do(req)
}

type Film struct {
   Title string
   Id    int
}

func FetchEpisodes(slug string, season int) ([]*Film, error) {
   req, err := http.NewRequest("GET", fmt.Sprintf("https://api.mubi.com/v4/series/%v/seasons/season-%v/episodes", slug, season), nil)
   if err != nil {
      return nil, err
   }
   req.Header.Set("client", client)
   req.Header.Set("client-country", ClientCountry)
   resp, err := do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   var result struct {
      Episodes []*Film
   }
   err = json.NewDecoder(resp.Body).Decode(&result)
   if err != nil {
      return nil, err
   }
   return result.Episodes, nil
}

func FetchFilm(slug string) (*Film, error) {
   req, err := http.NewRequest("GET", "https://api.mubi.com/v3/films/"+slug, nil)
   if err != nil {
      return nil, err
   }
   req.Header.Set("client", client)
   req.Header.Set("client-country", ClientCountry)
   resp, err := do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   result := &Film{}
   err = json.NewDecoder(resp.Body).Decode(result)
   if err != nil {
      return nil, err
   }
   return result, nil
}

func (f *Film) String() string {
   data := &strings.Builder{}
   data.WriteString("title: ")
   data.WriteString(f.Title)
   data.WriteString("\nid: ")
   fmt.Fprint(data, f.Id)
   return data.String()
}

type LinkCode struct {
   AuthToken string `json:"auth_token"`
   LinkCode  string `json:"link_code"`
}

func FetchLinkCode() (*LinkCode, error) {
   req, err := http.NewRequest("GET", "https://api.mubi.com/v3/link_code", nil)
   if err != nil {
      return nil, err
   }
   req.Header.Set("client", client)
   req.Header.Set("client-country", ClientCountry)
   resp, err := do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   result := &LinkCode{}
   err = json.NewDecoder(resp.Body).Decode(result)
   if err != nil {
      return nil, err
   }
   return result, nil
}

func (*LinkCode) CachePath() string {
   return "rosso/mubi/LinkCode"
}

func (l *LinkCode) FetchSession() (*Session, error) {
   body, err := json.Marshal(map[string]string{"auth_token": l.AuthToken})
   if err != nil {
      return nil, err
   }
   req, err := http.NewRequest("POST", "https://api.mubi.com/v3/authenticate", bytes.NewReader(body))
   if err != nil {
      return nil, err
   }
   req.Header.Set("client", client)
   req.Header.Set("client-country", ClientCountry)
   req.Header.Set("content-type", "application/json")
   resp, err := do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   result := &Session{}
   err = json.NewDecoder(resp.Body).Decode(result)
   if err != nil {
      return nil, err
   }
   return result, nil
}

func (l *LinkCode) String() string {
   var data strings.Builder
   data.WriteString("TO LOG IN AND START WATCHING\n")
   data.WriteString("Go to\n")
   data.WriteString("mubi.com/en/android\n")
   data.WriteString("and enter the code below\n")
   data.WriteString(l.LinkCode)
   return data.String()
}

type SecureUrl struct {
   TextTrackUrls []struct {
      Id  string
      Url string
   } `json:"text_track_urls"`
   Url         string // MPD
   UserMessage string `json:"user_message"`
}

func (s *SecureUrl) GetManifest() (*url.URL, error) {
   manifest, err := url.Parse(s.Url)
   if err != nil {
      return nil, err
   }
   manifest.Path = strings.NewReplacer(
      ".AVC1", "",
      ".ex-eac3", "",
      ".ex-vtt", "",
   ).Replace(manifest.Path)
   return manifest, nil
}

type Session struct {
   Token string
   User  struct {
      Id int
   }
}

func (*Session) CachePath() string {
   return "rosso/mubi/Session"
}

func (s *Session) FetchSecureUrl(id int) (*SecureUrl, error) {
   req, err := http.NewRequest("GET", fmt.Sprintf("https://api.mubi.com/v3/films/%v/viewing/secure_url", id), nil)
   if err != nil {
      return nil, err
   }
   req.Header.Set("authorization", "Bearer "+s.Token)
   req.Header.Set("client", client)
   req.Header.Set("client-country", ClientCountry)
   req.Header.Set("user-agent", "Firefox")
   resp, err := do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   var result SecureUrl
   err = json.NewDecoder(resp.Body).Decode(&result)
   if err != nil {
      return nil, err
   }
   if result.UserMessage != "" {
      return nil, errors.New(result.UserMessage)
   }
   return &result, nil
}

// to get the MPD you have to call this or view video on the website. request
// is hard geo blocked only the first time
func (s *Session) FetchViewing(id int) error {
   req, err := http.NewRequest("POST", fmt.Sprintf("https://api.mubi.com/v3/films/%v/viewing", id), nil)
   if err != nil {
      return err
   }
   req.Header.Set("authorization", "Bearer "+s.Token)
   req.Header.Set("client", client)
   req.Header.Set("client-country", ClientCountry)
   resp, err := do(req)
   if err != nil {
      return err
   }
   defer resp.Body.Close()
   var result struct {
      UserMessage string `json:"user_message"`
   }
   err = json.NewDecoder(resp.Body).Decode(&result)
   if err != nil {
      return err
   }
   if result.UserMessage != "" {
      return errors.New(result.UserMessage)
   }
   return nil
}

func (s *Session) FetchWidevine(body []byte) ([]byte, error) {
   data, err := json.Marshal(map[string]any{
      "merchant":  "mubi",
      "sessionId": s.Token,
      "userId":    s.User.Id,
   })
   if err != nil {
      return nil, err
   }
   req, err := http.NewRequest("POST", "https://lic.drmtoday.com/license-proxy-widevine/cenc/", bytes.NewReader(body))
   if err != nil {
      return nil, err
   }
   req.Header.Set("dt-custom-data", base64.StdEncoding.EncodeToString(data))
   resp, err := do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   if resp.StatusCode != http.StatusOK {
      return nil, errors.New(resp.Status)
   }
   var result struct {
      License []byte
   }
   err = json.NewDecoder(resp.Body).Decode(&result)
   if err != nil {
      return nil, err
   }
   return result.License, nil
}

// mubi.go
