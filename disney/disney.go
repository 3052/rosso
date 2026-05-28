package disney

import (
   "41.neocities.org/maya"
   _ "embed"
   "encoding/json"
   "errors"
   "io"
   "net/url"
   "path"
   "strings"
)

// ZGlzbmV5JmJyb3dzZXImMS4wLjA
// disney&browser&1.0.0
const client_api_key = "ZGlzbmV5JmJyb3dzZXImMS4wLjA.Cu56AgSfBTDag5NiRA81oLHkDZfu5L3CKadnefEAY84"

//go:embed authenticateWithOtp.gql
var mutation_authenticate_with_otp string

//go:embed loginWithActionGrant.gql
var mutation_login_with_action_grant string

//go:embed registerDevice.gql
var mutation_register_device string

//go:embed login.gql
var mutation_login string

//go:embed requestOtp.gql
var mutation_request_otp string

//go:embed refreshToken.gql
var mutation_refresh_token string

//go:embed switchProfile.gql
var mutation_switch_profile string

// https://disneyplus.com/browse/entity-7df81cf5-6be5-4e05-9ff6-da33baf0b94d
// https://disneyplus.com/cs-cz/browse/entity-7df81cf5-6be5-4e05-9ff6-da33baf0b94d
// https://disneyplus.com/play/7df81cf5-6be5-4e05-9ff6-da33baf0b94d
func GetEntityId(rawUrl string) (string, error) {
   parsed, err := url.Parse(rawUrl)
   if err != nil {
      return "", err
   }
   base := path.Base(parsed.Path)
   if !strings.HasPrefix(base, "entity-") {
      return "", errors.New("entity value missing from URL")
   }
   return base, nil
}

type AuthenticateWithOtp struct {
   ActionGrant string
}

type Error struct {
   Code        string // 2026-04-05
   Description string // 2026-04-05
}

type Login struct {
   Account struct {
      Profiles []Profile
   }
}

type LoginWithActionGrant struct {
   Account struct {
      Profiles []Profile
   }
}

type Page struct {
   Actions []struct {
      InternalTitle string // movie
   }
   Containers []struct {
      Seasons []struct { // series
         Visuals struct {
            Name string
         }
         Id string
      }
   }
   Visuals struct {
      Restriction struct {
         Message string
      }
   }
}

func (p *Profile) String() string {
   var data strings.Builder
   data.WriteString("name: ")
   data.WriteString(p.Name)
   data.WriteString("\nid: ")
   data.WriteString(p.Id)
   return data.String()
}

type Profile struct {
   Name string
   Id   string
}

type RequestOtp struct {
   Accepted bool
}

func (r *RequestOtp) String() string {
   if r.Accepted {
      return "accepted = true"
   }
   return "accepted = false"
}

func (s Season) String() string {
   var (
      data strings.Builder
      line bool
   )
   for _, item := range s.Items {
      for _, action := range item.Actions {
         if line {
            data.WriteByte('\n')
         } else {
            line = true
         }
         data.WriteString(action.InternalTitle)
      }
   }
   return data.String()
}

type Season struct {
   Items []struct {
      Actions []struct {
         InternalTitle string
      }
   }
}

func (t *Token) String() string {
   var data strings.Builder
   data.WriteString("type: ")
   data.WriteString(t.AccessTokenType)
   data.WriteString("\naccess token: ")
   data.WriteString(t.AccessToken)
   if t.RefreshToken != "" {
      data.WriteString("\nrefresh token: ")
      data.WriteString(t.RefreshToken)
   }
   return data.String()
}

// request: Account
func (t *Token) FetchSeason(id string) (*Season, error) {
   if err := t.assert("Account"); err != nil {
      return nil, err
   }
   resp, err := maya.Get(
      &url.URL{
         Scheme:   "https",
         Host:     "disney.api.edge.bamgrid.com",
         Path:     "/explore/v1.12/season/" + id,
         RawQuery: "limit=99",
      },
      map[string]string{"authorization": "Bearer " + t.AccessToken},
   )
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   var result struct {
      Data struct {
         Season Season
      }
   }
   err = json.NewDecoder(resp.Body).Decode(&result)
   if err != nil {
      return nil, err
   }
   return &result.Data.Season, nil
}

