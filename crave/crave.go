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

type Account struct {
   AccessToken  string `json:"access_token"`
   AccountId    string `json:"account_id"`
   RefreshToken string `json:"refresh_token"`
}

// PasswordLogin performs the initial login to get the first set of tokens
func PasswordLogin(username, password string) (*Account, error) {
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

type Profile struct {
   Nickname string `json:"nickname"`
   HasPin   bool   `json:"hasPin"`
   Id       string `json:"id"`
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
   data.WriteString("\nid = ")
   data.WriteString(p.Id)
   return data.String()
}

// ProfileLogin exchanges a refresh token for a fully authorized
// profile-specific Bearer token
func (a *Account) ProfileLogin(profileId string) error {
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

///

var Language = "EN"

//go:embed GetShowpage.gql
var get_showpage string

// GetContentID queries the GraphQL API to translate a Media ID to a Content ID
func GetContentId(mediaId string) (string, error) {
   data, err := json.Marshal(map[string]any{
      "query": get_showpage,
      "variables": map[string]any{
         "ids": []string{mediaId},
         "sessionContext": map[string]string{
            "userLanguage": Language,
            "userMaturity": "ADULT",
         },
      },
   })
   if err != nil {
      return "", err
   }
   req, _ := http.NewRequest(http.MethodPost, graphQlUrl, bytes.NewBuffer(data))
   // The GraphQL endpoint uses a base64 encoded JSON string that includes the access token
   authData := map[string]string{"platform": "platform_web"}
   authBytes, _ := json.Marshal(authData)
   encodedAuth := base64.StdEncoding.EncodeToString(authBytes)
   req.Header.Set("Authorization", "Bearer "+encodedAuth)
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return "", err
   }
   defer resp.Body.Close()
   var result struct {
      Data struct {
         Medias []struct {
            FirstContent struct {
               Id string `json:"id"`
            } `json:"firstContent"`
         } `json:"medias"`
      } `json:"data"`
   }
   if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
      return "", err
   }
   if len(result.Data.Medias) == 0 || result.Data.Medias[0].FirstContent.Id == "" {
      return "", fmt.Errorf("content ID not found in GraphQL response")
   }
   return result.Data.Medias[0].FirstContent.Id, nil
}

// GetPlaybackDetails retrieves the ContentPackage ID and Destination ID
func GetPlaybackDetails(contentId string) (int, int, error) {
   targetUrl := fmt.Sprintf(playbackUrl, contentId)
   req, _ := http.NewRequest(http.MethodGet, targetUrl, nil)
   req.Header.Set("x-playback-language", Language)
   req.Header.Set("x-client-platform", "platform_jasper_web")
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return 0, 0, err
   }
   defer resp.Body.Close()
   var result struct {
      ContentPackage struct {
         Id            int `json:"id"`
         DestinationId int `json:"destinationId"`
      } `json:"contentPackage"`
   }
   if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
      return 0, 0, err
   }
   if result.ContentPackage.Id == 0 {
      return 0, 0, fmt.Errorf("invalid content package ID received")
   }
   return result.ContentPackage.Id, result.ContentPackage.DestinationId, nil
}

// GetManifest retrieves the .mpd playback manifest URL from the 9c9media metadata API
func (a *Account) GetManifest(contentId string, contentPackageId, destinationId int) (string, error) {
   targetUrl := fmt.Sprintf(manifestUrl, contentId, contentPackageId, destinationId)
   req, _ := http.NewRequest(http.MethodGet, targetUrl, nil)
   // Append requested query parameters
   q := req.URL.Query()
   q.Add("format", "mpd")
   req.Header.Set("Authorization", "Bearer "+a.AccessToken)
   req.URL.RawQuery = q.Encode()
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return "", err
   }
   defer resp.Body.Close()
   var result struct {
      Playback string `json:"playback"`
   }
   if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
      return "", err
   }
   if result.Playback == "" {
      return "", fmt.Errorf("playback URL missing in manifest response")
   }
   return result.Playback, nil
}

const graphQlUrl = "https://rte-api.bellmedia.ca/graphql"

const playbackUrl = "https://playback.rte-api.bellmedia.ca/contents/%s"

const manifestUrl = "https://stream.video.9c9media.com/meta/content/%s/contentpackage/%d/destination/%d/platform/1"

// GetWidevineLicense issues the DRM license request using the provided payload
// and the session details
func (a *Account) GetWidevineLicense(session *PlaybackSession, payload string) ([]byte, error) {
   // The API expects the contentId as an integer
   contentIdInt, err := strconv.Atoi(session.ContentId)
   if err != nil {
      return nil, fmt.Errorf("failed to parse content ID to int: %w", err)
   }
   data, err := json.Marshal(WidevineRequest{
      Payload: payload,
      PlaybackContext: PlaybackContext{
         ContentId:        contentIdInt,
         ContentPackageId: session.ContentPackageId,
         PlatformId:       1, // Hardcoded to 1 for Web
         DestinationId:    session.DestinationId,
         Jwt:              a.AccessToken,
      },
   })
   if err != nil {
      return nil, err
   }
   req, err := http.NewRequest(
      http.MethodPost, "https://license.9c9media.com/widevine",
      bytes.NewBuffer(data),
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

// WidevineRequest represents the JSON body needed for the DRM license request
type WidevineRequest struct {
   Payload         string          `json:"payload"`
   PlaybackContext PlaybackContext `json:"playbackContext"`
}

type PlaybackContext struct {
   ContentId        int    `json:"contentId"`
   ContentPackageId int    `json:"contentpackageId"` // Note: lower-case 'p' as per their API
   PlatformId       int    `json:"platformId"`
   DestinationId    int    `json:"destinationId"`
   Jwt              string `json:"jwt"`
}

// PlaybackSession holds the necessary IDs to make subsequent requests (like licensing)
type PlaybackSession struct {
   ContentId        string
   ContentPackageId int
   DestinationId    int
}

// https://www.crave.ca/en/movie/goldeneye-38860"
func extractMediaId(url_data string) (string, error) {
   url_parse, err := url.Parse(url_data)
   if err != nil {
      return "", err
   }
   parts := strings.Split(url_parse.Path, "-")
   if len(parts) == 0 {
      return "", fmt.Errorf("invalid url format")
   }
   return parts[len(parts)-1], nil
}
