package crave

import (
   "encoding/json"
   "fmt"
   "io"
   "net/http"
   "net/url"
   "strings"
)

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
