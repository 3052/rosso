package draken

import (
   "bytes"
   _ "embed"
   "encoding/json"
   "errors"
   "io"
   "net/http"
   "net/url"
   "strings"
)

func ParseDash(playlist string) (*url.URL, error) {
   return url.Parse(playlist)
}

type Playback struct {
   Headers struct {
      MaginePlayEntitlementId string `json:"Magine-Play-EntitlementId"`
      MaginePlaySession       string `json:"Magine-Play-Session"`
   }
   Playlist string // MPD
}

func FetchPlayback(loginToken, playableId, entitlementId string) (*Playback, error) {
   req := http.Request{
      Method: "POST",
      URL: &url.URL{
         Scheme: "https",
         Host:   "client-api.magine.com",
         Path:   "/api/playback/v1/preflight/asset/" + playableId,
      },
      Header: http.Header{},
   }
   setBaseHeaders(&req, loginToken)
   setPlaybackHeaders(&req)
   req.Header.Set("magine-play-entitlementid", entitlementId)
   // this value is important, with the wrong value you get random failures
   req.Header.Set("x-forwarded-for", "95.192.0.0")
   resp, err := http.DefaultClient.Do(&req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   result := &Playback{}
   err = json.NewDecoder(resp.Body).Decode(result)
   if err != nil {
      return nil, err
   }
   return result, nil
}

//go:embed GetCustomIdFullMovie.gql
var get_custom_id_full_movie string

// setPlaybackHeaders adds the headers specific to playback functionality.
func setPlaybackHeaders(req *http.Request) {
   req.Header.Set("magine-play-deviceid", "!")
   req.Header.Set("magine-play-devicemodel", "firefox 111.0 / windows 10")
   req.Header.Set("magine-play-deviceplatform", "firefox")
   req.Header.Set("magine-play-devicetype", "web")
   req.Header.Set("magine-play-drm", "widevine")
   req.Header.Set("magine-play-protocol", "dashs")
}

// setBaseHeaders adds the common authentication and access tokens to a request.
func setBaseHeaders(req *http.Request, loginToken string) {
   req.Header.Set("magine-accesstoken", "22cc71a2-8b77-4819-95b0-8c90f4cf5663")
   if loginToken != "" {
      req.Header.Set("authorization", "Bearer "+loginToken)
   }
}

type Entitlement struct {
   Error *Error
   Token string
}

func (e *Error) Error() string {
   var data strings.Builder
   data.WriteString("message = ")
   data.WriteString(e.Message)
   data.WriteString("\nuser message = ")
   data.WriteString(e.UserMessage)
   return data.String()
}

type Error struct {
   Message     string
   UserMessage string `json:"user_message"`
}

func FetchLogin(identity, accessKey string) (*Login, error) {
   body, err := json.Marshal(map[string]string{
      "accessKey": accessKey,
      "identity":  identity,
   })
   if err != nil {
      return nil, err
   }
   req, err := http.NewRequest(
      "POST", "https://client-api.magine.com/api/login/v2/auth/email",
      bytes.NewReader(body),
   )
   if err != nil {
      return nil, err
   }
   setBaseHeaders(req, "") // No login token for this request
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   var result Login
   err = json.NewDecoder(resp.Body).Decode(&result)
   if err != nil {
      return nil, err
   }
   if result.Message != "" {
      return nil, errors.New(result.Message)
   }
   return &result, nil
}

type Login struct {
   Message string
   Token   string
}

func FetchPlayableId(customId string) (string, error) {
   body, err := json.Marshal(map[string]any{
      "query":     get_custom_id_full_movie,
      "variables": map[string]string{"customId": customId},
   })
   if err != nil {
      return "", err
   }
   req, err := http.NewRequest(
      "POST", "https://client-api.magine.com/api/apiql/v2",
      bytes.NewReader(body),
   )
   if err != nil {
      return "", err
   }
   setBaseHeaders(req, "") // No login token for this request
   // this value is important, with the wrong value you get random failures
   req.Header.Set("x-forwarded-for", "95.192.0.0")
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return "", err
   }
   defer resp.Body.Close()
   var result struct {
      Data struct {
         Viewer struct {
            ViewableCustomId *struct {
               DefaultPlayable struct {
                  Id string
               }
            }
         }
      }
   }
   err = json.NewDecoder(resp.Body).Decode(&result)
   if err != nil {
      return "", err
   }
   if result.Data.Viewer.ViewableCustomId == nil {
      return "", errors.New("ViewableCustomId")
   }
   return result.Data.Viewer.ViewableCustomId.DefaultPlayable.Id, nil
}

func FetchEntitlement(loginToken, playableId string) (*Entitlement, error) {
   req := http.Request{
      Method: "POST",
      URL: &url.URL{
         Scheme: "https",
         Host:   "client-api.magine.com",
         Path:   "/api/entitlement/v2/asset/" + playableId,
      },
      Header: http.Header{},
   }
   setBaseHeaders(&req, loginToken)
   resp, err := http.DefaultClient.Do(&req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   var result Entitlement
   err = json.NewDecoder(resp.Body).Decode(&result)
   if err != nil {
      return nil, err
   }
   if result.Error != nil {
      return nil, result.Error
   }
   return &result, nil
}

func (p *Playback) FetchWidevine(loginToken string, body []byte) ([]byte, error) {
   req, err := http.NewRequest(
      "POST", "https://client-api.magine.com/api/playback/v1/widevine/license",
      bytes.NewReader(body),
   )
   if err != nil {
      return nil, err
   }
   setBaseHeaders(req, loginToken)
   setPlaybackHeaders(req)
   req.Header.Set("magine-play-session", p.Headers.MaginePlaySession)
   req.Header.Set("magine-play-entitlementId", p.Headers.MaginePlayEntitlementId)
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   return io.ReadAll(resp.Body)
}
