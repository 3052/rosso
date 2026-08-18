package kanopy

import (
   "bytes"
   "encoding/json"
   "errors"
   "io"
   "net/http"
   "net/url"
   "path"
   "strconv"
   "strings"
)

func CreateLicense(loginData *Login, manifestData *Manifest, challenge []byte) ([]byte, error) {
   endpoint := &url.URL{
      Scheme: "https",
      Host:   "www.kanopy.com",
      Path:   "/kapi/licenses/widevine/" + manifestData.DrmLicenseId,
   }
   req, err := http.NewRequest("POST", endpoint.String(), bytes.NewReader(challenge))
   if err != nil {
      return nil, err
   }
   req.Header.Set("authorization", "Bearer "+loginData.Jwt)
   req.Header.Set("x-version", "web/undefined/undefined/undefined")
   resp, err := do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()

   return io.ReadAll(resp.Body)
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
   req, err := http.NewRequest("POST", "https://www.kanopy.com/kapi/login", bytes.NewReader(body))
   if err != nil {
      return nil, err
   }
   req.Header.Set("content-type", "application/json")
   resp, err := do(req)
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

   req, err := http.NewRequest("GET", endpoint.String(), nil)
   if err != nil {
      return nil, err
   }
   req.Header.Set("authorization", "Bearer "+loginData.Jwt)
   req.Header.Set("x-version", "web/undefined/undefined/undefined")

   resp, err := do(req)
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

   req, err := http.NewRequest("GET", endpoint.String(), nil)
   if err != nil {
      return nil, err
   }
   req.Header.Set("authorization", "Bearer "+loginData.Jwt)

   resp, err := do(req)
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

func (*Login) CachePath() string {
   return "rosso/kanopy/Login"
}

func (*Manifest) CachePath() string {
   return "rosso/kanopy/Manifest"
}

func (p *PlayResponse) GetDash() (*Manifest, error) {
   for _, manifest_data := range p.Manifests {
      if manifest_data.ManifestType == "dash" {
         return &manifest_data, nil
      }
   }
   return nil, errors.New("dash manifest not found")
}

// kanopy.go