// request: Device
// response: AccountWithoutActiveProfile
func (t *Token) LoginWithActionGrant(actionGrant string) (*LoginWithActionGrant, error) {
   if err := t.assert("Device"); err != nil {
      return nil, err
   }
   body, err := json.Marshal(map[string]any{
      "query": mutation_login_with_action_grant,
      "variables": map[string]any{
         "input": map[string]string{
            "actionGrant": actionGrant,
         },
      },
   })
   if err != nil {
      return nil, err
   }
   resp, err := maya.Post(
      &url.URL{
         Scheme: "https",
         Host:   "disney.api.edge.bamgrid.com",
         Path:   "/v1/public/graphql",
      },
      map[string]string{"authorization": "Bearer " + t.AccessToken},
      body,
   )
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   var result struct {
      Data struct {
         LoginWithActionGrant LoginWithActionGrant
      }
      Extensions struct {
         Sdk struct {
            Token Token
         }
      }
   }
   err = json.NewDecoder(resp.Body).Decode(&result)
   if err != nil {
      return nil, err
   }
   *t = result.Extensions.Sdk.Token
   return &result.Data.LoginWithActionGrant, nil
}

// request: Device
// response: AccountWithoutActiveProfile
func (t *Token) FetchLogin(email, password string) (*Login, error) {
   if err := t.assert("Device"); err != nil {
      return nil, err
   }
   body, err := json.Marshal(map[string]any{
      "query": mutation_login,
      "variables": map[string]any{
         "input": map[string]string{
            "email":    email,
            "password": password,
         },
      },
   })
   if err != nil {
      return nil, err
   }
   resp, err := maya.Post(
      &url.URL{
         Scheme: "https",
         Host:   "disney.api.edge.bamgrid.com",
         Path:   "/v1/public/graphql",
      },
      map[string]string{"authorization": "Bearer " + t.AccessToken},
      body,
   )
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   var result struct {
      Data struct {
         Login Login
      }
      Extensions struct {
         Sdk struct {
            Token Token
         }
      }
   }
   err = json.NewDecoder(resp.Body).Decode(&result)
   if err != nil {
      return nil, err
   }
   *t = result.Extensions.Sdk.Token
   return &result.Data.Login, nil
}

// THIS REQUEST SETS THE LOCATION BASED ON YOUR IP
// request: AccountWithoutActiveProfile
// response: Account
func (t *Token) SwitchProfile(profileId string) error {
   if err := t.assert("AccountWithoutActiveProfile"); err != nil {
      return err
   }
   body, err := json.Marshal(map[string]any{
      "query": mutation_switch_profile,
      "variables": map[string]any{
         "input": map[string]string{
            "profileId": profileId,
         },
      },
   })
   if err != nil {
      return err
   }
   resp, err := maya.Post(
      &url.URL{
         Scheme: "https",
         Host:   "disney.api.edge.bamgrid.com",
         Path:   "/v1/public/graphql",
      },
      map[string]string{"authorization": "Bearer " + t.AccessToken},
      body,
   )
   if err != nil {
      return err
   }
   defer resp.Body.Close()
   var result struct {
      Extensions struct {
         Sdk struct {
            Token Token
         }
      }
   }
   err = json.NewDecoder(resp.Body).Decode(&result)
   if err != nil {
      return err
   }
   *t = result.Extensions.Sdk.Token
   return nil
}

// Response: Device
func RegisterDevice() (*Token, error) {
   body, err := json.Marshal(map[string]any{
      "query": mutation_register_device,
      "variables": map[string]any{
         "input": map[string]any{
            "deviceProfile":      "!",
            "deviceFamily":       "!",
            "applicationRuntime": "!",
            "attributes": map[string]string{
               "operatingSystem":        "",
               "operatingSystemVersion": "",
            },
         },
      },
   })
   if err != nil {
      return nil, err
   }
   resp, err := maya.Post(
      &url.URL{
         Scheme: "https",
         Host:   "disney.api.edge.bamgrid.com",
         Path:   "/graph/v1/device/graphql",
      },
      map[string]string{"authorization": "Bearer " + client_api_key},
      body,
   )
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   var result struct {
      Data struct {
         RegisterDevice struct {
            Token Token
         }
      }
   }
   err = json.NewDecoder(resp.Body).Decode(&result)
   if err != nil {
      return nil, err
   }
   return &result.Data.RegisterDevice.Token, nil
}

