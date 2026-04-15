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

func (m *Manifest) GetManifest() (*url.URL, error) {
   return url.Parse(m.Url)
}

func (l *Login) FetchVideo(alias string) (*Video, error) {
   req := http.Request{
      URL: &url.URL{
         Scheme: "https",
         Host:   "www.kanopy.com",
         Path:   "/kapi/videos/alias/" + alias,
      },
      Header: http.Header{},
   }
   req.Header.Set("x-version", x_version)
   req.Header.Set("authorization", "Bearer "+l.Jwt)
   resp, err := http.DefaultClient.Do(&req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   var result struct {
      Video Video
   }
   err = json.NewDecoder(resp.Body).Decode(&result)
   if err != nil {
      return nil, err
   }
   return &result.Video, nil
}

func (l *Login) FetchMembership() (*Membership, error) {
   req := http.Request{
      URL: &url.URL{
         Scheme:   "https",
         Host:     "www.kanopy.com",
         Path:     "/kapi/memberships",
         RawQuery: "userId=" + strconv.Itoa(l.UserId),
      },
      Header: http.Header{},
   }
   req.Header.Set("authorization", "Bearer "+l.Jwt)
   req.Header.Set("user-agent", user_agent)
   req.Header.Set("x-version", x_version)
   resp, err := http.DefaultClient.Do(&req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   if resp.StatusCode != http.StatusOK {
      return nil, errors.New(resp.Status)
   }
   var result struct {
      List []Membership
   }
   err = json.NewDecoder(resp.Body).Decode(&result)
   if err != nil {
      return nil, err
   }
   return &result.List[0], nil
}

func (l *Login) FetchPlays(domainId, videoId int) (*Plays, error) {
   data, err := json.Marshal(map[string]int{
      "domainId": domainId,
      "userId":   l.UserId,
      "videoId":  videoId,
   })
   if err != nil {
      return nil, err
   }
   req, err := http.NewRequest(
      "POST", "https://www.kanopy.com/kapi/plays", bytes.NewReader(data),
   )
   if err != nil {
      return nil, err
   }
   req.Header.Set("authorization", "Bearer "+l.Jwt)
   req.Header.Set("content-type", "application/json")
   req.Header.Set("user-agent", user_agent)
   req.Header.Set("x-version", x_version)
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   var result Plays
   err = json.NewDecoder(resp.Body).Decode(&result)
   if err != nil {
      return nil, err
   }
   if result.ErrorMsgLong != "" {
      return nil, errors.New(result.ErrorMsgLong)
   }
   return &result, nil
}

func (l *Login) FetchWidevine(licenseId string, data []byte) ([]byte, error) {
   req, err := http.NewRequest(
      "POST", "https://www.kanopy.com/kapi/licenses/widevine/"+licenseId,
      bytes.NewReader(data),
   )
   if err != nil {
      return nil, err
   }
   req.Header.Set("authorization", "Bearer "+l.Jwt)
   req.Header.Set("user-agent", user_agent)
   req.Header.Set("x-version", x_version)
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   if resp.StatusCode != http.StatusOK {
      return nil, errors.New(resp.Status)
   }
   return io.ReadAll(resp.Body)
}

func (p *Plays) GetDash() (*Manifest, error) {
   for _, manifest_data := range p.Manifests {
      if manifest_data.ManifestType == "dash" {
         return &manifest_data, nil
      }
   }
   return nil, errors.New("dash manifest not found")
}

type Plays struct {
   Captions []struct {
      Files []struct {
         Url string
      }
   }
   ErrorMsgLong string `json:"error_msg_long"`
   Manifests    []Manifest
}

// good for 10 years
type Login struct {
   Jwt    string
   UserId int
}

type Membership struct {
   DomainId int
}

type Manifest struct {
   DrmLicenseId string
   ManifestType string
   Url          string
}

const user_agent = "!"

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

func FetchLogin(email, password string) (*Login, error) {
   data, err := json.Marshal(map[string]any{
      "credentialType": "email",
      "emailUser": map[string]string{
         "email":    email,
         "password": password,
      },
   })
   if err != nil {
      return nil, err
   }
   req, err := http.NewRequest(
      "POST", "https://www.kanopy.com/kapi/login", bytes.NewReader(data),
   )
   if err != nil {
      return nil, err
   }
   req.Header.Set("content-type", "application/json")
   req.Header.Set("user-agent", user_agent)
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   if resp.StatusCode != http.StatusOK {
      return nil, errors.New(resp.Status)
   }
   result := &Login{}
   err = json.NewDecoder(resp.Body).Decode(result)
   if err != nil {
      return nil, err
   }
   return result, nil
}

type Video struct {
   Alias   string
   VideoId int
}

const x_version = "!/!/!/!"
