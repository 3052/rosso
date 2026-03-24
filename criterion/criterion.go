package criterion

import (
   "bytes"
   "encoding/json"
   "errors"
   "fmt"
   "io"
   "net/http"
   "net/url"
)

func FetchToken(username, password string) (*Token, error) {
   resp, err := http.PostForm("https://auth.vhx.com/v1/oauth/token", url.Values{
      "client_id":  {client_id},
      "grant_type": {"password"},
      "password":   {password},
      "username":   {username},
   })
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   var result Token
   err = json.NewDecoder(resp.Body).Decode(&result)
   if err != nil {
      return nil, err
   }
   if err := result.AsError(); err != nil {
      return nil, err
   }
   return &result, nil
}

type Token struct {
   AccessToken      string `json:"access_token"`
   Error            string
   ErrorDescription string `json:"error_description"`
   RefreshToken     string `json:"refresh_token"`
}

// AsError returns a standard Go error if the token response was an error,
// otherwise it returns nil.
func (t *Token) AsError() error {
   if t.Error == "" {
      return nil
   }
   return fmt.Errorf("%s: %s", t.Error, t.ErrorDescription)
}

func (m *MediaFile) Dash() (*Dash, error) {
   resp, err := http.Get(m.Links.Source.Href)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   body, err := io.ReadAll(resp.Body)
   if err != nil {
      return nil, err
   }
   return &Dash{Body: body, Url: resp.Request.URL}, nil
}

func (t *Token) Item(slug string) (*VideoItem, error) {
   var req http.Request
   req.URL = &url.URL{
      Scheme:   "https",
      Host:     "api.vhx.com",
      Path: fmt.Sprintf("/collections/%v/items", slug),
      RawQuery: "site_id=59054",
   }
   req.Header = http.Header{}
   req.Header.Set("authorization", "Bearer "+t.AccessToken)
   resp, err := http.DefaultClient.Do(&req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   var result struct {
      Embedded struct {
         Items []VideoItem
      } `json:"_embedded"`
   }
   err = json.NewDecoder(resp.Body).Decode(&result)
   if err != nil {
      return nil, err
   }
   return &result.Embedded.Items[0], nil
}

type VideoItem struct {
   Links struct {
      Files struct {
         Href string // https://api.vhx.tv/videos/3460957/files
      }
   } `json:"_links"`
}

const client_id = "9a87f110f79cd25250f6c7f3a6ec8b9851063ca156dae493bf362a7faf146c78"

type Dash struct {
   Body []byte
   Url  *url.URL
}

type MediaFile struct {
   DrmAuthorizationToken string `json:"drm_authorization_token"`
   Links                 struct {
      Source struct {
         Href string // MPD
      }
   } `json:"_links"`
   Method string
}

func (m *MediaFile) Widevine(data []byte) ([]byte, error) {
   var req http.Request
   req.Method = "POST"
   req.URL = &url.URL{
      Scheme:   "https",
      Host:     "drm.vhx.com",
      Path:     "/v2/widevine",
      RawQuery: url.Values{"token": {m.DrmAuthorizationToken}}.Encode(),
   }
   req.Header = http.Header{}
   req.Body = io.NopCloser(bytes.NewReader(data))
   resp, err := http.DefaultClient.Do(&req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   return io.ReadAll(resp.Body)
}

type MediaFiles []MediaFile

func (m MediaFiles) Dash() (*MediaFile, error) {
   for _, file := range m {
      if file.Method == "dash" {
         return &file, nil
      }
   }
   return nil, errors.New("DASH media file not found")
}

func (t *Token) Refresh() error {
   resp, err := http.PostForm("https://auth.vhx.com/v1/oauth/token", url.Values{
      "client_id":     {client_id},
      "grant_type":    {"refresh_token"},
      "refresh_token": {t.RefreshToken},
   })
   if err != nil {
      return err
   }
   defer resp.Body.Close()
   err = json.NewDecoder(resp.Body).Decode(t)
   if err != nil {
      return err
   }
   return t.AsError()
}

func (t *Token) Files(item *VideoItem) (MediaFiles, error) {
   req, err := http.NewRequest("", item.Links.Files.Href, nil)
   if err != nil {
      return nil, err
   }
   req.Header.Set("authorization", "Bearer "+t.AccessToken)
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   if resp.StatusCode != http.StatusOK {
      return nil, errors.New(resp.Status)
   }
   var result MediaFiles
   err = json.NewDecoder(resp.Body).Decode(&result)
   if err != nil {
      return nil, err
   }
   return result, nil
}
