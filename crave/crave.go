package crave

import (
   "41.neocities.org/maya"
   _ "embed"
   "encoding/base64"
   "encoding/json"
   "errors"
   "fmt"
   "io"
   "net/url"
   "strconv"
   "strings"
)

func FetchContentPackage(accessToken string, contentId int) (*ContentPackage, error) {
   resp, err := maya.Get(
      &url.URL{
         Scheme: "https",
         Host:   "playback.rte-api.bellmedia.ca",
         Path:   fmt.Sprint("/contents/", contentId),
      },
      map[string]string{
         "authorization":       "Bearer " + accessToken,
         "x-client-platform":   "platform_jasper_web", // platform_jasper_html
         "x-playback-language": Language,
      },
   )
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   var result struct {
      ContentPackage ContentPackage
      Error          string // 2026-04-14
      Message string // 2026-04-29
   }
   if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
      return nil, err
   }
   if result.Error != "" {
      return nil, errors.New(result.Error)
   }
   if result.Message != "" {
      return nil, errors.New(result.Message)
   }
   return &result.ContentPackage, nil
}

func (c *ContentPackage) fetchManifest(contentId int, accessToken string, platformId int) (*Manifest, error) {
   resp, err := maya.Get(
      &url.URL{
         Scheme: "https",
         Host:   "stream.video.9c9media.com",
         Path: fmt.Sprintf(
            "/meta/content/%v/contentpackage/%v/destination/%v/platform/%v",
            contentId, c.Id, c.DestinationId, platformId,
         ),
         RawQuery: url.Values{
            "filter": {"ff"}, // 2160p HEVC
            "format": {"mpd"},
            "hd":     {"true"}, // 1080p H.264
            "mcv":    {"true"}, // H.264 + HEVC
            "uhd":    {"true"}, // 2160p HEVC
         }.Encode(),
      },
      map[string]string{"authorization": "Bearer " + accessToken},
   )
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

func (c *ContentPackage) fetchLicense(contentId int, accessToken string, payload []byte, platformId int, path string) ([]byte, error) {
   body, err := json.Marshal(map[string]any{
      "payload": payload,
      "playbackContext": map[string]any{
         "contentId":        contentId,
         "contentpackageId": c.Id, // lower-case 'p' as per their API
         "platformId":       platformId,
         "destinationId":    c.DestinationId,
         "jwt":              accessToken,
      },
   })
   if err != nil {
      return nil, err
   }
   resp, err := maya.Post(
      &url.URL{Scheme: "https", Host: "license.9c9media.com", Path: path},
      nil,
      body,
   )
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()

   body, err = io.ReadAll(resp.Body)
   if err != nil {
      return nil, err
   }
   if resp.StatusCode != 200 {
      var result struct {
         Message string
      }
      err = json.Unmarshal(body, &result)
      if err != nil {
         return nil, err
      }
      return nil, errors.New(result.Message)
   }

   return body, nil
}

func Login(username, password string) (*Account, error) {
   body := url.Values{
      "grant_type": {"password"},
      "password":   {password},
      "username":   {username},
   }.Encode()
   resp, err := maya.Post(
      &url.URL{
         Scheme: "https",
         Host:   "account.bellmedia.ca",
         Path:   "/api/login/v2.1",
      },
      map[string]string{
         "authorization": crave_web,
         "content-type":  "application/x-www-form-urlencoded",
      },
      []byte(body),
   )
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   if resp.StatusCode != 200 {
      return nil, fmt.Errorf("password login failed with: %v", resp.Status)
   }
   result := &Account{}
   if err := json.NewDecoder(resp.Body).Decode(result); err != nil {
      return nil, err
   }
   return result, nil
}

func FetchMedia(id int) (*Media, error) {
   body, err := json.Marshal(map[string]any{
      "query": get_showpage,
      "variables": map[string]any{
         "sessionContext": map[string]string{
            "userLanguage": Language,
            "userMaturity": "ADULT",
         },
         "ids": []string{fmt.Sprint(id)},
      },
   })
   if err != nil {
      return nil, err
   }
   bearer := base64.StdEncoding.EncodeToString(
      []byte(`{ "platform": "platform_web" }`),
   )
   resp, err := maya.Post(
      &url.URL{
         Scheme: "https",
         Host:   "rte-api.bellmedia.ca",
         Path:   "/graphql",
      },
      map[string]string{"authorization": "Bearer " + bearer},
      body,
   )
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

// crave-web:default
const crave_web = "Basic Y3JhdmUtd2ViOmRlZmF1bHQ="

// 699710369328da351ac33c63
func (a *Account) Login(profileId string) error {
   body := url.Values{
      "grant_type":    {"refresh_token"},
      "profile_id":    {profileId},
      "refresh_token": {a.RefreshToken},
   }.Encode()
   resp, err := maya.Post(
      &url.URL{
         Scheme: "https",
         Host:   "account.bellmedia.ca",
         Path:   "/api/login/v2.2",
      },
      map[string]string{
         "authorization": crave_web,
         "content-type":  "application/x-www-form-urlencoded",
      },
      []byte(body),
   )
   if err != nil {
      return err
   }
   defer resp.Body.Close()
   if resp.StatusCode != 200 {
      return fmt.Errorf("profile login failed with: %v", resp.Status)
   }
   return json.NewDecoder(resp.Body).Decode(a)
}

func (a *Account) FetchProfiles() ([]*Profile, error) {
   resp, err := maya.Get(
      &url.URL{
         Scheme: "https",
         Host:   "account.bellmedia.ca",
         Path:   "/api/profile/v2/account/" + a.AccountId,
      },
      map[string]string{"authorization": "Bearer " + a.AccessToken},
   )
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   if resp.StatusCode != 200 {
      return nil, fmt.Errorf("failed to fetch profiles with: %v", resp.Status)
   }
   var profiles []*Profile
   if err := json.NewDecoder(resp.Body).Decode(&profiles); err != nil {
      return nil, err
   }
   return profiles, nil
}

func FetchSubscriptions(accessToken string) ([]Subscription, error) {
   resp, err := maya.Get(
      &url.URL{
         Scheme: "https",
         Host:   "account.bellmedia.ca",
         Path:   "/api/subscription/v5",
      },
      map[string]string{"authorization": "Bearer " + accessToken},
   )
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   var result struct {
      Subscriptions []Subscription
   }
   err = json.NewDecoder(resp.Body).Decode(&result)
   if err != nil {
      return nil, err
   }
   return result.Subscriptions, nil
}

func (m *Manifest) GetManifest() (*url.URL, error) {
   return url.Parse(m.Playback)
}

type Manifest struct {
   Message  string
   Playback string
}

type Account struct {
   AccessToken  string `json:"access_token"`
   AccountId    string `json:"account_id"`
   RefreshToken string `json:"refresh_token"`
}

type Media struct {
   FirstContent struct {
      Id int `json:"id,string"`
   }
   Id int `json:"id,string"`
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

func (s *Subscription) String() string {
   var data strings.Builder
   data.WriteString("display name = ")
   data.WriteString(s.Experience.DisplayName)
   data.WriteString("\nexpiration date = ")
   data.WriteString(s.ExpirationDate)
   return data.String()
}

type Subscription struct {
   Experience struct {
      DisplayName string
   }
   ExpirationDate string
}

var Language = "EN"

//go:embed GetShowpage.gql
var get_showpage string

/*
https://crave.ca/en/movie/anaconda-2025-59881
https://crave.ca/en/play/anaconda-2025-3300246
https://crave.ca/movie/anaconda-2025-59881
https://crave.ca/play/anaconda-2025-3300246
*/
func ParseMedia(rawUrl string) (*Media, error) {
   parsedUrl, err := url.Parse(rawUrl)
   if err != nil {
      return nil, err
   }
   // Split the path directly.
   // e.g., "/en/movie/anaconda-2025-59881" -> ["", "en", "movie", "anaconda-2025-59881"]
   parts := strings.Split(parsedUrl.Path, "/")
   // We need at least 3 parts: the empty string (before the first "/"), the type, and the slug
   if len(parts) < 3 {
      return nil, errors.New("invalid URL path format")
   }
   // Safely grab the last two segments
   lastPart := parts[len(parts)-1] // e.g., "anaconda-2025-59881"
   typePart := parts[len(parts)-2] // e.g., "movie" or "play"
   // Find the last dash to extract the ID
   dashIdx := strings.LastIndex(lastPart, "-")
   if dashIdx == -1 || dashIdx == len(lastPart)-1 {
      return nil, errors.New("no ID found at the end of the URL")
   }
   idStr := lastPart[dashIdx+1:]
   // Convert extracted string to integer
   id, err := strconv.Atoi(idStr)
   if err != nil {
      return nil, fmt.Errorf("invalid ID format: %w", err)
   }
   // Populate struct based on the type
   media := &Media{}
   switch typePart {
   case "movie":
      media.Id = id
   case "play":
      media.FirstContent.Id = id
   default:
      return nil, fmt.Errorf("unknown media type: %s", typePart)
   }
   return media, nil
}

type ContentPackage struct {
   DestinationId int
   Id            int
}

// SL2000 max 2160p
func (c *ContentPackage) LicensePlayReady(contentId int, accessToken string, payload []byte) ([]byte, error) {
   return c.fetchLicense(contentId, accessToken, payload, 48, "/playready")
}

// L3 max 720p
func (c *ContentPackage) LicenseWidevine(contentId int, accessToken string, payload []byte) ([]byte, error) {
   return c.fetchLicense(contentId, accessToken, payload, 1, "/widevine")
}

func (c *ContentPackage) ManifestWidevine(contentId int, accessToken string) (*Manifest, error) {
   return c.fetchManifest(contentId, accessToken, 1)
}

func (c *ContentPackage) ManifestPlayReady(contentId int, accessToken string) (*Manifest, error) {
   return c.fetchManifest(contentId, accessToken, 48)
}
