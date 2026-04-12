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
   "strings"
)

func (s *SecureUrl) FetchDash() (*Dash, error) {
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

type Session struct {
   Token string
   User  struct {
      Id int
   }
}

func (s *Session) FetchWidevine(body []byte) ([]byte, error) {
   // final slash is needed
   req, err := http.NewRequest(
      "POST", "https://lic.drmtoday.com/license-proxy-widevine/cenc/",
      bytes.NewReader(body),
   )
   if err != nil {
      return nil, err
   }
   data, err := json.Marshal(map[string]any{
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

// to get the MPD you have to call this or view video on the website. request
// is hard geo blocked only the first time
func (s *Session) FetchViewing(id int) error {
   req := http.Request{
      Method: "POST",
      URL: &url.URL{
         Scheme: "https",
         Host:   "api.mubi.com",
         Path:   fmt.Sprintf("/v3/films/%v/viewing", id),
      },
      Header: http.Header{},
   }
   req.Header.Set("authorization", "Bearer "+s.Token)
   req.Header.Set("client", client)
   req.Header.Set("client-country", ClientCountry)
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

func (s *Session) FetchSecureUrl(id int) (*SecureUrl, error) {
   req := http.Request{
      URL: &url.URL{
         Scheme: "https",
         Host:   "api.mubi.com",
         Path:   fmt.Sprintf("/v3/films/%v/viewing/secure_url", id),
      },
      Header: http.Header{},
   }
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

///

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

// "android" requires headers:
// client-device-identifier
// client-version
const client = "web"

var ClientCountry = "US"

func FetchLinkCode() (*LinkCode, error) {
   req := http.Request{
      URL: &url.URL{
         Scheme: "https",
         Host:   "api.mubi.com",
         Path:   "/v3/link_code",
      },
      Header: http.Header{},
   }
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

func (f *Film) String() string {
   data := &strings.Builder{}
   data.WriteString("title = ")
   data.WriteString(f.Title)
   data.WriteString("\nid = ")
   fmt.Fprint(data, f.Id)
   return data.String()
}

type Film struct {
   Title string
   Id    int
}

func FetchEpisodes(slug string, season int) ([]Film, error) {
   req := http.Request{
      URL: &url.URL{
         Scheme: "https",
         Host:   "api.mubi.com",
         Path:   fmt.Sprintf("/v4/series/%v/seasons/season-%v/episodes", slug, season),
      },
      Header: http.Header{},
   }
   req.Header.Set("client", client)
   req.Header.Set("client-country", ClientCountry)
   resp, err := http.DefaultClient.Do(&req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   var result struct {
      Episodes []Film
   }
   err = json.NewDecoder(resp.Body).Decode(&result)
   if err != nil {
      return nil, err
   }
   return result.Episodes, nil
}

func FetchFilm(slug string) (*Film, error) {
   req := http.Request{
      URL: &url.URL{
         Scheme: "https",
         Host:   "api.mubi.com",
         Path:   "/v3/films/" + slug,
      },
      Header: http.Header{},
   }
   req.Header.Set("client", client)
   req.Header.Set("client-country", ClientCountry)
   resp, err := http.DefaultClient.Do(&req)
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
type Dash struct {
   Body []byte
   Url  *url.URL
}

func (l *LinkCode) FetchSession() (*Session, error) {
   body, err := json.Marshal(map[string]string{"auth_token": l.AuthToken})
   if err != nil {
      return nil, err
   }
   req, err := http.NewRequest(
      "POST", "https://api.mubi.com/v3/authenticate", bytes.NewReader(body),
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

type SecureUrl struct {
   TextTrackUrls []struct {
      Id  string
      Url string
   } `json:"text_track_urls"`
   Url         string // MPD
   UserMessage string `json:"user_message"`
}
