package kanopy

import (
   "41.neocities.org/maya"
   "encoding/json"
   "errors"
   "io"
   "net/url"
   "path"
   "strconv"
   "strings"
)

// Supports URLs such as:
// - https://kanopy.com/video/6440418
// - https://kanopy.com/video/genius-party
// - https://kanopy.com/en/video/genius-party
// - https://kanopy.com/en/product/genius-party
func ParseVideo(address *url.URL) (*Video, error) {
   if !strings.Contains(address.Host, "kanopy.com") {
      return nil, errors.New("invalid domain")
   }
   // Get the directory of the path (removes the final identifier).
   // e.g., "/en/product/genius-party" -> "/en/product"
   dir := path.Dir(address.Path)
   // Check if the directory ends with "/video" OR "/product".
   // This supports:
   // - /video/{id}
   // - /en/video/{id}
   // - /en/product/{id}
   if !strings.HasSuffix(dir, "/video") && !strings.HasSuffix(dir, "/product") {
      return nil, errors.New("invalid path structure")
   }
   var result Video
   identifier := path.Base(address.Path)
   numeric_id, err := strconv.Atoi(identifier)
   if err != nil {
      result.Alias = identifier
   } else {
      result.VideoId = numeric_id
   }
   return &result, nil
}

type PlayResponse struct {
   PlayId    string     `json:"playId"`
   Manifests []Manifest `json:"manifests"`
   Captions  []Caption  `json:"captions"`
}

type Manifest struct {
   Url            *Url
   ManifestType   string `json:"manifestType"`
   DrmType        string `json:"drmType"`
   StorageService string `json:"storageService"`
   Cdn            string `json:"cdn"`
   DrmLicenseId   string `json:"drmLicenseID"`
}

type Video struct {
   VideoId         int    `json:"videoId"`
   Title           string `json:"title"`
   DescriptionHtml string `json:"descriptionHtml"`
   DurationSeconds int    `json:"durationSeconds"`
   Alias           string `json:"alias"`
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

type PlayRequest struct {
   DomainId int `json:"domainId"`
   UserId   int `json:"userId"`
   VideoId  int `json:"videoId"`
}

type File struct {
   Type string `json:"type"`
   Url  string
}

type Caption struct {
   Language string `json:"language"`
   Files    []File `json:"files"`
   Label    string `json:"label"`
}

type EmailUser struct {
   Email    string `json:"email"`
   Password string `json:"password"`
}

type LoginRequest struct {
   CredentialType string    `json:"credentialType"`
   EmailUser      EmailUser `json:"emailUser"`
}

///

func CreatePlay(loginData *Login, membershipData *Membership, videoData *Video) (*PlayResponse, error) {
   endpoint := &url.URL{
      Scheme: "https",
      Host:   "www.kanopy.com",
      Path:   "/kapi/plays",
   }

   payload := PlayRequest{
      DomainId: membershipData.DomainId,
      UserId:   loginData.UserId,
      VideoId:  videoData.VideoId,
   }

   body, err := json.Marshal(payload)
   if err != nil {
      return nil, err
   }

   headers := map[string]string{
      "authorization": "Bearer " + loginData.Jwt,
   }

   resp, err := maya.Post(endpoint, headers, body)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()

   var play PlayResponse
   if err := json.NewDecoder(resp.Body).Decode(&play); err != nil {
      return nil, err
   }

   return &play, nil
}

type Login struct {
   Jwt               string `json:"jwt"`
   VisitorId         string `json:"visitorId"`
   UserId            int    `json:"userId"`
   KanopyKidsEnabled bool   `json:"kanopyKidsEnabled"`
   WebshopId         int    `json:"webshopId"`
   WebshopCode       string `json:"webshopCode"`
   UserRole          string `json:"userRole"`
}

func LoginUser(email string, password string) (*Login, error) {
   endpoint := &url.URL{
      Scheme: "https",
      Host:   "www.kanopy.com",
      Path:   "/kapi/login",
   }

   payload := LoginRequest{
      CredentialType: "email",
      EmailUser: EmailUser{
         Email:    email,
         Password: password,
      },
   }

   body, err := json.Marshal(payload)
   if err != nil {
      return nil, err
   }

   resp, err := maya.Post(endpoint, nil, body)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()

   var loginData Login
   if err := json.NewDecoder(resp.Body).Decode(&loginData); err != nil {
      return nil, err
   }

   return &loginData, nil
}

func GetMemberships(loginData *Login) ([]Membership, error) {
   endpoint := &url.URL{
      Scheme: "https",
      Host:   "www.kanopy.com",
      Path:   "/kapi/memberships",
   }

   query := url.Values{}
   query.Set("userId", strconv.Itoa(loginData.UserId))
   endpoint.RawQuery = query.Encode()

   headers := map[string]string{
      "authorization": "Bearer " + loginData.Jwt,
   }

   resp, err := maya.Get(endpoint, headers)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   var result struct {
      List []Membership `json:"list"`
   }
   if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
      return nil, err
   }
   return result.List, nil
}

func GetVideo(loginData *Login, alias string) (*Video, error) {
   endpoint := &url.URL{
      Scheme: "https",
      Host:   "www.kanopy.com",
      Path:   "/kapi/videos/alias/" + alias,
   }

   headers := map[string]string{
      "authorization": "Bearer " + loginData.Jwt,
   }

   resp, err := maya.Get(endpoint, headers)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   var result struct {
      Video Video `json:"video"`
   }
   if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
      return nil, err
   }
   return &result.Video, nil
}

func CreateLicense(loginData *Login, manifestData *Manifest, challenge []byte) ([]byte, error) {
   endpoint := &url.URL{
      Scheme: "https",
      Host:   "www.kanopy.com",
      Path:   "/kapi/licenses/widevine/" + manifestData.DrmLicenseId,
   }

   headers := map[string]string{
      "authorization": "Bearer " + loginData.Jwt,
   }

   resp, err := maya.Post(endpoint, headers, challenge)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()

   return io.ReadAll(resp.Body)
}

func (p *PlayResponse) GetDash() (*Manifest, error) {
   for _, manifest_data := range p.Manifests {
      if manifest_data.ManifestType == "dash" {
         return &manifest_data, nil
      }
   }
   return nil, errors.New("dash manifest not found")
}

type Url struct {
   Url url.URL
}

func (u *Url) UnmarshalText(text []byte) error {
   return u.Url.UnmarshalBinary(text)
}

func (u *Url) MarshalText() ([]byte, error) {
   return u.Url.MarshalBinary()
}
