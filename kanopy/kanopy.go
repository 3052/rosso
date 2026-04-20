package kanopy

import (
   "41.neocities.org/maya"
   "encoding/json"
   "errors"
   "fmt"
   "io"
   "net/url"
   "path"
   "strconv"
   "strings"
)

// CreatePlay registers a playback event using the DomainId from a Membership
// and the VideoId from a Video.
func (s *Session) CreatePlay(membership *Membership, video *Video) (*PlayResponse, error) {
   body, err := json.Marshal(map[string]int{
      "domainId": membership.DomainId,
      "userId":   s.UserId,
      "videoId":  video.VideoId,
   })
   if err != nil {
      return nil, err
   }
   resp, err := maya.Post(
      &url.URL{
         Scheme: "https",
         Host:   host,
         Path:   "/kapi/plays",
      },
      map[string]string{
         "authorization": "Bearer " + s.Jwt,
         "content-type":  "application/json",
         "user-agent":    UserAgent,
         "x-version":     Xversion,
      },
      body,
   )
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()

   if resp.StatusCode != 200 {
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

// GetVideo fetches video metadata and strips the outer wrapper to return the Video object.
func (s *Session) GetVideo(alias string) (*Video, error) {
   resp, err := maya.Get(
      &url.URL{
         Scheme: "https",
         Host:   host,
         Path:   "/kapi/videos/alias/" + alias,
      },
      map[string]string{
         "authorization": "Bearer " + s.Jwt,
         "x-version":     Xversion,
      },
   )
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()

   if resp.StatusCode != 200 {
      return nil, fmt.Errorf("get video failed with status: %d", resp.StatusCode)
   }

   respBody, err := io.ReadAll(resp.Body)
   if err != nil {
      return nil, err
   }

   var wrapper struct {
      Video Video `json:"video"`
   }

   if err := json.Unmarshal(respBody, &wrapper); err != nil {
      return nil, err
   }

   return &wrapper.Video, nil
}

// Session represents an authenticated user context.
type Session struct {
   Jwt       string `json:"jwt"`
   UserId    int    `json:"userId"`
   UserRole  string `json:"userRole"`
   VisitorId string `json:"visitorId"`
}

// Supports URLs such as:
// - https://kanopy.com/video/6440418
// - https://kanopy.com/video/genius-party
// - https://kanopy.com/en/video/genius-party
// - https://kanopy.com/en/product/genius-party
func ParseVideo(urlData string) (*Video, error) {
   url_parse, err := url.Parse(urlData)
   if err != nil {
      return nil, err
   }
   if !strings.Contains(url_parse.Host, "kanopy.com") {
      return nil, errors.New("invalid domain")
   }
   // Get the directory of the path (removes the final identifier).
   // e.g., "/en/product/genius-party" -> "/en/product"
   dir := path.Dir(url_parse.Path)
   // Check if the directory ends with "/video" OR "/product".
   // This supports:
   // - /video/{id}
   // - /en/video/{id}
   // - /en/product/{id}
   if !strings.HasSuffix(dir, "/video") && !strings.HasSuffix(dir, "/product") {
      return nil, errors.New("invalid path structure")
   }
   v := &Video{}
   identifier := path.Base(url_parse.Path)
   numericId, err := strconv.Atoi(identifier)
   if err != nil {
      v.Alias = identifier
   } else {
      v.VideoId = numericId
   }
   return v, nil
}

// Video represents the flattened video metadata, omitting truncated nested objects.
type Video struct {
   Alias                      string `json:"alias"`
   AncestorVideoIds           []int  `json:"ancestorVideoIds"`
   DescriptionHtml            string `json:"descriptionHtml"`
   DurationSeconds            int    `json:"durationSeconds"`
   FeedId                     int    `json:"feedId"`
   HasBurntInCaptions         bool   `json:"hasBurntInCaptions"`
   HasCaptions                bool   `json:"hasCaptions"`
   HasPublicPerformanceRights bool   `json:"hasPublicPerformanceRights"`
   IsFree                     bool   `json:"isFree"`
   IsKids                     bool   `json:"isKids"`
   IsRequestable              bool   `json:"isRequestable"`
   IsSilent                   bool   `json:"isSilent"`
   NeedsTitleTreatment        bool   `json:"needsTitleTreatment"`
   ProductionYear             int    `json:"productionYear"`
   Title                      string `json:"title"`
   VideoId                    int    `json:"videoId"`
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

const (
   host      = "www.kanopy.com"
   UserAgent = "!"
   Xversion  = "!/!/!/!"
)
type Manifest struct {
   Cdn            string    `json:"cdn"`
   DrmLicenseId   string    `json:"drmLicenseID"`
   DrmType        string    `json:"drmType"`
   ManifestType   string    `json:"manifestType"`
   StorageService string    `json:"storageService"`
   StudioDrm      StudioDrm `json:"studioDrm"`
   Url            string    `json:"url"`
}

func (s *Session) GetWidevine(manifest *Manifest, body []byte) ([]byte, error) {
   resp, err := maya.Post(
      &url.URL{
         Scheme: "https",
         Host:   host,
         Path:   "/kapi/licenses/widevine/" + manifest.DrmLicenseId,
      },
      map[string]string{
         "authorization": "Bearer " + s.Jwt,
         "user-agent":    UserAgent,
         "x-version":     Xversion,
      },
      body,
   )
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()

   if resp.StatusCode != 200 {
      return nil, fmt.Errorf("widevine license request failed with status: %d", resp.StatusCode)
   }

   return io.ReadAll(resp.Body)
}

// Login authenticates the user and returns an initialized Session.
func Login(email, password string) (*Session, error) {
   body, err := json.Marshal(map[string]any{
      "credentialType": "email",
      "emailUser": map[string]string{
         "email":    email,
         "password": password,
      },
   })
   if err != nil {
      return nil, err
   }
   resp, err := maya.Post(
      &url.URL{
         Scheme: "https",
         Host:   host,
         Path:   "/kapi/login",
      },
      map[string]string{
         "content-type": "application/json",
         "user-agent":   UserAgent,
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
// GetMemberships fetches the library memberships associated with the session's UserId
// and returns the list of memberships directly.
func (s *Session) GetMemberships() ([]Membership, error) {
   resp, err := maya.Get(
      &url.URL{
         Scheme:   "https",
         Host:     host,
         Path:     "/kapi/memberships",
         RawQuery: fmt.Sprint("userId=", s.UserId),
      },
      map[string]string{
         "authorization": "Bearer " + s.Jwt,
         "user-agent":    UserAgent,
         "x-version":     Xversion,
      },
   )
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()

   if resp.StatusCode != 200 {
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
