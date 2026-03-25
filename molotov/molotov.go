package molotov

import (
   "bytes"
   "encoding/json"
   "errors"
   "fmt"
   "io"
   "net/http"
   "net/url"
   "strconv"
   "strings"
)

func (l *Login) Asset(programData *Program) (*Asset, error) {
   url_parse, err := url.Parse(programData.Actions.Play.Url)
   if err != nil {
      return nil, err
   }
   query := url_parse.Query() // keep existing query string
   query.Set("access_token", l.Auth.AccessToken)
   url_parse.RawQuery = query.Encode()
   var req http.Request
   req.URL = url_parse
   req.Header = http.Header{}
   req.Header.Set("x-forwarded-for", "138.199.15.158")
   req.Header.Set("x-molotov-agent", browser_app)
   resp, err := http.DefaultClient.Do(&req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   var result Asset
   err = json.NewDecoder(resp.Body).Decode(&result)
   if err != nil {
      return nil, err
   }
   if result.Error != nil {
      return nil, result.Error
   }
   return &result, nil
}

func (a *AssetError) Error() string {
   var data strings.Builder
   data.WriteString("developer message = ")
   data.WriteString(a.DeveloperMessage)
   data.WriteString("\nuser message = ")
   data.WriteString(a.UserMessage)
   return data.String()
}

type AssetError struct {
   DeveloperMessage string `json:"developer_message"`
   UserMessage      string `json:"user_message"`
}

type Login struct {
   Auth struct {
      AccessToken  string `json:"access_token"`
      RefreshToken string `json:"refresh_token"`
   }
}

const (
   browser_app   = `{ "app_build": 4, "app_id": "browser_app", "inner_app_version_name": "5.7.0" }`
   customer_area = `{ "app_build": 1, "app_id": "customer_area" }`
)

type Dash struct {
   Body []byte
   Url  *url.URL
}

// authorization server issues a new refresh token, in which case the
// client MUST discard the old refresh token and replace it with the new
// refresh token
func (l *Login) Refresh() error {
   var req http.Request
   req.Header = http.Header{}
   req.Header.Set("x-molotov-agent", customer_area)
   req.URL = &url.URL{
      Scheme: "https",
      Host:   "fapi.molotov.tv",
      Path:   "/v3/auth/refresh/" + l.Auth.RefreshToken,
   }
   resp, err := http.DefaultClient.Do(&req)
   if err != nil {
      return err
   }
   defer resp.Body.Close()
   return json.NewDecoder(resp.Body).Decode(l)
}

func (a *Asset) Dash() (*Dash, error) {
   resp, err := http.Get(strings.Replace(a.Stream.Url, "high", "fhdready", 1))
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

type MediaId struct {
   Program int
   Channel int
}

func ParseMediaId(urlData string) (*MediaId, error) {
   _, remainder, found := strings.Cut(urlData, "/p/")
   if !found {
      return nil, errors.New("url does not contain the /p/ marker")
   }
   ids, _, _ := strings.Cut(remainder, "/")
   program, channel, found := strings.Cut(ids, "-")
   if !found {
      return nil, errors.New("invalid format: hyphen not found between IDs")
   }
   m := &MediaId{}
   var err error
   // Assign directly to the struct fields
   if m.Program, err = strconv.Atoi(program); err != nil {
      return nil, errors.New("program ID is not a valid integer")
   }
   if m.Channel, err = strconv.Atoi(channel); err != nil {
      return nil, errors.New("channel ID is not a valid integer")
   }
   return m, nil
}

func FetchLogin(email, password string) (*Login, error) {
   data, err := json.Marshal(map[string]string{
      "grant_type": "password",
      "email":      email,
      "password":   password,
   })
   if err != nil {
      return nil, err
   }
   req, err := http.NewRequest(
      "POST", "https://fapi.molotov.tv/v3.1/auth/login",
      bytes.NewReader(data),
   )
   if err != nil {
      return nil, err
   }
   req.Header.Set("x-molotov-agent", customer_area)
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   result := &Login{}
   err = json.NewDecoder(resp.Body).Decode(result)
   if err != nil {
      return nil, err
   }
   return result, nil
}

type Program struct {
   Actions struct {
      Play *struct { // FIXME check for nil
         Url string // fapi.molotov.tv/v2/me/assets
      }
   }
}

func (l *Login) Program(media *MediaId) (*Program, error) {
   var req http.Request
   req.Header = http.Header{}
   req.Header.Set("x-molotov-agent", customer_area)
   req.URL = &url.URL{
      Scheme: "https",
      Host:   "fapi.molotov.tv",
      Path: fmt.Sprintf(
         "/v2/channels/%v/programs/%v/view", media.Channel, media.Program,
      ),
      RawQuery: url.Values{"access_token": {l.Auth.AccessToken}}.Encode(),
   }
   resp, err := http.DefaultClient.Do(&req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   var result struct {
      Program Program
   }
   err = json.NewDecoder(resp.Body).Decode(&result)
   if err != nil {
      return nil, err
   }
   if result.Program.Actions.Play == nil {
      return nil, errors.New("program is not available for playback")
   }
   return &result.Program, nil
}

type Asset struct {
   Drm struct {
      Token string
   }
   Error  *AssetError
   Stream struct {
      Url string // MPD
   }
}

func (a *Asset) Widevine(data []byte) ([]byte, error) {
   req, err := http.NewRequest(
      "POST", "https://lic.drmtoday.com/license-proxy-widevine/cenc/",
      bytes.NewReader(data),
   )
   if err != nil {
      return nil, err
   }
   req.Header.Set("x-dt-auth-token", a.Drm.Token)
   resp, err := http.DefaultClient.Do(req)
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
