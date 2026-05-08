package amc

import (
   "41.neocities.org/maya"
   _ "embed"
   "encoding/json"
   "fmt"
   "io"
   "net/url"
)

type Subheading struct {
   ID    string `json:"id,omitempty"`
   Title string `json:"title,omitempty"`
   Type  string `json:"type,omitempty"`
}

type Text struct {
   Title       *TextElement `json:"title,omitempty"`
   Description *TextElement `json:"description,omitempty"`
   Subheadings []Subheading `json:"subheadings,omitempty"`
}

type TextElement struct {
   Title string `json:"title,omitempty"`
}

type TTS struct {
   SpeechText string `json:"speechText,omitempty"`
}

func (s *Source) GetManifest() (*url.URL, error) {
   return url.Parse(s.Src)
}

type Source struct {
   Codecs     string
   KeySystems KeySystems `json:"key_systems"`
   Src        string     // MPD
   Type       string
}

func (a *AuthData) Refresh() error {
   resp, err := maya.Post(
      &url.URL{
         Scheme: "https",
         Host:   "gw.cds.amcn.com",
         Path:   "/auth-orchestration-id/api/v1/refresh",
      },
      map[string]string{"authorization": "Bearer " + a.RefreshToken},
      nil,
   )
   if err != nil {
      return err
   }
   defer resp.Body.Close()
   if resp.StatusCode != 200 {
      return fmt.Errorf("refresh failed with status: %d", resp.StatusCode)
   }
   var result struct {
      Success bool     `json:"success"`
      Status  int      `json:"status"`
      Data    AuthData `json:"data"`
   }
   if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
      return err
   }
   *a = result.Data
   return nil
}

//go:embed playback.json
var playback_json []byte

func License(licenseUrl, bcovAuth string, challenge []byte) ([]byte, error) {
   target, err := url.Parse(licenseUrl)
   if err != nil {
      return nil, err
   }
   resp, err := maya.Post(
      target, map[string]string{"bcov-auth": bcovAuth}, challenge,
   )
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   if resp.StatusCode != 200 {
      return nil, fmt.Errorf("license request failed with status: %d", resp.StatusCode)
   }
   return io.ReadAll(resp.Body)
}

// AuthData represents the inner payload of authentication responses.
type AuthData struct {
   AccessToken  string `json:"access_token"`
   RefreshToken string `json:"refresh_token"`
   TokenType    string `json:"token_type"`
   ExpiresIn    int    `json:"expires_in"`
}

// Login authenticates the user. It requires the guest token (access_token)
// retrieved from calling the Unauth() function.
func Login(guestToken, email, password string) (*AuthData, error) {
   // Body
   body, err := json.Marshal(map[string]string{
      "email":    email,
      "password": password,
   })
   if err != nil {
      return nil, err
   }
   resp, err := maya.Post(
      &url.URL{
         Scheme: "https",
         Host:   "gw.cds.amcn.com",
         Path:   "/auth-orchestration-id/api/v1/login",
      },
      map[string]string{
         "authorization":           "Bearer " + guestToken,
         "content-type":            "application/json",
         "x-amcn-language":         "en",
         "x-amcn-network":          "amcplus",
         "x-amcn-platform":         "web",
         "x-amcn-service-group-id": "10",
         "x-amcn-tenant":           "amcn",
         "x-amcn-device-ad-id":     "-",
         "x-amcn-device-id":        "-",
         "x-amcn-service-id":       "amcplus",
         "x-ccpa-do-not-sell":      "doNotPassData",
      },
      body,
   )
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   if resp.StatusCode != 200 {
      return nil, fmt.Errorf("login failed with status: %d", resp.StatusCode)
   }
   // Internal envelope to strip the first layer
   var envelope struct {
      Success bool     `json:"success"`
      Status  int      `json:"status"`
      Data    AuthData `json:"data"`
   }
   if err := json.NewDecoder(resp.Body).Decode(&envelope); err != nil {
      return nil, err
   }
   return &envelope.Data, nil
}

func Unauth() (*AuthData, error) {
   resp, err := maya.Post(
      &url.URL{
         Scheme: "https",
         Host:   "gw.cds.amcn.com",
         Path:   "/auth-orchestration-id/api/v1/unauth",
      },
      map[string]string{
         "x-amcn-network":   "amcplus",
         "x-amcn-platform":  "web",
         "x-amcn-tenant":    "amcn",
         "x-amcn-device-id": "-",
         "x-amcn-language":  "en",
      },
      nil,
   )
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   if resp.StatusCode != 200 {
      return nil, fmt.Errorf("unauth failed with status: %d", resp.StatusCode)
   }
   // Internal envelope to strip the first layer
   var envelope struct {
      Success bool     `json:"success"`
      Status  int      `json:"status"`
      Data    AuthData `json:"data"`
   }
   if err := json.NewDecoder(resp.Body).Decode(&envelope); err != nil {
      return nil, err
   }
   return &envelope.Data, nil
}

