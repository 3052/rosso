package mubi

import (
   "bytes"
   "encoding/base64"
   "encoding/json"
   "errors"
   "fmt"
   "io"
   "net/http"
   "net/url"
   "strconv"
   "strings"
)

type Film struct {
   Id   int
   Slug string
}

func (f *Film) FetchId() error {
   var req http.Request
   req.URL = &url.URL{
      Scheme: "https",
      Host:   "api.mubi.com",
      Path:   "/v3/films/" + f.Slug,
   }
   req.Header = http.Header{}
   req.Header.Set("client", client)
   req.Header.Set("client-country", ClientCountry)
   resp, err := http.DefaultClient.Do(&req)
   if err != nil {
      return err
   }
   defer resp.Body.Close()
   return json.NewDecoder(resp.Body).Decode(f)
}

// https://mubi.com/films/346537
// https://mubi.com/en/films/346537
// https://mubi.com/films/346537/player
// https://mubi.com/en/films/346537/player
// https://mubi.com/films/fallen-leaves-2023
// https://mubi.com/en/films/fallen-leaves-2023
// https://mubi.com/us/films/fallen-leaves-2023
// https://mubi.com/en/us/films/fallen-leaves-2023
func ParseFilm(data string) (*Film, error) {
   url_data, err := url.Parse(data)
   if err != nil {
      return nil, err
   }
   if url_data.Host != "mubi.com" {
      return nil, errors.New("not a valid mubi URL")
   }
   parts := strings.Split(url_data.Path, "/")
   for i, part := range parts {
      if part == "films" && i+1 < len(parts) {
         film := &Film{}
         identifier := parts[i+1]
         film.Id, err = strconv.Atoi(identifier)
         if err != nil {
            film.Slug = identifier
         }
         return film, nil
      }
   }
   return nil, errors.New("film identifier not found in URL")
}

type Session struct {
   Token string
   User  struct {
      Id int
   }
}

func (s *Session) Widevine(data []byte) ([]byte, error) {
   // final slash is needed
   req, err := http.NewRequest(
      "POST", "https://lic.drmtoday.com/license-proxy-widevine/cenc/",
      bytes.NewReader(data),
   )
   if err != nil {
      return nil, err
   }

   data, err = json.Marshal(map[string]any{
      "merchant":  "mubi",
      "sessionId": s.Token,
      "userId":    s.User.Id,
   })
   if err != nil {
      return nil, err
   }

   req.Header.Set("dt-custom-data", base64.StdEncoding.EncodeToString(data))

   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   // Check if the response is not a 200 OK
   if resp.StatusCode != http.StatusOK {
      return nil, fmt.Errorf("unexpected HTTP error %v", resp.StatusCode)
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

///

// to get the MPD you have to call this or view video on the website. request
// is hard geo blocked only the first time
func (s *Session) Viewing(filmId int) error {
   var req http.Request
   req.Header = http.Header{}
   req.Header.Set("authorization", "Bearer "+s.Token)
   req.Header.Set("client", client)
   req.Header.Set("client-country", ClientCountry)
   req.Method = "POST"
   req.URL = &url.URL{
      Scheme: "https",
      Host:   "api.mubi.com",
      Path:   fmt.Sprintf("/v3/films/%v/viewing", filmId),
   }
   resp, err := http.DefaultClient.Do(&req)
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

func FetchLinkCode() (*LinkCode, error) {
   var req http.Request
   req.URL = &url.URL{
      Scheme: "https",
      Host:   "api.mubi.com",
      Path:   "/v3/link_code",
   }
   req.Header = http.Header{}
   req.Header.Set("client", client)
   req.Header.Set("client-country", ClientCountry)
   resp, err := http.DefaultClient.Do(&req)
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

func (s *SecureUrl) Dash() (*Dash, error) {
   s.Url = strings.NewReplacer(
      ".AVC1", "",
      ".ex-eac3", "",
      ".ex-vtt", "",
   ).Replace(s.Url)
   resp, err := http.Get(s.Url)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   body, err := io.ReadAll(resp.Body)
   if err != nil {
      return nil, err
   }
   return &Dash{Body: body, Url: resp.Request.URL}, nil
}

func (s *Session) SecureUrl(filmId int) (*SecureUrl, error) {
   var req http.Request
   req.URL = &url.URL{
      Scheme: "https",
      Host:   "api.mubi.com",
      Path:   fmt.Sprintf("/v3/films/%v/viewing/secure_url", filmId),
   }
   req.Header = http.Header{}
   req.Header.Set("authorization", "Bearer "+s.Token)
   req.Header.Set("client", client)
   req.Header.Set("client-country", ClientCountry)
   req.Header.Set("user-agent", "Firefox")
   resp, err := http.DefaultClient.Do(&req)
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

type SecureUrl struct {
   TextTrackUrls []struct {
      Id  string
      Url string
   } `json:"text_track_urls"`
   Url         string // MPD
   UserMessage string `json:"user_message"`
}

type Dash struct {
   Body []byte
   Url  *url.URL
}

// "android" requires headers:
// client-device-identifier
// client-version
const client = "web"

var ClientCountry = "US"

type LinkCode struct {
   AuthToken string `json:"auth_token"`
   LinkCode  string `json:"link_code"`
}

func (l *LinkCode) String() string {
   var data strings.Builder
   data.WriteString("TO LOG IN AND START WATCHING\n")
   data.WriteString("Go to\n")
   data.WriteString("mubi.com/android\n")
   data.WriteString("and enter the code below\n")
   data.WriteString(l.LinkCode)
   return data.String()
}

func (l *LinkCode) Session() (*Session, error) {
   data, err := json.Marshal(map[string]string{"auth_token": l.AuthToken})
   if err != nil {
      return nil, err
   }
   req, err := http.NewRequest(
      "POST", "https://api.mubi.com/v3/authenticate", bytes.NewReader(data),
   )
   if err != nil {
      return nil, err
   }
   req.Header.Set("client", client)
   req.Header.Set("client-country", ClientCountry)
   req.Header.Set("content-type", "application/json")
   resp, err := http.DefaultClient.Do(req)
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
