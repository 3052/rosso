package mubi

import (
   "encoding/json"
   "errors"
   "fmt"
   "net/http"
   "net/url"
   "strings"
)

func (s *SecureUrl) GetManifest() (*url.URL, error) {
   s.Url = strings.NewReplacer(
      ".AVC1", "",
      ".ex-eac3", "",
      ".ex-vtt", "",
   ).Replace(s.Url)
   return url.Parse(s.Url)
}

type SecureUrl struct {
   TextTrackUrls []struct {
      Id  string
      Url string
   } `json:"text_track_urls"`
   Url         string // MPD
   UserMessage string `json:"user_message"`
}

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

type Session struct {
   Token string
   User  struct {
      Id int
   }
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