func SeriesDetail(authToken string, seriesId int) (*ContentNode, error) {
   resp, err := maya.Get(
      &url.URL{
         Scheme: "https",
         Host:   "gw.cds.amcn.com",
         Path: fmt.Sprint(
            "/content-compiler-cr/api/v1/content/amcn/amcplus/type/series-detail/id/",
            seriesId,
         ),
      },
      map[string]string{
         "authorization":   "Bearer " + authToken,
         "x-amcn-network":  "amcplus",
         "x-amcn-platform": "android",
         "x-amcn-tenant":   "amcn",
      },
   )
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   if resp.StatusCode != 200 {
      return nil, fmt.Errorf("series detail failed with status: %d", resp.StatusCode)
   }
   // Internal envelope to strip the first layer
   var envelope struct {
      Success bool        `json:"success"`
      Status  int         `json:"status"`
      Data    ContentNode `json:"data"`
   }
   if err := json.NewDecoder(resp.Body).Decode(&envelope); err != nil {
      return nil, err
   }
   return &envelope.Data, nil
}

func SeasonEpisodes(authToken string, seasonId int) (*ContentNode, error) {
   resp, err := maya.Get(
      &url.URL{
         Scheme: "https",
         Host:   "gw.cds.amcn.com",
         Path: fmt.Sprint(
            "/content-compiler-cr/api/v1/content/amcn/amcplus/type/season-episodes/id/",
            seasonId,
         ),
      },
      map[string]string{
         "authorization":   "Bearer " + authToken,
         "x-amcn-network":  "amcplus",
         "x-amcn-platform": "android",
         "x-amcn-tenant":   "amcn",
      },
   )
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   if resp.StatusCode != 200 {
      return nil, fmt.Errorf("season episodes failed with status: %d", resp.StatusCode)
   }
   // Internal envelope to strip the first layer
   var envelope struct {
      Success bool        `json:"success"`
      Status  int         `json:"status"`
      Data    ContentNode `json:"data"`
   }
   if err := json.NewDecoder(resp.Body).Decode(&envelope); err != nil {
      return nil, err
   }
   return &envelope.Data, nil
}

type DownloadData struct {
   Downloadable        bool `json:"downloadable,omitempty"`
   DownloadingExpireIn int  `json:"downloadingExpireIn,omitempty"`
   DownloadingEndDate  int  `json:"downloadingEndDate,omitempty"`
}

type Images struct {
   Default string `json:"default,omitempty"`
   Mobile  string `json:"mobile,omitempty"`
   Tablet  string `json:"tablet,omitempty"`
}

type KeySystems struct {
   ComWidevineAlpha struct {
      LicenseURL string `json:"license_url"`
   } `json:"com.widevine.alpha"`
   ComMicrosoftPlayready struct {
      LicenseURL string `json:"license_url"`
   } `json:"com.microsoft.playready"`
}

type Navigation struct {
   ClientRequest struct {
      Endpoint string `json:"endpoint,omitempty"`
   } `json:"client_request,omitempty"`
   ContentID    string `json:"content_id,omitempty"`
   ContentType  string `json:"contentType,omitempty"`
   MicroAppType string `json:"micro_app_type,omitempty"`
   Properties   struct {
      Fullscreen bool   `json:"fullscreen,omitempty"`
      IsLive     bool   `json:"isLive,omitempty"`
      VideoTitle string `json:"videoTitle,omitempty"`
   } `json:"properties,omitempty"`
   ScreenDesignType string `json:"screenDesignType,omitempty"`
}

func GetPlayback(authToken string, videoId int) (*Playback, error) {
   resp, err := maya.Post(
      &url.URL{
         Scheme: "https",
         Host:   "gw.cds.amcn.com",
         Path:   fmt.Sprint("/playback-id/api/v1/playback/", videoId),
      },
      map[string]string{
         "authorization":       "Bearer " + authToken,
         "content-type":        "application/json",
         "x-amcn-language":     "en",
         "x-amcn-network":      "amcplus",
         "x-amcn-platform":     "web",
         "x-amcn-service-id":   "amcplus",
         "x-amcn-tenant":       "amcn",
         "x-amcn-device-ad-id": "-",
         "x-ccpa-do-not-sell":  "doNotPassData",
      },
      playback_json,
   )
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   if resp.StatusCode != 200 {
      return nil, fmt.Errorf("playback failed with status: %d", resp.StatusCode)
   }
   var result struct {
      Data struct {
         PlaybackJsonData struct {
            Sources []Source
         }
      }
   }
   if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
      return nil, err
   }
   return &Playback{
      BcovAuth: resp.Header.Get("x-amcn-bc-jwt"),
      Sources:  result.Data.PlaybackJsonData.Sources,
   }, nil
}

type Playback struct {
   BcovAuth string
   Sources  []Source
}

// DashSource finds and returns the first Source with the type
// "application/dash+xml"
func (p *Playback) Dash() (*Source, error) {
   for _, source_data := range p.Sources {
      if source_data.Type == "application/dash+xml" {
         return &source_data, nil
      }
   }
   return nil, fmt.Errorf("application/dash+xml source not found")
}
