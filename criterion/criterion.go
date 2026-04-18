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

func (f *File) GetManifest() (*url.URL, error) {
   return url.Parse(f.Links.Source.Href)
}

func GetDash(files []File) (*File, error) {
   for _, file_data := range files {
      if file_data.Method == "dash" {
         return &file_data, nil
      }
   }
   return nil, errors.New("DASH media file not found")
}

// AsError returns a standard Go error if the token response was an error,
// otherwise it returns nil.
func (t *Token) AsError() error {
   if t.Error == "" {
      return nil
   }
   return fmt.Errorf("%s: %s", t.Error, t.ErrorDescription)
}

const client_id = "9a87f110f79cd25250f6c7f3a6ec8b9851063ca156dae493bf362a7faf146c78"

type File struct {
   DrmAuthorizationToken string `json:"drm_authorization_token"`
   Links                 struct {
      Source struct {
         Href string // MPD
      }
   } `json:"_links"`
   Method string
}

func (f *File) FetchWidevine(data []byte) ([]byte, error) {
   req, err := http.NewRequest(
      "POST", "https://drm.vhx.com/v2/widevine", bytes.NewReader(data),
   )
   if err != nil {
      return nil, err
   }
   req.URL.RawQuery = url.Values{"token": {f.DrmAuthorizationToken}}.Encode()
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   return io.ReadAll(resp.Body)
}

type Token struct {
   AccessToken      string `json:"access_token"`
   Error            string
   ErrorDescription string `json:"error_description"`
   RefreshToken     string `json:"refresh_token"`
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
