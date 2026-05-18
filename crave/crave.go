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

/*
https://crave.ca/en/movie/anaconda-2025-59881
https://crave.ca/en/play/anaconda-2025-3300246
https://crave.ca/movie/anaconda-2025-59881
https://crave.ca/play/anaconda-2025-3300246
https://crave.ca/play/heated-rivalry/ill-believe-in-anything-s1e5-3233873
*/
func ParseMedia(address *url.URL) (*Media, error) {
   // Split the path directly.
   parts := strings.Split(address.Path, "/")
   if len(parts) < 3 {
      return nil, errors.New("invalid URL path format")
   }
   // Anchor the URL by looking for the explicit media type
   var typePart string
   for _, part := range parts {
      if part == "movie" || part == "play" {
         typePart = part
         break
      }
   }
   if typePart == "" {
      return nil, errors.New("missing media type (movie/play) in URL")
   }
   // Safely grab the last segment (the slug containing the ID)
   lastPart := parts[len(parts)-1]
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
   // Populate struct based on the anchored type
   media_data := &Media{}
   switch typePart {
   case "movie":
      media_data.Id = id
   case "play":
      media_data.FirstContent.Id = id
   }
   return media_data, nil
}

func GetPlayback(token *ProfileToken, activeMedia *Media) (*Playback, error) {
   endpoint := &url.URL{
      Scheme: "https",
      Host:   "playback.rte-api.bellmedia.ca",
      Path:   "/contents/" + strconv.Itoa(activeMedia.FirstContent.Id),
   }

   headers := map[string]string{
      "x-client-platform":   "platform_jasper_web", // platform_jasper_html
      "authorization":       "Bearer " + token.AccessToken,
      "x-playback-language": "EN",
   }

   resp, err := maya.Get(endpoint, headers)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   var result Playback
   if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
      return nil, err
   }
   if result.Error != "" {
      return nil, errors.New(result.Error)
   }
   return &result, nil
}

type Playback struct {
   ContentId      int            `json:"contentId,string"`
   ContentPackage ContentPackage `json:"contentPackage"`
   DestinationId  int            `json:"destinationId"`
   Error          string         // 2026-05-03
}

type Profile struct {
   Id                string   `json:"id"`
   AccountId         string   `json:"accountId"`
   Nickname          string   `json:"nickname"`
   HasPin            bool     `json:"hasPin"`
   Master            bool     `json:"master"`
   Maturity          string   `json:"maturity"`
   Onboarded         bool     `json:"onboarded"`
   UiLanguage        string   `json:"uiLanguage"`
   PlaybackLanguages []string `json:"playbackLanguages"`
   LastModifiedDate  string   `json:"lastModifiedDate"`
   AvatarUrl         string   `json:"avatarUrl"`
}

func GetProfiles(account *AccountToken) ([]Profile, error) {
   endpoint := &url.URL{
      Scheme: "https",
      Host:   "account.bellmedia.ca",
      Path:   "/api/profile/v2/account/" + account.AccountId,
   }

   headers := map[string]string{
      "authorization": "Bearer " + account.AccessToken,
   }

   resp, err := maya.Get(endpoint, headers)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()

   var activeProfiles []Profile
   if err := json.NewDecoder(resp.Body).Decode(&activeProfiles); err != nil {
      return nil, err
   }

   return activeProfiles, nil
}

func (p *Profile) String() string {
   var data strings.Builder
   data.WriteString("nickname: ")
   data.WriteString(p.Nickname)
   if p.HasPin {
      data.WriteString("\nhas pin: true")
   } else {
      data.WriteString("\nhas pin: false")
   }
   if p.Master {
      data.WriteString("\nmaster: true")
   } else {
      data.WriteString("\nmaster: false")
   }
   data.WriteString("\nmaturity: ")
   data.WriteString(p.Maturity)
   data.WriteString("\nid: ")
   data.WriteString(p.Id)
   return data.String()
}

