package kanopy

import (
   "bytes"
   "encoding/json"
   "fmt"
   "io"
   "net/http"
   "net/url"
)

type PlayRequest struct {
   DomainId int `json:"domainId"`
   UserId   int `json:"userId"`
   VideoId  int `json:"videoId"`
}

type CaptionFile struct {
   Type string `json:"type"`
   Url  string `json:"url"`
}

type Caption struct {
   Files    []CaptionFile `json:"files"`
   Label    string        `json:"label"`
   Language string        `json:"language"`
}

type Dva struct {
   U int `json:"u"`
}

type StudioDrm struct {
   AuthXml      string `json:"authXml"`
   DrmLicenseId string `json:"drmLicenseId"`
}

func (m *Manifest) GetManifest() (*url.URL, error) {
   return url.Parse(m.Url)
}

type Manifest struct {
   Cdn            string    `json:"cdn"`
   DrmLicenseId   string    `json:"drmLicenseID"`
   DrmType        string    `json:"drmType"`
   ManifestType   string    `json:"manifestType"`
   StorageService string    `json:"storageService"`
   StudioDrm      StudioDrm `json:"studioDrm"`
   Url            string    `json:"url"`
}

type PlayResponse struct {
   Captions  []Caption   `json:"captions"`
   Dva       Dva         `json:"dva"`
   Manifests []*Manifest `json:"manifests"`
   PlayId    string      `json:"playId"`
}

// DashManifest returns the manifest with type "dash" or an error if it is not found.
func (p *PlayResponse) DashManifest() (*Manifest, error) {
   for _, manifest := range p.Manifests {
      if manifest.ManifestType == "dash" {
         return manifest, nil
      }
   }
   return nil, fmt.Errorf("dash manifest not found in play response")
}

// CreatePlay registers a playback event using the DomainId from a Membership
// and the VideoId from a Video.
func (s *Session) CreatePlay(membership *Membership, video *Video) (*PlayResponse, error) {
   if membership == nil {
      return nil, fmt.Errorf("membership context is required to create a play")
   }
   if video == nil {
      return nil, fmt.Errorf("video context is required to create a play")
   }

   payload := PlayRequest{
      DomainId: membership.DomainId,
      UserId:   s.UserId,
      VideoId:  video.VideoId,
   }

   body, err := json.Marshal(payload)
   if err != nil {
      return nil, err
   }

   req, err := http.NewRequest("POST", BaseUrl+"/kapi/plays", bytes.NewBuffer(body))
   if err != nil {
      return nil, err
   }

   req.Header.Set("X-Version", Xversion)
   req.Header.Set("Authorization", "Bearer "+s.Jwt)
   req.Header.Set("Content-Type", "application/json")
   req.Header.Set("User-Agent", UserAgent)

   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()

   if resp.StatusCode != http.StatusOK {
      return nil, fmt.Errorf("create play failed with status: %d", resp.StatusCode)
   }

   respBody, err := io.ReadAll(resp.Body)
   if err != nil {
      return nil, err
   }

   var playResp PlayResponse
   if err := json.Unmarshal(respBody, &playResp); err != nil {
      return nil, err
   }

   return &playResp, nil
}

type Membership struct {
   IdentityId         int    `json:"identityId"`
   DomainId           int    `json:"domainId"`
   UserId             int    `json:"userId"`
   Status             string `json:"status"`
   IsDefault          bool   `json:"isDefault"`
   Sitename           string `json:"sitename"`
   Subdomain          string `json:"subdomain"`
   TicketsAvailable   int    `json:"ticketsAvailable"`
   MaxTicketsPerMonth int    `json:"maxTicketsPerMonth"`
}

// GetMemberships fetches the library memberships associated with the session's UserId
// and returns the list of memberships directly.
func (s *Session) GetMemberships() ([]Membership, error) {
   url := fmt.Sprintf("%s/kapi/memberships?userId=%d", BaseUrl, s.UserId)

   req, err := http.NewRequest("GET", url, nil)
   if err != nil {
      return nil, err
   }

   req.Header.Set("User-Agent", UserAgent)
   req.Header.Set("X-Version", Xversion)
   req.Header.Set("Authorization", "Bearer "+s.Jwt)

   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()

   if resp.StatusCode != http.StatusOK {
      return nil, fmt.Errorf("get memberships failed with status: %d", resp.StatusCode)
   }

   respBody, err := io.ReadAll(resp.Body)
   if err != nil {
      return nil, err
   }

   var wrapper struct {
      List []Membership `json:"list"`
   }

   if err := json.Unmarshal(respBody, &wrapper); err != nil {
      return nil, err
   }

   return wrapper.List, nil
}

type LoginRequest struct {
   CredentialType string `json:"credentialType"`
   EmailUser      struct {
      Email    string `json:"email"`
      Password string `json:"password"`
   } `json:"emailUser"`
}

// Login authenticates the user and returns an initialized Session.
func Login(email, password string) (*Session, error) {
   payload := LoginRequest{
      CredentialType: "email",
   }
   payload.EmailUser.Email = email
   payload.EmailUser.Password = password

   body, err := json.Marshal(payload)
   if err != nil {
      return nil, err
   }

   req, err := http.NewRequest("POST", BaseUrl+"/kapi/login", bytes.NewBuffer(body))
   if err != nil {
      return nil, err
   }

   req.Header.Set("Content-Type", "application/json")
   req.Header.Set("User-Agent", UserAgent)

   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()

   if resp.StatusCode != http.StatusOK {
      return nil, fmt.Errorf("login failed with status: %d", resp.StatusCode)
   }

   respBody, err := io.ReadAll(resp.Body)
   if err != nil {
      return nil, err
   }

   var session Session
   if err := json.Unmarshal(respBody, &session); err != nil {
      return nil, err
   }

   return &session, nil
}

func (s *Session) GetWidevine(manifest *Manifest, challenge []byte) ([]byte, error) {
   if manifest == nil {
      return nil, fmt.Errorf("a valid stream manifest is required to request a DRM license")
   }
   if manifest.DrmLicenseId == "" {
      return nil, fmt.Errorf("manifest does not contain a DRM license ID")
   }

   url := fmt.Sprintf("%s/kapi/licenses/widevine/%s", BaseUrl, manifest.DrmLicenseId)

   req, err := http.NewRequest("POST", url, bytes.NewBuffer(challenge))
   if err != nil {
      return nil, err
   }

   req.Header.Set("Authorization", "Bearer "+s.Jwt)
   req.Header.Set("User-Agent", UserAgent)
   req.Header.Set("X-Version", Xversion)

   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()

   if resp.StatusCode != http.StatusOK {
      return nil, fmt.Errorf("widevine license request failed with status: %d", resp.StatusCode)
   }

   return io.ReadAll(resp.Body)
}

const (
   BaseUrl   = "https://www.kanopy.com"
   UserAgent = "!"
   Xversion  = "!/!/!/!"
)

// Session represents an authenticated user context.
type Session struct {
   Jwt       string `json:"jwt"`
   VisitorId string `json:"visitorId"`
   UserId    int    `json:"userId"`
   UserRole  string `json:"userRole"`
}
