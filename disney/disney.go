package disney

import (
   "bytes"
   _ "embed"
   "encoding/json"
   "errors"
   "io"
   "net/http"
   "net/url"
   "path"
   "strings"
)

// token.go
// ZGlzbmV5JmJyb3dzZXImMS4wLjA
// disney&browser&1.0.0
const client_api_key = "ZGlzbmV5JmJyb3dzZXImMS4wLjA.Cu56AgSfBTDag5NiRA81oLHkDZfu5L3CKadnefEAY84"

//go:embed authenticateWithOtp.gql
var mutation_authenticate_with_otp string

//go:embed login.gql
var mutation_login string

//go:embed loginWithActionGrant.gql
var mutation_login_with_action_grant string

//go:embed refreshToken.gql
var mutation_refresh_token string

//go:embed registerDevice.gql
var mutation_register_device string

//go:embed requestOtp.gql
var mutation_request_otp string

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
   req, err := http.NewRequest(
      "POST",
      "https://disney.api.edge.bamgrid.com/graph/v1/device/graphql",
      bytes.NewReader(body),
   )
   if err != nil {
      return nil, err
   }
   req.Header.Set("authorization", "Bearer "+client_api_key)
   resp, err := do(req)
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

type AuthenticateWithOtp struct {
   ActionGrant string
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

type Profile struct {
   Name string
   Id   string
}

func (p *Profile) String() string {
   var data strings.Builder
   data.WriteString("name: ")
   data.WriteString(p.Name)
   data.WriteString("\nid: ")
   data.WriteString(p.Id)
   return data.String()
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

type Season struct {
   Items []struct {
      Actions []struct {
         InternalTitle string
      }
   }
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
   req, err := http.NewRequest(
      "POST",
      "https://disney.api.edge.bamgrid.com/v1/public/graphql",
      bytes.NewReader(body),
   )
   if err != nil {
      return nil, err
   }
   req.Header.Set("authorization", "Bearer "+t.AccessToken)
   resp, err := do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
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

func (*Token) CachePath() string {
   return "rosso/disney/Token"
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
   req, err := http.NewRequest(
      "POST",
      "https://disney.api.edge.bamgrid.com/v1/public/graphql",
      bytes.NewReader(body),
   )
   if err != nil {
      return nil, err
   }
   req.Header.Set("authorization", "Bearer "+t.AccessToken)
   resp, err := do(req)
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

// SL2000 max: 720p
// SL3000 max: 2160p
// request: Account
func (t *Token) FetchPlayReady(body []byte) ([]byte, error) {
   if err := t.assert("Account"); err != nil {
      return nil, err
   }
   req, err := http.NewRequest(
      "POST",
      "https://disney.playback.edge.bamgrid.com/playready/v1/obtain-license.asmx",
      bytes.NewReader(body),
   )
   if err != nil {
      return nil, err
   }
   req.Header.Set("authorization", t.AccessToken)
   resp, err := do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   if resp.StatusCode != 200 {
      return nil, errors.New(resp.Status)
   }
   return io.ReadAll(resp.Body)
}

// request: Account
func (t *Token) FetchSeason(id string) (*Season, error) {
   if err := t.assert("Account"); err != nil {
      return nil, err
   }
   req, err := http.NewRequest(
      "GET",
      "https://disney.api.edge.bamgrid.com/explore/v1.12/season/"+id+"?limit=99",
      nil,
   )
   if err != nil {
      return nil, err
   }
   req.Header.Set("authorization", "Bearer "+t.AccessToken)
   resp, err := do(req)
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
   req, err := http.NewRequest(
      "POST",
      "https://disney.playback.edge.bamgrid.com/v7/playback/ctr-regular",
      bytes.NewReader(body),
   )
   if err != nil {
      return nil, err
   }
   req.Header.Set("authorization", "Bearer "+t.AccessToken)
   req.Header.Set("content-type", "application/json")
   req.Header.Set("x-application-version", "")
   req.Header.Set("x-bamsdk-client-id", "")
   req.Header.Set("x-bamsdk-platform", "")
   req.Header.Set("x-bamsdk-version", "")
   req.Header.Set("x-dss-feature-filtering", "true")
   resp, err := do(req)
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

// L3 max: 720p
// request: Account
func (t *Token) FetchWidevine(body []byte) ([]byte, error) {
   if err := t.assert("Account"); err != nil {
      return nil, err
   }
   req, err := http.NewRequest(
      "POST",
      "https://disney.playback.edge.bamgrid.com/widevine/v1/obtain-license",
      bytes.NewReader(body),
   )
   if err != nil {
      return nil, err
   }
   req.Header.Set("authorization", t.AccessToken)
   resp, err := do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   return io.ReadAll(resp.Body)
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
   req, err := http.NewRequest(
      "POST",
      "https://disney.api.edge.bamgrid.com/v1/public/graphql",
      bytes.NewReader(body),
   )
   if err != nil {
      return nil, err
   }
   req.Header.Set("authorization", "Bearer "+t.AccessToken)
   resp, err := do(req)
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
   req, err := http.NewRequest(
      "POST",
      "https://disney.api.edge.bamgrid.com/graph/v1/device/graphql",
      bytes.NewReader(body),
   )
   if err != nil {
      return err
   }
   req.Header.Set("authorization", "Bearer "+client_api_key)
   resp, err := do(req)
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
   req, err := http.NewRequest(
      "POST",
      "https://disney.api.edge.bamgrid.com/v1/public/graphql",
      bytes.NewReader(body),
   )
   if err != nil {
      return nil, err
   }
   req.Header.Set("authorization", "Bearer "+t.AccessToken)
   resp, err := do(req)
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
   req, err := http.NewRequest(
      "POST",
      "https://disney.api.edge.bamgrid.com/v1/public/graphql",
      bytes.NewReader(body),
   )
   if err != nil {
      return err
   }
   req.Header.Set("authorization", "Bearer "+t.AccessToken)
   resp, err := do(req)
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

// disney.go
