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

func (f *File) ParseDash() (*url.URL, error) {
   return url.Parse(f.Links.Source.Href)
}

type File struct {
   DrmAuthorizationToken string `json:"drm_authorization_token"`
   Links                 struct {
      Source struct {
         Href string // MPD
      }
   } `json:"_links"`
   Method string
}

func GetDash(files []File) (*File, error) {
   for _, file_data := range files {
      if file_data.Method == "dash" {
         return &file_data, nil
      }
   }
   return nil, errors.New("DASH media file not found")
}

type Item struct {
   Links struct {
      Files struct {
         Href string // https://api.vhx.tv/videos/3460957/files
      }
   } `json:"_links"`
}

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

// AsError returns a standard Go error if the token response was an error,
// otherwise it returns nil.
func (t *Token) AsError() error {
   if t.Error == "" {
      return nil
   }
   return fmt.Errorf("%s: %s", t.Error, t.ErrorDescription)
}

const client_id = "9a87f110f79cd25250f6c7f3a6ec8b9851063ca156dae493bf362a7faf146c78"

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

func FetchItem(accessToken, slug string) (*Item, error) {
   req := http.Request{
      URL: &url.URL{
         Scheme:   "https",
         Host:     "api.vhx.com",
         Path:     fmt.Sprintf("/collections/%v/items", slug),
         RawQuery: "site_id=59054",
      },
      Header: http.Header{},
   }
   req.Header.Set("authorization", "Bearer "+accessToken)
   resp, err := http.DefaultClient.Do(&req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   var result struct {
      Embedded struct {
         Items []Item
      } `json:"_embedded"`
   }
   err = json.NewDecoder(resp.Body).Decode(&result)
   if err != nil {
      return nil, err
   }
   return &result.Embedded.Items[0], nil
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

func FetchFiles(accessToken, filesHref string) ([]File, error) {
   req := http.Request{
      Header: http.Header{},
   }
   var err error
   req.URL, err = url.Parse(filesHref)
   if err != nil {
      return nil, err
   }
   req.Header.Set("authorization", "Bearer "+accessToken)
   resp, err := http.DefaultClient.Do(&req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   if resp.StatusCode != http.StatusOK {
      return nil, errors.New(resp.Status)
   }
   var result []File
   err = json.NewDecoder(resp.Body).Decode(&result)
   if err != nil {
      return nil, err
   }
   return result, nil
}
