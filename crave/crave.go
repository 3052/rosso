package crave

import (
   "bytes"
   _ "embed"
   "encoding/base64"
   "encoding/json"
   "errors"
   "fmt"
   "io"
   "net/http"
   "net/url"
   "strconv"
   "strings"
)

// https://crave.ca/movie/anaconda-2025-59881
// https://crave.ca/movie/goldeneye-38860
func ParseMediaId(urlData string) (int, error) {
   idx := strings.LastIndex(urlData, "-")
   if idx == -1 {
      return 0, strconv.ErrSyntax
   }
   return strconv.Atoi(urlData[idx+1:])
}

type Account struct {
   AccessToken  string `json:"access_token"`
   AccountId    string `json:"account_id"`
   RefreshToken string `json:"refresh_token"`
}

func Login(username, password string) (*Account, error) {
   data := url.Values{
      "grant_type": {"password"},
      "password":   {password},
      "username":   {username},
   }.Encode()
   req, err := http.NewRequest(
      "POST", "https://account.bellmedia.ca/api/login/v2.1",
      strings.NewReader(data),
   )
   if err != nil {
      return nil, err
   }
   req.Header.Set("content-type", "application/x-www-form-urlencoded")
   req.SetBasicAuth("crave-web", "default")
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   if resp.StatusCode != http.StatusOK {
      return nil, fmt.Errorf("password login failed with: %v", resp.Status)
   }
   result := &Account{}
   if err := json.NewDecoder(resp.Body).Decode(result); err != nil {
      return nil, err
   }
   return result, nil
}

// 699710369328da351ac33c63
func (a *Account) Login(profileId string) error {
   data := url.Values{
      "grant_type":    {"refresh_token"},
      "profile_id":    {profileId},
      "refresh_token": {a.RefreshToken},
   }.Encode()
   req, err := http.NewRequest(
      "POST", "https://account.bellmedia.ca/api/login/v2.2",
      strings.NewReader(data),
   )
   if err != nil {
      return err
   }
   req.Header.Set("content-type", "application/x-www-form-urlencoded")
   req.SetBasicAuth("crave-web", "default")
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return err
   }
   defer resp.Body.Close()
   if resp.StatusCode != http.StatusOK {
      return fmt.Errorf("profile login failed with: %v", resp.Status)
   }
   return json.NewDecoder(resp.Body).Decode(a)
}

func (c *ContentPackage) FetchWidevine(contentId int, accessToken string, payload []byte) ([]byte, error) {
   data, err := json.Marshal(map[string]any{
      "payload": payload,
      "playbackContext": map[string]any{
         "contentId":        contentId,
         "contentpackageId": c.Id, // lower-case 'p' as per their API
         "platformId":       1,    // Hardcoded to 1 for Web
         "destinationId":    c.DestinationId,
         "jwt":              accessToken,
      },
   })
   if err != nil {
      return nil, err
   }
   req, err := http.NewRequest(
      "POST", "https://license.9c9media.com/widevine", bytes.NewBuffer(data),
   )
   if err != nil {
      return nil, err
   }
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   data, err = io.ReadAll(resp.Body)
   if err != nil {
      return nil, err
   }
   if resp.StatusCode != http.StatusOK {
      var result struct {
         Message string
      }
      err = json.Unmarshal(data, &result)
      if err != nil {
         return nil, err
      }
      return nil, errors.New(result.Message)
   }
   // The response is usually a binary widevine license
   return data, nil
}

type ContentPackage struct {
   Id            int
   DestinationId int
}

type Dash struct {
   Body []byte
   Url  *url.URL
}

type Manifest struct {
   Message  string
   Playback string
}

