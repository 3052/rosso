package criterion

import (
   "bytes"
   "encoding/json"
   "errors"
   "fmt"
   "io"
   "log"
   "net/http"
   "net/url"
   "strings"
)

const client_id = "9a87f110f79cd25250f6c7f3a6ec8b9851063ca156dae493bf362a7faf146c78"

func FetchFilesHref(accessToken, slug string) (string, error) {
   req, err := http.NewRequest(
      "GET",
      fmt.Sprintf("https://api.vhx.com/collections/%v/items?site_id=59054", slug),
      nil,
   )
   if err != nil {
      return "", err
   }
   req.Header.Set("authorization", "Bearer "+accessToken)
   resp, err := do(req)
   if err != nil {
      return "", err
   }
   defer resp.Body.Close()
   var result struct {
      Embedded struct {
         Items []struct {
            Links struct {
               Files struct {
                  Href string // https://api.vhx.tv/videos/3460957/files
               }
            } `json:"_links"`
         }
      } `json:"_embedded"`
   }
   err = json.NewDecoder(resp.Body).Decode(&result)
   if err != nil {
      return "", err
   }
   return result.Embedded.Items[0].Links.Files.Href, nil
}

func do(req *http.Request) (*http.Response, error) {
   log.Println(req.Method, req.URL)
   return http.DefaultClient.Do(req)
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

func FetchFiles(accessToken string, files string) ([]File, error) {
   req, err := http.NewRequest("GET", files, nil)
   if err != nil {
      return nil, err
   }
   req.Header.Set("authorization", "Bearer "+accessToken)
   resp, err := do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   if resp.StatusCode != 200 {
      return nil, errors.New(resp.Status)
   }
   var result []File
   err = json.NewDecoder(resp.Body).Decode(&result)
   if err != nil {
      return nil, err
   }
   return result, nil
}

func GetDash(files []File) (*File, error) {
   for _, file_data := range files {
      if file_data.Method == "dash" {
         return &file_data, nil
      }
   }
   return nil, errors.New("DASH media file not found")
}

func (*File) CachePath() string {
   return "rosso/criterion/File"
}

func (f *File) FetchWidevine(body []byte) ([]byte, error) {
   req, err := http.NewRequest(
      "POST",
      "https://drm.vhx.com/v2/widevine?"+
         url.Values{"token": {f.DrmAuthorizationToken}}.Encode(),
      bytes.NewReader(body),
   )
   if err != nil {
      return nil, err
   }
   resp, err := do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   return io.ReadAll(resp.Body)
}

type Token struct {
   AccessToken  string `json:"access_token"`
   RefreshToken string `json:"refresh_token"`
}

func FetchToken(username, password string) (*Token, error) {
   body := url.Values{
      "client_id":  {client_id},
      "grant_type": {"password"},
      "password":   {password},
      "username":   {username},
   }.Encode()
   req, err := http.NewRequest(
      "POST",
      "https://auth.vhx.com/v1/oauth/token",
      strings.NewReader(body),
   )
   if err != nil {
      return nil, err
   }
   req.Header.Set("content-type", "application/x-www-form-urlencoded")
   resp, err := do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   var result Token
   err = json.NewDecoder(resp.Body).Decode(&result)
   if err != nil {
      return nil, err
   }
   return &result, nil
}

func (*Token) CachePath() string {
   return "rosso/criterion/Token"
}

func (t *Token) Refresh() error {
   body := url.Values{
      "client_id":     {client_id},
      "grant_type":    {"refresh_token"},
      "refresh_token": {t.RefreshToken},
   }.Encode()
   req, err := http.NewRequest(
      "POST",
      "https://auth.vhx.com/v1/oauth/token",
      strings.NewReader(body),
   )
   if err != nil {
      return err
   }
   req.Header.Set("content-type", "application/x-www-form-urlencoded")
   resp, err := do(req)
   if err != nil {
      return err
   }
   defer resp.Body.Close()
   return json.NewDecoder(resp.Body).Decode(t)
}

// criterion.go
