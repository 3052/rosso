package crave

import (
   "bytes"
   "encoding/base64"
   "encoding/json"
   "fmt"
   "io"
   "net/http"
   "net/url"
   "strings"
)

// GetManifest retrieves the .mpd playback manifest URL from the 9c9media metadata API
func (t *TokenResponse) GetManifest(contentId string, contentPackageId, destinationId int) (string, error) {
   targetURL := fmt.Sprintf(manifestURL, contentId, contentPackageId, destinationId)
   req, _ := http.NewRequest(http.MethodGet, targetURL, nil)
   // Append requested query parameters
   q := req.URL.Query()
   q.Add("format", "mpd")
   req.Header.Set("Authorization", "Bearer "+ t.AccessToken)
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

const (
   graphqlURL  = "https://rte-api.bellmedia.ca/graphql"
   playbackURL = "https://playback.rte-api.bellmedia.ca/contents/%s"
   manifestURL = "https://stream.video.9c9media.com/meta/content/%s/contentpackage/%d/destination/%d/platform/1"
)

const get_showpage = `
query GetShowpage($sessionContext: SessionContext!, $ids: [String!]!) {
   medias(sessionContext: $sessionContext, ids: $ids) {
      firstContent {
         id
      }
   }
}
`

// GetContentID queries the GraphQL API to translate a Media ID to a Content ID
func GetContentId(mediaId string) (string, error) {
   payload := map[string]any{
      "query": get_showpage,
      "variables": map[string]any{
         "ids": []string{mediaId},
         "sessionContext": map[string]string{
            "userLanguage": "EN",
            "userMaturity": "ADULT",
         },
      },
   }
   body, _ := json.Marshal(payload)
   req, _ := http.NewRequest(http.MethodPost, graphqlURL, bytes.NewBuffer(body))
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
         Medias[]struct {
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
   targetURL := fmt.Sprintf(playbackURL, contentId)
   req, _ := http.NewRequest(http.MethodGet, targetURL, nil)
   req.Header.Set("x-playback-language", "EN")
   req.Header.Set("x-client-platform", "platform_jasper_web")
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return 0, 0, err
   }
   defer resp.Body.Close()
   var result struct {
      ContentPackage struct {
         Id            int `json:"id"`
         DestinationID int `json:"destinationId"`
      } `json:"contentPackage"`
   }
   if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
      return 0, 0, err
   }
   if result.ContentPackage.Id == 0 {
      return 0, 0, fmt.Errorf("invalid content package ID received")
   }
   return result.ContentPackage.Id, result.ContentPackage.DestinationID, nil
}
type TokenResponse struct {
   AccessToken  string `json:"access_token"`
   RefreshToken string `json:"refresh_token"`
   AccountId    string `json:"account_id,omitempty"`
   ExpiresIn    int    `json:"expires_in"`
}

type Profile struct {
   Id        string `json:"id"`
   AccountId string `json:"accountId"`
   Nickname  string `json:"nickname"`
   HasPin    bool   `json:"hasPin"`
   Master    bool   `json:"master"`
   Maturity  string `json:"maturity"`
}

const BaseURL = "https://account.bellmedia.ca"

// Basic base64("crave-web:default")
const BasicAuth = "Basic Y3JhdmUtd2ViOmRlZmF1bHQ="

// PasswordLogin performs the initial login to get the first set of tokens.
func PasswordLogin(username, password string) (*TokenResponse, error) {
   endpoint := fmt.Sprintf("%s/api/login/v2.1", BaseURL)
   data := url.Values{}
   data.Set("username", username)
   data.Set("password", password)
   data.Set("grant_type", "password")
   req, err := http.NewRequest("POST", endpoint, strings.NewReader(data.Encode()))
   if err != nil {
      return nil, err
   }
   req.Header.Set("Authorization", BasicAuth)
   req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   if resp.StatusCode < 200 || resp.StatusCode >= 300 {
      body, _ := io.ReadAll(resp.Body)
      return nil, fmt.Errorf("password login failed with status %d: %s", resp.StatusCode, string(body))
   }
   var tokenResp TokenResponse
   if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
      return nil, err
   }
   return &tokenResp, nil
}

// GetProfiles fetches the list of profiles associated with the account.
func GetProfiles(accountId, accessToken string) ([]*Profile, error) {
   endpoint := fmt.Sprintf("%s/api/profile/v2/account/%s", BaseURL, accountId)
   req, err := http.NewRequest("GET", endpoint, nil)
   if err != nil {
      return nil, err
   }
   req.Header.Set("Authorization", "Bearer "+accessToken)
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   if resp.StatusCode < 200 || resp.StatusCode >= 300 {
      body, _ := io.ReadAll(resp.Body)
      return nil, fmt.Errorf("failed to fetch profiles with status %d: %s", resp.StatusCode, string(body))
   }
   var profiles []*Profile
   if err := json.NewDecoder(resp.Body).Decode(&profiles); err != nil {
      return nil, err
   }
   return profiles, nil
}

// ProfileLogin exchanges a refresh token for a fully authorized
// profile-specific Bearer token
func ProfileLogin(refreshToken, profileID string) (*TokenResponse, error) {
   endpoint := fmt.Sprintf("%s/api/login/v2.2", BaseURL)
   data := url.Values{}
   data.Set("grant_type", "refresh_token")
   data.Set("refresh_token", refreshToken)
   data.Set("profile_id", profileID)
   req, err := http.NewRequest("POST", endpoint, strings.NewReader(data.Encode()))
   if err != nil {
      return nil, err
   }
   req.Header.Set("Authorization", BasicAuth)
   req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   if resp.StatusCode < 200 || resp.StatusCode >= 300 {
      body, _ := io.ReadAll(resp.Body)
      return nil, fmt.Errorf("profile login failed with status %d: %s", resp.StatusCode, string(body))
   }
   var tokenResp TokenResponse
   if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
      return nil, err
   }
   return &tokenResp, nil
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