type ProfileToken struct {
   AccessToken  string `json:"access_token"`
   RefreshToken string `json:"refresh_token"`
   Scope        string `json:"scope"`
   TokenType    string `json:"token_type"`
   ExpiresIn    int    `json:"expires_in"`
}

func SwitchProfile(account *AccountToken, profileId string) (*ProfileToken, error) {
   endpoint := &url.URL{
      Scheme: "https",
      Host:   "account.bellmedia.ca",
      Path:   "/api/login/v2.2",
   }

   headers := map[string]string{
      // crave-web:default
      "authorization": "Basic Y3JhdmUtd2ViOmRlZmF1bHQ=",
      "content-type":  "application/x-www-form-urlencoded",
   }

   values := url.Values{}
   values.Set("grant_type", "refresh_token")
   values.Set("profile_id", profileId)
   values.Set("refresh_token", account.RefreshToken)

   body := []byte(values.Encode())

   resp, err := maya.Post(endpoint, headers, body)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()

   token := &ProfileToken{}
   if err := json.NewDecoder(resp.Body).Decode(token); err != nil {
      return nil, err
   }

   return token, nil
}

type Subscription struct {
   Type              string     `json:"type"`
   Experience        Experience `json:"experience"`
   SubscriptionState string     `json:"subscriptionState"`
   StoreName         string     `json:"storeName"`
   ExpirationDate    string     `json:"expirationDate"`
   AutoRenewStatus   bool       `json:"autoRenewStatus"`
}

func GetSubscriptions(token *ProfileToken) ([]Subscription, error) {
   endpoint := &url.URL{
      Scheme: "https",
      Host:   "account.bellmedia.ca",
      Path:   "/api/subscription/v5",
   }

   headers := map[string]string{
      "authorization": "Bearer " + token.AccessToken,
   }

   resp, err := maya.Get(endpoint, headers)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()

   var wrapper struct {
      Subscriptions []Subscription `json:"subscriptions"`
   }
   if err := json.NewDecoder(resp.Body).Decode(&wrapper); err != nil {
      return nil, err
   }

   return wrapper.Subscriptions, nil
}

func (s *Subscription) String() string {
   var data strings.Builder
   data.WriteString("display name: ")
   data.WriteString(s.Experience.DisplayName)
   data.WriteString("\nexpiration date: ")
   data.WriteString(s.ExpirationDate)
   return data.String()
}

//go:embed GetShowpage.gql
var get_showpage string

func GetStream(token *ProfileToken, activePlayback *Playback) (*url.URL, error) {
   endpoint := &url.URL{
      Scheme: "https",
      Host:   "stream.video.9c9media.com",
      Path: fmt.Sprintf(
         "/meta/content/%d/contentpackage/%d/destination/%d/platform/48",
         activePlayback.ContentId, activePlayback.ContentPackage.Id, activePlayback.DestinationId,
      ),
   }
   values := url.Values{}
   values.Set("filter", "ff") // 2160p HEVC
   values.Set("format", "mpd")
   values.Set("hd", "true")  // 1080p H.264
   values.Set("mcv", "true") // H.264 + HEVC
   values.Set("uhd", "true") // 2160p HEVC
   endpoint.RawQuery = values.Encode()

   headers := map[string]string{
      "authorization": "Bearer " + token.AccessToken,
   }

   resp, err := maya.Get(endpoint, headers)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   var result struct {
      Message  string // 2026-05-01
      Playback string `json:"playback"`
   }
   if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
      return nil, err
   }
   if result.Message != "" {
      return nil, errors.New(result.Message)
   }
   return url.Parse(result.Playback)
}

