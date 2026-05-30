package criterion

import (
   "41.neocities.org/maya"
   "encoding/json"
   "errors"
   "fmt"
   "io"
   "net/url"
)

func (*File) CachePath() string {
   return "rosso/criterion/File"
}

type File struct {
   DrmAuthorizationToken string `json:"drm_authorization_token"`
   Links                 struct {
      Source struct {
         Href *Url // MPD
      }
   } `json:"_links"`
   Method string
}

func (*Token) CachePath() string {
   return "rosso/criterion/Token"
}

type Token struct {
   AccessToken  string `json:"access_token"`
   RefreshToken string `json:"refresh_token"`
}

const client_id = "9a87f110f79cd25250f6c7f3a6ec8b9851063ca156dae493bf362a7faf146c78"

func FetchFilesHref(accessToken, slug string) (*url.URL, error) {
   resp, err := maya.Get(
      &url.URL{
         Scheme:   "https",
         Host:     "api.vhx.com",
         Path:     fmt.Sprintf("/collections/%v/items", slug),
         RawQuery: "site_id=59054",
      },
      map[string]string{"authorization": "Bearer " + accessToken},
   )
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   var result struct {
      Embedded struct {
         Items []struct {
            Links struct {
               Files struct {
                  Href Url // https://api.vhx.tv/videos/3460957/files
               }
            } `json:"_links"`
         }
      } `json:"_embedded"`
   }
   err = json.NewDecoder(resp.Body).Decode(&result)
   if err != nil {
      return nil, err
   }
   return &result.Embedded.Items[0].Links.Files.Href.Url, nil
}

func (f *File) FetchWidevine(body []byte) ([]byte, error) {
   resp, err := maya.Post(
      &url.URL{
         Scheme:   "https",
         Host:     "drm.vhx.com",
         Path:     "/v2/widevine",
         RawQuery: url.Values{"token": {f.DrmAuthorizationToken}}.Encode(),
      },
      nil,
      body,
   )
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   return io.ReadAll(resp.Body)
}

func GetDash(files []File) (*File, error) {
   for _, file_data := range files {
      if file_data.Method == "dash" {
         return &file_data, nil
      }
   }
   return nil, errors.New("DASH media file not found")
}

func FetchFiles(accessToken string, files *url.URL) ([]File, error) {
   resp, err := maya.Get(
      files, map[string]string{"authorization": "Bearer " + accessToken},
   )
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

func FetchToken(username, password string) (*Token, error) {
   body := url.Values{
      "client_id":  {client_id},
      "grant_type": {"password"},
      "password":   {password},
      "username":   {username},
   }.Encode()
   resp, err := maya.Post(
      &url.URL{
         Scheme: "https",
         Host:   "auth.vhx.com",
         Path:   "/v1/oauth/token",
      },
      map[string]string{"content-type": "application/x-www-form-urlencoded"},
      []byte(body),
   )
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

func (t *Token) Refresh() error {
   body := url.Values{
      "client_id":     {client_id},
      "grant_type":    {"refresh_token"},
      "refresh_token": {t.RefreshToken},
   }.Encode()
   resp, err := maya.Post(
      &url.URL{
         Scheme: "https",
         Host:   "auth.vhx.com",
         Path:   "/v1/oauth/token",
      },
      map[string]string{"content-type": "application/x-www-form-urlencoded"},
      []byte(body),
   )
   if err != nil {
      return err
   }
   defer resp.Body.Close()
   return json.NewDecoder(resp.Body).Decode(t)
}

type Url struct {
   Url url.URL
}

func (u *Url) MarshalText() ([]byte, error) {
   return u.Url.MarshalBinary()
}

func (u *Url) UnmarshalText(text []byte) error {
   return u.Url.UnmarshalBinary(text)
}
