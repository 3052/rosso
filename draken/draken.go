package draken

import (
   _ "embed"
   "encoding/json"
   "errors"
   "io"
   "net/url"
   "strings"

   "41.neocities.org/maya"
)

func (p *Playback) GetManifest() (*url.URL, error) {
   return url.Parse(p.Playlist)
}

type Playback struct {
   Headers struct {
      MaginePlayEntitlementId string `json:"Magine-Play-EntitlementId"`
      MaginePlaySession       string `json:"Magine-Play-Session"`
   }
   Playlist string // MPD
}

func FetchPlayback(loginToken, playableId, entitlementId string) (*Playback, error) {
   headers := map[string]string{
      "magine-play-entitlementid": entitlementId,
      // this value is important, with the wrong value you get random failures
      "x-forwarded-for": "95.192.0.0",
   }
   setBaseHeaders(headers, loginToken)
   setPlaybackHeaders(headers)

   resp, err := maya.Post(&url.URL{
      Scheme: "https",
      Host:   "client-api.magine.com",
      Path:   "/api/playback/v1/preflight/asset/" + playableId,
   }, headers, nil)
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
func setPlaybackHeaders(headers map[string]string) {
   headers["magine-play-deviceid"] = "!"
   headers["magine-play-devicemodel"] = "firefox 111.0 / windows 10"
   headers["magine-play-deviceplatform"] = "firefox"
   headers["magine-play-devicetype"] = "web"
   headers["magine-play-drm"] = "widevine"
   headers["magine-play-protocol"] = "dashs"
}

// setBaseHeaders adds the common authentication and access tokens to a request.
func setBaseHeaders(headers map[string]string, loginToken string) {
   headers["magine-accesstoken"] = "22cc71a2-8b77-4819-95b0-8c90f4cf5663"
   if loginToken != "" {
      headers["authorization"] = "Bearer " + loginToken
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
   headers := make(map[string]string)
   setBaseHeaders(headers, "") // No login token for this request

   resp, err := maya.Post(&url.URL{
      Scheme: "https",
      Host:   "client-api.magine.com",
      Path:   "/api/login/v2/auth/email",
   }, headers, body)
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
   headers := map[string]string{
      // this value is important, with the wrong value you get random failures
      "x-forwarded-for": "95.192.0.0",
   }
   setBaseHeaders(headers, "") // No login token for this request

   resp, err := maya.Post(&url.URL{
      Scheme: "https",
      Host:   "client-api.magine.com",
      Path:   "/api/apiql/v2",
   }, headers, body)
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
   headers := make(map[string]string)
   setBaseHeaders(headers, loginToken)

   resp, err := maya.Post(&url.URL{
      Scheme: "https",
      Host:   "client-api.magine.com",
      Path:   "/api/entitlement/v2/asset/" + playableId,
   }, headers, nil)
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
   headers := map[string]string{
      "magine-play-session":       p.Headers.MaginePlaySession,
      "magine-play-entitlementId": p.Headers.MaginePlayEntitlementId,
   }
   setBaseHeaders(headers, loginToken)
   setPlaybackHeaders(headers)

   resp, err := maya.Post(&url.URL{
      Scheme: "https",
      Host:   "client-api.magine.com",
      Path:   "/api/playback/v1/widevine/license",
   }, headers, body)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   return io.ReadAll(resp.Body)
}