// SL2000 max 2160p
func AcquireLicense(challenge []byte, token *ProfileToken, activePlayback *Playback) ([]byte, error) {
   endpoint := &url.URL{
      Scheme: "https",
      Host:   "license.9c9media.com",
      Path:   "/playready",
   }

   bodyMap := map[string]interface{}{
      "payload": base64.StdEncoding.EncodeToString(challenge),
      "playbackContext": map[string]interface{}{
         "contentId": activePlayback.ContentId,
         // lower-case 'p' as per their API
         "contentpackageId": activePlayback.ContentPackage.Id,
         "destinationId":    activePlayback.DestinationId,
         "jwt":              token.AccessToken,
         "platformId":       48,
      },
   }

   body, err := json.Marshal(bodyMap)
   if err != nil {
      return nil, err
   }

   resp, err := maya.Post(endpoint, nil, body)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   if resp.StatusCode != 200 {
      return nil, errors.New(resp.Status)
   }
   return io.ReadAll(resp.Body)
}

type AccountToken struct {
   AccessToken  string `json:"access_token"`
   RefreshToken string `json:"refresh_token"`
   AccountId    string `json:"account_id"`
   Jti          string `json:"jti"`
}

func PerformLogin(username string, password string) (*AccountToken, error) {
   endpoint := &url.URL{
      Scheme: "https",
      Host:   "account.bellmedia.ca",
      Path:   "/api/login/v2.1",
   }

   headers := map[string]string{
      // crave-web:default
      "authorization": "Basic Y3JhdmUtd2ViOmRlZmF1bHQ=",
      "content-type":  "application/x-www-form-urlencoded",
   }

   values := url.Values{}
   values.Set("grant_type", "password")
   values.Set("password", password)
   values.Set("username", username)

   body := []byte(values.Encode())

   resp, err := maya.Post(endpoint, headers, body)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()

   account := &AccountToken{}
   if err := json.NewDecoder(resp.Body).Decode(account); err != nil {
      return nil, err
   }

   return account, nil
}

type Circle struct {
   Svg ImageSet `json:"svg"`
   Png ImageSet `json:"png"`
}

type ContentPackage struct {
   DurationInSeconds int    `json:"durationInSeconds"`
   Id                int    `json:"id"`
   IsDescribedVideo  bool   `json:"isDescribedVideo"`
   Language          string `json:"language"`
}

type ContentPolicy struct {
   Sku string `json:"sku"`
}

type Experience struct {
   Id              string          `json:"id"`
   Sku             string          `json:"sku"`
   BrandId         string          `json:"brandId"`
   DisplayName     string          `json:"displayName"`
   Logos           Logos           `json:"logos"`
   ContentPolicies []ContentPolicy `json:"contentPolicies"`
}

type FirstContent struct {
   Id int `json:"id,string"`
}

type ImageSet struct {
   Small LocalizedUrl `json:"small"`
}

type LocalizedUrl struct {
   Fr string `json:"fr"`
   En string `json:"en"`
}

type Logos struct {
   Circle Circle `json:"circle"`
}

func GetMedia(showId int) (*Media, error) {
   endpoint := &url.URL{
      Scheme: "https",
      Host:   "rte-api.bellmedia.ca",
      Path:   "/graphql",
   }

   headers := map[string]string{
      // {"platform":"platform_web"}
      "authorization": "Bearer eyJwbGF0Zm9ybSI6InBsYXRmb3JtX3dlYiJ9",
   }

   bodyMap := map[string]interface{}{
      "query": get_showpage,
      "variables": map[string]interface{}{
         "ids": []string{strconv.Itoa(showId)},
         "sessionContext": map[string]interface{}{
            "userLanguage": "EN",
            "userMaturity": "ADULT",
         },
      },
   }

   body, err := json.Marshal(bodyMap)
   if err != nil {
      return nil, err
   }

   resp, err := maya.Post(endpoint, headers, body)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()

   var result struct {
      Data struct {
         Medias []Media `json:"medias"`
      } `json:"data"`
   }
   if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
      return nil, err
   }
   return &result.Data.Medias[0], nil
}

type Media struct {
   FirstContent FirstContent `json:"firstContent"`
   Id           int          `json:"id,string"`
}
