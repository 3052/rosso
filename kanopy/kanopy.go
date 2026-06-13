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

func (*Login) CachePath() string {
   return "rosso/kanopy/Login"
}

type PlayResponse struct {
   Captions  []Caption
   Manifests []Manifest
   PlayId    string
}

func CreatePlay(loginData *Login, membershipData *Membership, videoData *Video) (*PlayResponse, error) {
   body, err := json.Marshal(PlayRequest{
      DomainId: membershipData.DomainId,
      UserId:   loginData.UserId,
      VideoId:  videoData.VideoId,
   })
   if err != nil {
      return nil, err
   }
   resp, err := maya.Post(
      &url.URL{
         Scheme: "https",
         Host:   "www.kanopy.com",
         Path:   "/kapi/plays",
      },
      map[string]string{
         "authorization": "Bearer " + loginData.Jwt,
         "content-type":  "application/json",
         "x-version":     "web/undefined/undefined/undefined",
      },
      body,
   )
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

type Video struct {
   VideoId         int    `json:"videoId"`
   Title           string `json:"title"`
   DescriptionHtml string `json:"descriptionHtml"`
   DurationSeconds int    `json:"durationSeconds"`
   Alias           string `json:"alias"`
}

// https://kanopy.com/product/6440418
// https://kanopy.com/en/product/6440418
// https://kanopy.com/video/6440418
// https://kanopy.com/en/irving/product/6440418
// https://kanopy.com/en/irving/product/genius-party
// https://kanopy.com/en/irving/video/6440418
// https://kanopy.com/en/irving/video/genius-party
// https://kanopy.com/en/irving/video/justwatch-6440418
// https://kanopy.com/en/product/genius-party
// https://kanopy.com/en/video/6440418
// https://kanopy.com/en/video/genius-party
// https://kanopy.com/en/video/justwatch-6440418
// https://kanopy.com/irving/product/6440418
// https://kanopy.com/irving/product/genius-party
// https://kanopy.com/irving/video/6440418
// https://kanopy.com/irving/video/genius-party
// https://kanopy.com/irving/video/justwatch-6440418
// https://kanopy.com/product/genius-party
// https://kanopy.com/product/justwatch-6440418
// https://kanopy.com/video/genius-party
// https://kanopy.com/video/justwatch-6440418
func ParseVideo(rawUrl string) (*Video, error) {
   parsedUrl, err := url.Parse(rawUrl)
   if err != nil {
      return nil, err
   }
   slug := path.Base(parsedUrl.Path)
   video := &Video{}
   idStr := strings.TrimPrefix(slug, "justwatch-")
   if id, err := strconv.Atoi(idStr); err == nil {
      video.VideoId = id
   } else {
      video.Alias = slug
   }
   return video, nil
}

///

func (*Manifest) CachePath() string {
   return "rosso/kanopy/Manifest"
}

type Manifest struct {
   Url            *Url
   ManifestType   string `json:"manifestType"`
   DrmType        string `json:"drmType"`
   StorageService string `json:"storageService"`
   Cdn            string `json:"cdn"`
   DrmLicenseId   string `json:"drmLicenseID"`
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

func CreateLicense(loginData *Login, manifestData *Manifest, challenge []byte) ([]byte, error) {
   resp, err := maya.Post(
      &url.URL{
         Scheme: "https",
         Host:   "www.kanopy.com",
         Path:   "/kapi/licenses/widevine/" + manifestData.DrmLicenseId,
      },
      map[string]string{
         "authorization": "Bearer " + loginData.Jwt,
         "x-version":     "web/undefined/undefined/undefined",
      },
      challenge,
   )
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()

   return io.ReadAll(resp.Body)
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
      "x-version":     "web/undefined/undefined/undefined",
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

func LoginUser(email string, password string) (*Login, error) {
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
         Host:   "www.kanopy.com",
         Path:   "/kapi/login",
      },
      map[string]string{"content-type": "application/json"},
      body,
   )
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
