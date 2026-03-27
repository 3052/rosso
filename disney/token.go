package disney

import (
   "bytes"
   "encoding/json"
   "errors"
   "io"
   "net/http"
   "net/url"
)

// request: Account
func (t *Token) Season(id string) (*Season, error) {
   if err := t.assert("Account"); err != nil {
      return nil, err
   }
   req := http.Request{
      URL: &url.URL{
         Scheme:   "https",
         Host:     "disney.api.edge.bamgrid.com",
         Path:     "/explore/v1.12/season/" + id,
         RawQuery: "limit=99",
      },
      Header: http.Header{},
   }
   req.Header.Set("authorization", "Bearer "+t.AccessToken)
   resp, err := http.DefaultClient.Do(&req)
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
func (t *Token) Page(entity string) (*Page, error) {
   if err := t.assert("Account"); err != nil {
      return nil, err
   }
   var req http.Request
   req.Header = http.Header{}
   req.Header.Set("authorization", "Bearer "+t.AccessToken)
   req.URL = &url.URL{
      Scheme:   "https",
      Host:     "disney.api.edge.bamgrid.com",
      Path:     "/explore/v1.12/page/entity-" + entity,
      RawQuery: "limit=0",
   }
   resp, err := http.DefaultClient.Do(&req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   var result struct {
      Data struct {
         Errors []Error // region
         Page   Page
      }
      Errors []Error // auth.expired
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

// request: Account
func (t *Token) Stream(mediaId string) (*Stream, error) {
   if err := t.assert("Account"); err != nil {
      return nil, err
   }
   playback_id, err := json.Marshal(map[string]string{
      "mediaId": mediaId,
   })
   if err != nil {
      return nil, err
   }
   data, err := json.Marshal(map[string]any{
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
   var req http.Request
   req.Method = "POST"
   req.URL = &url.URL{
      Scheme: "https",
      Host:   "disney.playback.edge.bamgrid.com",
      // Path: "/v7/playback/ctr-high",
      // Path: "/v7/playback/tv-drm-ctr-h265-atmos",
      Path: "/v7/playback/ctr-regular",
   }
   req.Header = http.Header{}
   req.Header.Set("authorization", "Bearer "+t.AccessToken)
   req.Header.Set("content-type", "application/json")
   req.Header.Set("x-application-version", "")
   req.Header.Set("x-bamsdk-client-id", "")
   req.Header.Set("x-bamsdk-platform", "")
   req.Header.Set("x-bamsdk-version", "")
   req.Header.Set("x-dss-feature-filtering", "true")
   req.Body = io.NopCloser(bytes.NewReader(data))
   resp, err := http.DefaultClient.Do(&req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   var result struct {
      Errors []Error
      Stream Stream
   }
   err = json.NewDecoder(resp.Body).Decode(&result)
   if err != nil {
      return nil, err
   }
   if len(result.Errors) >= 1 {
      return nil, &result.Errors[0]
   }
   return &result.Stream, nil
}

// THIS REQUEST SETS THE LOCATION BASED ON YOUR IP
// request: AccountWithoutActiveProfile
// response: Account
func (t *Token) SwitchProfile(profileId string) error {
   if err := t.assert("AccountWithoutActiveProfile"); err != nil {
      return err
   }
   data, err := json.Marshal(map[string]any{
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
      "POST", "https://disney.api.edge.bamgrid.com/v1/public/graphql",
      bytes.NewReader(data),
   )
   if err != nil {
      return err
   }
   req.Header.Set("authorization", "Bearer "+t.AccessToken)
   resp, err := http.DefaultClient.Do(req)
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
// expires: 4 hours
// request: Account
func RefreshToken(refresh *Token) error {
   if err := refresh.assert("Account"); err != nil {
      return err
   }
   data, err := json.Marshal(map[string]any{
      "query": mutation_refresh_token,
      "variables": map[string]any{
         "input": map[string]string{
            "refreshToken": refresh.RefreshToken,
         },
      },
   })
   if err != nil {
      return err
   }
   req, err := http.NewRequest(
      "POST", "https://disney.api.edge.bamgrid.com/graph/v1/device/graphql",
      bytes.NewReader(data),
   )
   if err != nil {
      return err
   }
   req.Header.Set("authorization", "Bearer "+client_api_key)
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return err
   }
   defer resp.Body.Close()
   return json.NewDecoder(resp.Body).Decode(refresh)
}

func (r *RequestOtp) String() string {
   if r.Accepted {
      return "accepted = true"
   }
   return "accepted = false"
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

// request: Device
// response: AccountWithoutActiveProfile
func (t *Token) Login(email, password string) (*Login, error) {
   if err := t.assert("Device"); err != nil {
      return nil, err
   }
   data, err := json.Marshal(map[string]any{
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
      "POST", "https://disney.api.edge.bamgrid.com/v1/public/graphql",
      bytes.NewReader(data),
   )
   if err != nil {
      return nil, err
   }
   req.Header.Set("authorization", "Bearer "+t.AccessToken)
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   var result struct {
      Data struct {
         Login Login
      }
      Errors     []Error
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
   if len(result.Errors) >= 1 {
      return nil, &result.Errors[0]
   }
   *t = result.Extensions.Sdk.Token
   return &result.Data.Login, nil
}

// request: Device
func (t *Token) RequestOtp(email string) (*RequestOtp, error) {
   if err := t.assert("Device"); err != nil {
      return nil, err
   }
   data, err := json.Marshal(map[string]any{
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
      "POST", "https://disney.api.edge.bamgrid.com/v1/public/graphql",
      bytes.NewReader(data),
   )
   if err != nil {
      return nil, err
   }
   req.Header.Set("authorization", "Bearer "+t.AccessToken)
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   var result struct {
      Data struct {
         RequestOtp RequestOtp
      }
      Errors []Error
   }
   err = json.NewDecoder(resp.Body).Decode(&result)
   if err != nil {
      return nil, err
   }
   if len(result.Errors) >= 1 {
      return nil, &result.Errors[0]
   }
   return &result.Data.RequestOtp, nil
}

// request: Device
func (t *Token) AuthenticateWithOtp(email, passcode string) (*AuthenticateWithOtp, error) {
   if err := t.assert("Device"); err != nil {
      return nil, err
   }
   data, err := json.Marshal(map[string]any{
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
      "POST", "https://disney.api.edge.bamgrid.com/v1/public/graphql",
      bytes.NewReader(data),
   )
   if err != nil {
      return nil, err
   }
   req.Header.Set("authorization", "Bearer "+t.AccessToken)
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   var result struct {
      Data struct {
         AuthenticateWithOtp AuthenticateWithOtp
      }
      Errors []Error
   }
   err = json.NewDecoder(resp.Body).Decode(&result)
   if err != nil {
      return nil, err
   }
   if len(result.Errors) >= 1 {
      return nil, &result.Errors[0]
   }
   return &result.Data.AuthenticateWithOtp, nil
}

// request: Device
// response: AccountWithoutActiveProfile
func (t *Token) LoginWithActionGrant(actionGrant string) (*LoginWithActionGrant, error) {
   if err := t.assert("Device"); err != nil {
      return nil, err
   }
   data, err := json.Marshal(map[string]any{
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
      "POST", "https://disney.api.edge.bamgrid.com/v1/public/graphql",
      bytes.NewReader(data),
   )
   if err != nil {
      return nil, err
   }
   req.Header.Set("authorization", "Bearer "+t.AccessToken)
   resp, err := http.DefaultClient.Do(req)
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

// SL2000: 720p
// SL3000: 2160p
// request: Account
func (t *Token) PlayReady(data []byte) ([]byte, error) {
   if err := t.assert("Account"); err != nil {
      return nil, err
   }
   req, err := http.NewRequest(
      "POST",
      "https://disney.playback.edge.bamgrid.com/playready/v1/obtain-license.asmx",
      bytes.NewReader(data),
   )
   if err != nil {
      return nil, err
   }
   req.Header.Set("authorization", t.AccessToken)
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   if resp.StatusCode != http.StatusOK {
      return nil, errors.New(resp.Status)
   }
   return io.ReadAll(resp.Body)
}

// L1: 2160p
// L3: 720p
// request: Account
func (t *Token) Widevine(data []byte) ([]byte, error) {
   if err := t.assert("Account"); err != nil {
      return nil, err
   }
   req, err := http.NewRequest(
      "POST",
      "https://disney.playback.edge.bamgrid.com/widevine/v1/obtain-license",
      bytes.NewReader(data),
   )
   if err != nil {
      return nil, err
   }
   req.Header.Set("authorization", t.AccessToken)
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   return io.ReadAll(resp.Body)
}