// expires: 4 hours
// request: Account
func (t *Token) Refresh() error {
   if err := t.assert("Account"); err != nil {
      return err
   }
   body, err := json.Marshal(map[string]any{
      "query": mutation_refresh_token,
      "variables": map[string]any{
         "input": map[string]string{
            "refreshToken": t.RefreshToken,
         },
      },
   })
   if err != nil {
      return err
   }
   resp, err := maya.Post(
      &url.URL{
         Scheme: "https",
         Host:   "disney.api.edge.bamgrid.com",
         Path:   "/graph/v1/device/graphql",
      },
      map[string]string{"authorization": "Bearer " + client_api_key},
      body,
   )
   if err != nil {
      return err
   }
   defer resp.Body.Close()
   var result struct {
      Extensions struct {
         Sdk struct {
            Token Token
         }
      }
   }
   err = json.NewDecoder(resp.Body).Decode(&result)
   if err != nil {
      return err
   }
   *t = result.Extensions.Sdk.Token
   return nil
}

// L3 max: 720p
// request: Account
func (t *Token) FetchWidevine(body []byte) ([]byte, error) {
   if err := t.assert("Account"); err != nil {
      return nil, err
   }
   resp, err := maya.Post(
      &url.URL{
         Scheme: "https",
         Host:   "disney.playback.edge.bamgrid.com",
         Path:   "/widevine/v1/obtain-license",
      },
      map[string]string{"authorization": t.AccessToken},
      body,
   )
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   return io.ReadAll(resp.Body)
}

// SL2000 max: 720p
// SL3000 max: 2160p
// request: Account
func (t *Token) FetchPlayReady(body []byte) ([]byte, error) {
   if err := t.assert("Account"); err != nil {
      return nil, err
   }
   resp, err := maya.Post(
      &url.URL{
         Scheme: "https",
         Host:   "disney.playback.edge.bamgrid.com",
         Path:   "/playready/v1/obtain-license.asmx",
      },
      map[string]string{"authorization": t.AccessToken},
      body,
   )
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   if resp.StatusCode != 200 {
      return nil, errors.New(resp.Status)
   }
   return io.ReadAll(resp.Body)
}

type Token struct {
   AccessTokenType string
   AccessToken     string
   RefreshToken    string
}

func (t *Token) assert(expected string) error {
   if t.AccessTokenType != expected {
      return errors.New("expected token type " + expected)
   }
   return nil
}

// request: Account
func (t *Token) FetchPage(entity string) (*Page, error) {
   if err := t.assert("Account"); err != nil {
      return nil, err
   }
   resp, err := maya.Get(
      &url.URL{
         Scheme:   "https",
         Host:     "disney.api.edge.bamgrid.com",
         Path:     "/explore/v1.12/page/" + entity,
         RawQuery: "limit=0",
      },
      map[string]string{"authorization": "Bearer " + t.AccessToken},
   )
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   var result struct {
      Data struct {
         Errors []Error // 2026-04-11
         Page   Page
      }
      Errors []Error // 2026-05-03
   }
   err = json.NewDecoder(resp.Body).Decode(&result)
   if err != nil {
      return nil, err
   }
   if len(result.Errors) >= 1 {
      return nil, &result.Errors[0]
   }
   if len(result.Data.Errors) >= 1 {
      return nil, &result.Data.Errors[0]
   }
   return &result.Data.Page, nil
}

