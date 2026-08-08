package amc

import (
   "bytes"
   "encoding/json"
   "fmt"
   "io"
   "log"
   "net/http"
)

func Do(req *http.Request) (*http.Response, error) {
   log.Println(req.Method, req.URL)
   return http.DefaultClient.Do(req)
}

func License(licenseUrl, bcovAuth string, challenge []byte) ([]byte, error) {
   resp, err := doPost(licenseUrl, map[string]string{"bcov-auth": bcovAuth}, challenge)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   if resp.StatusCode != http.StatusOK {
      return nil, fmt.Errorf("license request failed with status: %d", resp.StatusCode)
   }
   return io.ReadAll(resp.Body)
}

func doGet(rawUrl string, headers map[string]string) (*http.Response, error) {
   req, err := http.NewRequest("GET", rawUrl, nil)
   if err != nil {
      return nil, err
   }
   for key, val := range headers {
      req.Header.Set(key, val)
   }
   return Do(req)
}

func doPost(rawUrl string, headers map[string]string, body []byte) (*http.Response, error) {
   req, err := http.NewRequest("POST", rawUrl, bytes.NewReader(body))
   if err != nil {
      return nil, err
   }
   for key, val := range headers {
      req.Header.Set(key, val)
   }
   return Do(req)
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
   resp, err := doPost(
      "https://gw.cds.amcn.com/auth-orchestration-id/api/v1/login",
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

func Unauth() (*AuthData, error) {
   resp, err := doPost(
      "https://gw.cds.amcn.com/auth-orchestration-id/api/v1/unauth",
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

func (*AuthData) CachePath() string {
   return "rosso/amc/AuthData"
}

func (a *AuthData) Refresh() error {
   resp, err := doPost(
      "https://gw.cds.amcn.com/auth-orchestration-id/api/v1/refresh",
      map[string]string{"authorization": "Bearer " + a.RefreshToken},
      nil,
   )
   if err != nil {
      return err
   }
   defer resp.Body.Close()
   if resp.StatusCode != http.StatusOK {
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

// String implements the fmt.Stringer interface for easy printing.
func (m *Metadata) String() string {
   hasShow := m.ShowName != "" && m.ShowName != "none"
   if m.SeasonNumber > 0 && m.EpisodeNumber > 0 {
      if hasShow {
         return fmt.Sprintf("ShowName: %s\nSeasonNumber: %d\nEpisodeNumber: %d\nTitle: %s\nNID: %d",
            m.ShowName, m.SeasonNumber, m.EpisodeNumber, m.Title, m.Nid)
      }
      return fmt.Sprintf("SeasonNumber: %d\nEpisodeNumber: %d\nTitle: %s\nNID: %d",
         m.SeasonNumber, m.EpisodeNumber, m.Title, m.Nid)
   }
   if m.SeasonNumber > 0 {
      if hasShow {
         return fmt.Sprintf("ShowName: %s\nTitle: %s\nNID: %d", m.ShowName, m.Title, m.Nid)
      }
      return fmt.Sprintf("Title: %s\nNID: %d", m.Title, m.Nid)
   }
   if m.Title != "" {
      if hasShow && m.ShowName != m.Title {
         return fmt.Sprintf("ShowName: %s\nTitle: %s\nNID: %d", m.ShowName, m.Title, m.Nid)
      }
      return fmt.Sprintf("Title: %s\nNID: %d", m.Title, m.Nid)
   }
   return fmt.Sprintf("NID: %d", m.Nid)
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

// Properties holds all possible strongly-typed properties found in the UI
// nodes
type Properties struct {
   ID           string    `json:"id,omitempty"`
   PageType     string    `json:"pageType,omitempty"`
   ManifestType string    `json:"manifestType,omitempty"`
   CountryCode  string    `json:"countryCode,omitempty"`
   Mode         string    `json:"mode,omitempty"`
   Orientation  string    `json:"orientation,omitempty"`
   Layout       string    `json:"layout,omitempty"`
   Scrollable   bool      `json:"scrollable,omitempty"`
   ContentType  string    `json:"contentType,omitempty"`
   Nid          int       `json:"nid,omitempty"`
   Metadata     *Metadata `json:"metadata,omitempty"`
}

// auth.go
