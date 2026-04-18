package criterion

import (
   "41.neocities.org/maya"
   "encoding/json"
   "errors"
   "fmt"
   "net/http"
   "net/url"
)

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

func FetchFilesHref(accessToken, slug string) (string, error) {
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
   if err := result.AsError(); err != nil {
      return nil, err
   }
   return &result, nil
}