func (m *Manifest) FetchDash() (*Dash, error) {
   resp, err := http.Get(m.Playback)
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

func (a *Account) FetchProfiles() ([]*Profile, error) {
   req := http.Request{
      URL: &url.URL{
         Scheme: "https",
         Host:   "account.bellmedia.ca",
         Path:   "/api/profile/v2/account/" + a.AccountId,
      },
      Header: http.Header{},
   }
   req.Header.Set("authorization", "Bearer "+a.AccessToken)
   resp, err := http.DefaultClient.Do(&req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   if resp.StatusCode != http.StatusOK {
      return nil, fmt.Errorf("failed to fetch profiles with: %v", resp.Status)
   }
   var profiles []*Profile
   if err := json.NewDecoder(resp.Body).Decode(&profiles); err != nil {
      return nil, err
   }
   return profiles, nil
}

func (c *ContentPackage) FetchManifest(contentId int, accessToken string) (*Manifest, error) {
   req := http.Request{
      URL: &url.URL{
         Scheme: "https",
         Host:   "stream.video.9c9media.com",
         Path: fmt.Sprintf(
            "/meta/content/%v/contentpackage/%v/destination/%v/platform/1",
            contentId, c.Id, c.DestinationId,
         ),
         RawQuery: url.Values{
            "filter": {"ff"}, // 1080p
            "format": {"mpd"},
            "hd":     {"true"}, // 1080p
         }.Encode(),
      },
      Header: http.Header{},
   }
   req.Header.Set("authorization", "Bearer "+accessToken)
   resp, err := http.DefaultClient.Do(&req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   var result Manifest
   err = json.NewDecoder(resp.Body).Decode(&result)
   if err != nil {
      return nil, err
   }
   if result.Message != "" {
      return nil, errors.New(result.Message)
   }
   return &result, nil
}

type Media struct {
   FirstContent struct {
      Id int `json:"id,string"`
   }
}

func (m Media) FetchContentPackage() (*ContentPackage, error) {
   req := http.Request{
      URL: &url.URL{
         Scheme: "https",
         Host:   "playback.rte-api.bellmedia.ca",
         Path:   "/contents/" + strconv.Itoa(m.FirstContent.Id),
      },
      Header: http.Header{},
   }
   req.Header.Set("x-playback-language", Language)
   req.Header.Set("x-client-platform", "platform_jasper_web")
   resp, err := http.DefaultClient.Do(&req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   var result struct {
      ContentPackage ContentPackage
   }
   if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
      return nil, err
   }
   return &result.ContentPackage, nil
}

func FetchMedia(id int) (*Media, error) {
   body, err := json.Marshal(map[string]any{
      "query": get_showpage,
      "variables": map[string]any{
         "sessionContext": map[string]string{
            "userLanguage": Language,
            "userMaturity": "ADULT",
         },
         "ids": []string{strconv.Itoa(id)},
      },
   })
   if err != nil {
      return nil, err
   }
   req, err := http.NewRequest(
      "POST", "https://rte-api.bellmedia.ca/graphql", bytes.NewBuffer(body),
   )
   if err != nil {
      return nil, err
   }
   bearer := base64.StdEncoding.EncodeToString(
      []byte(`{ "platform": "platform_web" }`),
   )
   req.Header.Set("Authorization", "Bearer "+bearer)
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   var result struct {
      Data struct {
         Medias []Media
      }
   }
   if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
      return nil, err
   }
   if len(result.Data.Medias) == 0 || result.Data.Medias[0].FirstContent.Id == 0 {
      return nil, errors.New("content ID not found in GraphQL response")
   }
   return &result.Data.Medias[0], nil
}

func (p *Profile) String() string {
   var data strings.Builder
   data.WriteString("nickname = ")
   data.WriteString(p.Nickname)
   if p.HasPin {
      data.WriteString("\nhas pin = true")
   } else {
      data.WriteString("\nhas pin = false")
   }
   if p.Master {
      data.WriteString("\nmaster = true")
   } else {
      data.WriteString("\nmaster = false")
   }
   data.WriteString("\nmaturity = ")
   data.WriteString(p.Maturity)
   data.WriteString("\nid = ")
   data.WriteString(p.Id)
   return data.String()
}

type Profile struct {
   Nickname string `json:"nickname"`
   HasPin   bool   `json:"hasPin"`
   Master   bool
   Maturity string
   Id       string `json:"id"`
}

var Language = "EN"

//go:embed GetShowpage.gql
var get_showpage string
