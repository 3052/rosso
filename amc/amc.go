package amc

import (
   "bytes"
   "encoding/json"
   "fmt"
   "io"
   "net/http"
   "net/url"
)

func (s *Source) ParseDash() (*url.URL, error) {
   return url.Parse(s.Src)
}

// AuthData represents the inner payload of authentication responses.
type AuthData struct {
   AccessToken  string `json:"access_token"`
   RefreshToken string `json:"refresh_token"`
   TokenType    string `json:"token_type"`
   ExpiresIn    int    `json:"expires_in"`
}

// ContentNode represents the recursive Server-Driven UI tree used by AMC.
type ContentNode struct {
   Type             string        `json:"type"`
   Properties       *Properties   `json:"properties,omitempty"`
   TabletProperties *Properties   `json:"tablet_properties,omitempty"`
   Children         []ContentNode `json:"children,omitempty"`
   Callback         *Callback     `json:"callback,omitempty"`
}

// Properties holds all possible strongly-typed properties found in the UI nodes.
type Properties struct {
   ID           string `json:"id,omitempty"`
   PageType     string `json:"pageType,omitempty"`
   ManifestType string `json:"manifestType,omitempty"`
   CountryCode  string `json:"countryCode,omitempty"`
   Mode         string `json:"mode,omitempty"`
   Orientation  string `json:"orientation,omitempty"`
   Layout       string `json:"layout,omitempty"`
   Scrollable   bool   `json:"scrollable,omitempty"`
   ContentType  string `json:"contentType,omitempty"`
   Nid          int    `json:"nid,omitempty"`

   Images       *Images       `json:"images,omitempty"`
   Metadata     *Metadata     `json:"metadata,omitempty"`
   Text         *Text         `json:"text,omitempty"`
   DownloadData *DownloadData `json:"downloadData,omitempty"`
   TTS          *TTS          `json:"TTS,omitempty"`
   Navigation   *Navigation   `json:"navigation,omitempty"`
}

type Images struct {
   Default string `json:"default,omitempty"`
   Mobile  string `json:"mobile,omitempty"`
   Tablet  string `json:"tablet,omitempty"`
}

type Metadata struct {
   AmcnID                   string `json:"amcnId,omitempty"`
   EpisodeNumber            int    `json:"episodeNumber,omitempty"`
   ContentNetworkOfRecordID int    `json:"contentNetworkOfRecordId,omitempty"`
   SeasonNumber             int    `json:"seasonNumber,omitempty"`
   ShowName                 string `json:"showName,omitempty"`
   Title                    string `json:"title,omitempty"`
   Nid                      int    `json:"nid,omitempty"`
   PageType                 string `json:"pageType,omitempty"`
   URL                      string `json:"url,omitempty"`
   Action                   string `json:"action,omitempty"`
   ElementType              string `json:"elementType,omitempty"`
   ClickthroughURL          string `json:"clickthroughUrl,omitempty"`
   ElementName              string `json:"elementName,omitempty"`
   ItemText                 string `json:"itemText,omitempty"`
   Label                    string `json:"label,omitempty"`
   NavComponentName         string `json:"navComponentName,omitempty"`
   NavigationTitle          string `json:"navigationTitle,omitempty"`
   IsNavigation             bool   `json:"isNavigation,omitempty"`
   ListTitle                string `json:"listTitle,omitempty"`
   IsPlayback               bool   `json:"isPlayback,omitempty"`
   ListMode                 string `json:"listMode,omitempty"`
   SearchValue              string `json:"searchValue,omitempty"`
   ListPosition             int    `json:"listPosition,omitempty"`
   ComponentName            string `json:"componentName,omitempty"`
}

type Text struct {
   Title       *TextElement `json:"title,omitempty"`
   Description *TextElement `json:"description,omitempty"`
   Subheadings []Subheading `json:"subheadings,omitempty"`
}

type TextElement struct {
   Title string `json:"title,omitempty"`
}

type Subheading struct {
   ID    string `json:"id,omitempty"`
   Title string `json:"title,omitempty"`
   Type  string `json:"type,omitempty"`
}

type DownloadData struct {
   Downloadable        bool      `json:"downloadable,omitempty"`
   DownloadingExpireIn int       `json:"downloadingExpireIn,omitempty"`
   DownloadingEndDate  int       `json:"downloadingEndDate,omitempty"`
   Callback            *Callback `json:"callback,omitempty"`
}

