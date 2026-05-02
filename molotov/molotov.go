package molotov

import (
   "41.neocities.org/maya"
   "encoding/json"
   "errors"
   "fmt"
   "net/url"
   "strconv"
   "strings"
)

type Auth struct {
   AccessToken  string `json:"access_token"`
   RefreshToken string `json:"refresh_token"`
}

// authorization server issues a new refresh token, in which case the
// client MUST discard the old refresh token and replace it with the new
// refresh token
func (a *Auth) Refresh() error {
   resp, err := maya.Get(
      &url.URL{
         Scheme: "https",
         Host:   "fapi.molotov.tv",
         Path:   "/v3/auth/refresh/" + a.RefreshToken,
      },
      map[string]string{"x-molotov-agent": customer_area},
   )
   if err != nil {
      return err
   }
   defer resp.Body.Close()
   return json.NewDecoder(resp.Body).Decode(a)
}

func (a *Auth) FetchPlay(programData *Program) (*Play, error) {
   resp, err := maya.Get(
      &url.URL{
         Scheme: "https",
         Host:   "fapi.molotov.tv",
         Path: fmt.Sprintf(
            "/v2/channels/%v/programs/%v/view",
            programData.ChannelId, programData.Id,
         ),
         RawQuery: url.Values{"access_token": {a.AccessToken}}.Encode(),
      },
      map[string]string{"x-molotov-agent": customer_area},
   )
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   var result struct {
      Program struct {
         Actions struct {
            Play *Play
         }
      }
   }
   err = json.NewDecoder(resp.Body).Decode(&result)
   if err != nil {
      return nil, err
   }
   if result.Program.Actions.Play == nil {
      return nil, errors.New("program is not available for playback")
   }
   return result.Program.Actions.Play, nil
}

func FetchAuth(email, password string) (*Auth, error) {
   body, err := json.Marshal(map[string]string{
      "grant_type": "password",
      "email":      email,
      "password":   password,
   })
   if err != nil {
      return nil, err
   }
   resp, err := maya.Post(
      &url.URL{
         Scheme: "https", Host: "fapi.molotov.tv", Path: "/v3.1/auth/login",
      },
      map[string]string{"x-molotov-agent": customer_area},
      body,
   )
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   var result struct {
      Auth Auth
   }
   err = json.NewDecoder(resp.Body).Decode(&result)
   if err != nil {
      return nil, err
   }
   return &result.Auth, nil
}

func (a *Asset) FetchWidevine(body []byte) ([]byte, error) {
   resp, err := maya.Post(
      &url.URL{
         Scheme: "https",
         Host:   "lic.drmtoday.com",
         Path:   "/license-proxy-widevine/cenc/",
      },
      map[string]string{"x-dt-auth-token": a.Drm.Token},
      body,
   )
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   if resp.StatusCode != 200 {
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

func (a *Asset) GetManifest() (*url.URL, error) {
   return url.Parse(strings.Replace(a.Stream.Url, "high", "fhdready", 1))
}

const (
   browser_app   = `{ "app_build": 4, "app_id": "browser_app", "inner_app_version_name": "5.7.0" }`
   customer_area = `{ "app_build": 1, "app_id": "customer_area" }`
)

func (e *Error) Error() string {
   var data strings.Builder
   data.WriteString("developer message = ")
   data.WriteString(e.DeveloperMessage)
   data.WriteString("\nuser message = ")
   data.WriteString(e.UserMessage)
   return data.String()
}

type Error struct {
   DeveloperMessage string `json:"developer_message"`
   UserMessage      string `json:"user_message"`
}

// https://molotov.tv/fr_fr/p/15301-2328
// https://molotov.tv/fr_fr/p/15301-2328/closer-entre-adultes-consentants
func ParseProgram(data string) (*Program, error) {
   var found bool
   _, data, found = strings.Cut(data, "/p/")
   if !found {
      return nil, errors.New("url does not contain the /p/ marker")
   }
   data, _, _ = strings.Cut(data, "/")
   id, channel_id, found := strings.Cut(data, "-")
   if !found {
      return nil, errors.New("invalid format: hyphen not found between IDs")
   }
   var (
      p   Program
      err error
   )
   if p.Id, err = strconv.Atoi(id); err != nil {
      return nil, errors.New("program ID is not a valid integer")
   }
   if p.ChannelId, err = strconv.Atoi(channel_id); err != nil {
      return nil, errors.New("channel ID is not a valid integer")
   }
   return &p, nil
}

type Program struct {
   Id        int
   ChannelId int
}

type Asset struct {
   Drm struct {
      Token string
   }
   Error  *Error
   Stream struct {
      Url string // MPD
   }
}

type Play struct {
   Url string // fapi.molotov.tv/v2/me/assets
}

func (a *Auth) FetchAsset(playData *Play) (*Asset, error) {
   target, err := url.Parse(playData.Url)
   if err != nil {
      return nil, err
   }
   query := target.Query() // keep existing query string
   query.Set("access_token", a.AccessToken)
   target.RawQuery = query.Encode()
   resp, err := maya.Get(
      target,
      map[string]string{
         "x-forwarded-for": "138.199.15.158",
         "x-molotov-agent": browser_app,
      },
   )
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