// request: Device
func (t *Token) AuthenticateWithOtp(email, passcode string) (*AuthenticateWithOtp, error) {
   if err := t.assert("Device"); err != nil {
      return nil, err
   }
   body, err := json.Marshal(map[string]any{
      "query": mutation_authenticate_with_otp,
      "variables": map[string]any{
         "input": map[string]string{
            "email":    email,
            "passcode": passcode,
         },
      },
   })
   if err != nil {
      return nil, err
   }
   resp, err := maya.Post(
      &url.URL{
         Scheme: "https",
         Host:   "disney.api.edge.bamgrid.com",
         Path:   "/v1/public/graphql",
      },
      map[string]string{"authorization": "Bearer " + t.AccessToken},
      body,
   )
   if err != nil {
      return nil, err
   }
   var result struct {
      Data struct {
         AuthenticateWithOtp AuthenticateWithOtp
      }
   }
   err = json.NewDecoder(resp.Body).Decode(&result)
   if err != nil {
      return nil, err
   }
   return &result.Data.AuthenticateWithOtp, nil
}

///

// request: Account
func (t *Token) FetchStream(mediaId string) (*url.URL, error) {
   if err := t.assert("Account"); err != nil {
      return nil, err
   }
   playback_id, err := json.Marshal(map[string]string{
      "mediaId": mediaId,
   })
   if err != nil {
      return nil, err
   }
   body, err := json.Marshal(map[string]any{
      "playback": map[string]any{
         "attributes": map[string]any{
            "assetInsertionStrategy": "SGAI",
            "codecs": map[string]any{
               "supportsMultiCodecMaster": true, // 4K
               "video": []string{
                  "h.264",
                  "h.265",
               },
            },
            "videoRanges": []string{"HDR10"},
         },
      },
      "playbackId": playback_id,
   })
   if err != nil {
      return nil, err
   }
   resp, err := maya.Post(
      &url.URL{
         Scheme: "https",
         Host:   "disney.playback.edge.bamgrid.com",
         // /v7/playback/ctr-high
         // /v7/playback/tv-drm-ctr-h265-atmos
         Path: "/v7/playback/ctr-regular",
      },
      map[string]string{
         "authorization":           "Bearer " + t.AccessToken,
         "content-type":            "application/json",
         "x-application-version":   "",
         "x-bamsdk-client-id":      "",
         "x-bamsdk-platform":       "",
         "x-bamsdk-version":        "",
         "x-dss-feature-filtering": "true",
      },
      body,
   )
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   var result struct {
      Stream struct {
         Sources []struct {
            Complete struct {
               Url string
            }
         }
      }
   }
   err = json.NewDecoder(resp.Body).Decode(&result)
   if err != nil {
      return nil, err
   }
   return url.Parse(result.Stream.Sources[0].Complete.Url)
}

// request: Device
func (t *Token) RequestOtp(email string) (*RequestOtp, error) {
   if err := t.assert("Device"); err != nil {
      return nil, err
   }
   body, err := json.Marshal(map[string]any{
      "query": mutation_request_otp,
      "variables": map[string]any{
         "input": map[string]string{
            "email":  email,
            "reason": "Login",
         },
      },
   })
   if err != nil {
      return nil, err
   }
   resp, err := maya.Post(
      &url.URL{
         Scheme: "https",
         Host:   "disney.api.edge.bamgrid.com",
         Path:   "/v1/public/graphql",
      },
      map[string]string{"authorization": "Bearer " + t.AccessToken},
      body,
   )
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   var result struct {
      Data struct {
         RequestOtp RequestOtp
      }
   }
   err = json.NewDecoder(resp.Body).Decode(&result)
   if err != nil {
      return nil, err
   }
   return &result.Data.RequestOtp, nil
}

func (p *Page) String() string {
   var data strings.Builder
   if len(p.Containers[0].Seasons) >= 1 {
      var line bool
      for _, seasonItem := range p.Containers[0].Seasons {
         if line {
            data.WriteString("\n\n")
         } else {
            line = true
         }
         data.WriteString("name: ")
         data.WriteString(seasonItem.Visuals.Name)
         data.WriteString("\nid: ")
         data.WriteString(seasonItem.Id)
      }
   } else {
      data.WriteString(p.Actions[0].InternalTitle)
   }
   return data.String()
}

func (e *Error) Error() string {
   var data strings.Builder
   data.WriteString("code: ")
   data.WriteString(e.Code)
   data.WriteString("\ndescription: ")
   data.WriteString(e.Description)
   return data.String()
}