type TTS struct {
   SpeechText string `json:"speechText,omitempty"`
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

type Callback struct {
   Endpoint string `json:"endpoint,omitempty"`
   Type     string `json:"type,omitempty"`
}

// PlaybackData represents the inner streaming and DRM source data.
type PlaybackData struct {
   PlaybackJsonData struct {
      VideoID string   `json:"id"`
      Sources []Source `json:"sources"`
   } `json:"playbackJsonData"`
}

type KeySystems struct {
   ComWidevineAlpha struct {
      LicenseURL string `json:"license_url"`
   } `json:"com.widevine.alpha"`
   ComMicrosoftPlayready struct {
      LicenseURL string `json:"license_url"`
   } `json:"com.microsoft.playready"`
}

// PlaybackResult groups the parsed playback data with the Brightcove JWT needed for DRM.
type PlaybackResult struct {
   Data     PlaybackData
   BcovAuth string
}

// EpisodesMetadata recursively traverses the Server-Driven UI tree
// and extracts only the Metadata for playable episodes.
func (c *ContentNode) EpisodesMetadata() []*Metadata {
   var metadata []*Metadata

   var walk func(node ContentNode)
   walk = func(node ContentNode) {
      if node.Type == "card" && node.Properties != nil && node.Properties.ContentType == "episode" && node.Properties.Metadata != nil {
         metadata = append(metadata, node.Properties.Metadata)
      }
      for _, child := range node.Children {
         walk(child)
      }
   }

   walk(*c)
   return metadata
}

// SeasonsMetadata recursively traverses the Server-Driven UI tree
// and extracts only the Metadata for seasons.
func (c *ContentNode) SeasonsMetadata() []*Metadata {
   var metadata []*Metadata

   var walk func(node ContentNode)
   walk = func(node ContentNode) {
      // Season tabs are identified by being a tab_bar_item with a valid season number
      if node.Type == "tab_bar_item" && node.Properties != nil && node.Properties.Metadata != nil && node.Properties.Metadata.SeasonNumber > 0 {
         metadata = append(metadata, node.Properties.Metadata)
      }
      for _, child := range node.Children {
         walk(child)
      }
   }

   walk(*c)
   return metadata
}

// DashSource finds and returns the first Source with the type "application/dash+xml".
func (p *PlaybackData) DashSource() (*Source, error) {
   for _, src := range p.PlaybackJsonData.Sources {
      if src.Type == "application/dash+xml" {
         return &src, nil
      }
   }
   return nil, fmt.Errorf("application/dash+xml source not found")
}

// String implements the fmt.Stringer interface for easy printing.
func (m *Metadata) String() string {
   if m.SeasonNumber > 0 && m.EpisodeNumber > 0 {
      return fmt.Sprintf("%s S%02dE%02d: %s (ID: %d)", m.ShowName, m.SeasonNumber, m.EpisodeNumber, m.Title, m.Nid)
   }
   if m.SeasonNumber > 0 {
      if m.ShowName != "" {
         return fmt.Sprintf("%s %s (ID: %d)", m.ShowName, m.Title, m.Nid)
      }
      return fmt.Sprintf("%s (ID: %d)", m.Title, m.Nid)
   }
   if m.Title != "" {
      if m.ShowName != "" && m.ShowName != m.Title {
         return fmt.Sprintf("%s: %s (ID: %d)", m.ShowName, m.Title, m.Nid)
      }
      return fmt.Sprintf("%s (ID: %d)", m.Title, m.Nid)
   }
   return fmt.Sprintf("NID: %d", m.Nid)
}

// Login authenticates the user. It requires the guest token (access_token)
// retrieved from calling the Unauth() function.
func Login(guestToken, email, password string) (*AuthData, error) {
   body := map[string]string{
      "email":    email,
      "password": password,
   }
   payload, err := json.Marshal(body)
   if err != nil {
      return nil, err
   }

   req, err := http.NewRequest(http.MethodPost, "https://gw.cds.amcn.com/auth-orchestration-id/api/v1/login", bytes.NewReader(payload))
   if err != nil {
      return nil, err
   }

   req.Header.Set("authorization", "Bearer "+guestToken)
   req.Header.Set("content-type", "application/json")
   req.Header.Set("x-amcn-language", "en")
   req.Header.Set("x-amcn-network", "amcplus")
   req.Header.Set("x-amcn-platform", "web")
   req.Header.Set("x-amcn-service-group-id", "10")
   req.Header.Set("x-amcn-tenant", "amcn")
   req.Header.Set("x-amcn-device-ad-id", "-")
   req.Header.Set("x-amcn-device-id", "-")
   req.Header.Set("x-amcn-service-id", "amcplus")
   req.Header.Set("x-ccpa-do-not-sell", "doNotPassData")

   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()

   if resp.StatusCode != http.StatusOK {
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

func Playback(authToken string, videoID int) (*PlaybackResult, error) {
   url := fmt.Sprintf("https://gw.cds.amcn.com/playback-id/api/v1/playback/%d", videoID)

   payload := []byte(`{"adtags":{"lat":0,"mode":"on-demand","playerHeight":0,"playerWidth":0,"ppid":0,"url":"-"}}`)

   req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(payload))
   if err != nil {
      return nil, err
   }

   req.Header.Set("authorization", "Bearer "+authToken)
   req.Header.Set("content-type", "application/json")
   req.Header.Set("x-amcn-language", "en")
   req.Header.Set("x-amcn-network", "amcplus")
   req.Header.Set("x-amcn-platform", "web")
   req.Header.Set("x-amcn-service-id", "amcplus")
   req.Header.Set("x-amcn-tenant", "amcn")
   req.Header.Set("x-amcn-device-ad-id", "-")
   req.Header.Set("x-ccpa-do-not-sell", "doNotPassData")

   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()

   if resp.StatusCode != http.StatusOK {
      return nil, fmt.Errorf("playback failed with status: %d", resp.StatusCode)
   }

   // Internal envelope to strip the first layer
   var envelope struct {
      Success bool         `json:"success"`
      Status  int          `json:"status"`
      Data    PlaybackData `json:"data"`
   }

   if err := json.NewDecoder(resp.Body).Decode(&envelope); err != nil {
      return nil, err
   }

   return &PlaybackResult{
      Data:     envelope.Data,
      BcovAuth: resp.Header.Get("x-amcn-bc-jwt"),
   }, nil
}

func SeriesDetail(authToken string, seriesID int) (*ContentNode, error) {
   url := fmt.Sprintf("https://gw.cds.amcn.com/content-compiler-cr/api/v1/content/amcn/amcplus/type/series-detail/id/%d", seriesID)

   req, err := http.NewRequest(http.MethodGet, url, nil)
   if err != nil {
      return nil, err
   }

   req.Header.Set("authorization", "Bearer "+authToken)
   req.Header.Set("x-amcn-network", "amcplus")
   req.Header.Set("x-amcn-platform", "android")
   req.Header.Set("x-amcn-tenant", "amcn")

   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()

   if resp.StatusCode != http.StatusOK {
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

func Unauth() (*AuthData, error) {
   req, err := http.NewRequest(http.MethodPost, "https://gw.cds.amcn.com/auth-orchestration-id/api/v1/unauth", nil)
   if err != nil {
      return nil, err
   }

   req.Header.Set("x-amcn-network", "amcplus")
   req.Header.Set("x-amcn-platform", "web")
   req.Header.Set("x-amcn-tenant", "amcn")
   req.Header.Set("x-amcn-device-id", "-")
   req.Header.Set("x-amcn-language", "en")

   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()

   if resp.StatusCode != http.StatusOK {
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

func SeasonEpisodes(authToken string, seasonID int) (*ContentNode, error) {
   url := fmt.Sprintf("https://gw.cds.amcn.com/content-compiler-cr/api/v1/content/amcn/amcplus/type/season-episodes/id/%d", seasonID)

   req, err := http.NewRequest(http.MethodGet, url, nil)
   if err != nil {
      return nil, err
   }

   req.Header.Set("authorization", "Bearer "+authToken)
   req.Header.Set("x-amcn-network", "amcplus")
   req.Header.Set("x-amcn-platform", "android")
   req.Header.Set("x-amcn-tenant", "amcn")

   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()

   if resp.StatusCode != http.StatusOK {
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

func License(licenseURL, bcovAuth string, challengePayload []byte) ([]byte, error) {
   req, err := http.NewRequest(http.MethodPost, licenseURL, bytes.NewReader(challengePayload))
   if err != nil {
      return nil, err
   }

   req.Header.Set("bcov-auth", bcovAuth)

   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()

   if resp.StatusCode != http.StatusOK {
      return nil, fmt.Errorf("license request failed with status: %d", resp.StatusCode)
   }

   return io.ReadAll(resp.Body)
}

func Refresh(refreshToken string) (*AuthData, error) {
   req, err := http.NewRequest(http.MethodPost, "https://gw.cds.amcn.com/auth-orchestration-id/api/v1/refresh", nil)
   if err != nil {
      return nil, err
   }

   req.Header.Set("authorization", "Bearer "+refreshToken)

   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()

   if resp.StatusCode != http.StatusOK {
      return nil, fmt.Errorf("refresh failed with status: %d", resp.StatusCode)
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

type Source struct {
   Codecs     string     `json:"codecs"`
   Src        string     `json:"src"` // MPD
   Type       string     `json:"type"`
   KeySystems KeySystems `json:"key_systems"`
}
