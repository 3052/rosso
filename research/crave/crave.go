package crave

import (
   "encoding/json"
   "fmt"
   "io"
   "net/http"
   "net/url"
   "strings"
)

// PasswordLogin performs the initial login to get the first set of tokens.
func PasswordLogin(username, password, recaptchaToken string) (*TokenResponse, error) {
   endpoint := fmt.Sprintf("%s/api/login/v2.1", BaseURL)

   data := url.Values{}
   data.Set("username", username)
   data.Set("password", password)
   data.Set("grant_type", "password")
   if recaptchaToken != "" {
      data.Set("recaptcha_token", recaptchaToken)
   }

   req, err := http.NewRequest("POST", endpoint, strings.NewReader(data.Encode()))
   if err != nil {
      return nil, err
   }

   req.Header.Set("Authorization", BasicAuth)
   req.Header.Set("Content-Type", "application/x-www-form-urlencoded;charset=utf-8")
   req.Header.Set("User-Agent", UserAgent)
   req.Header.Set("Accept", "application/json, text/plain, */*")

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

type Profile struct {
   ID        string `json:"id"`
   AccountID string `json:"accountId"`
   Nickname  string `json:"nickname"`
   HasPin    bool   `json:"hasPin"`
   Master    bool   `json:"master"`
   Maturity  string `json:"maturity"`
}

// GetProfiles fetches the list of profiles associated with the account.
func GetProfiles(accountID, accessToken string) ([]*Profile, error) {
   endpoint := fmt.Sprintf("%s/api/profile/v2/account/%s", BaseURL, accountID)

   req, err := http.NewRequest("GET", endpoint, nil)
   if err != nil {
      return nil, err
   }

   req.Header.Set("Authorization", "Bearer "+accessToken)
   req.Header.Set("User-Agent", UserAgent)
   req.Header.Set("Accept", "application/json, text/plain, */*")
   req.Header.Set("Origin", "https://www.crave.ca")

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

// ProfileLogin exchanges a refresh token for a fully authorized profile-specific Bearer token.
func ProfileLogin(refreshToken, profileID, profilePin string) (*TokenResponse, error) {
   endpoint := fmt.Sprintf("%s/api/login/v2.2", BaseURL)

   data := url.Values{}
   data.Set("grant_type", "refresh_token")
   data.Set("refresh_token", refreshToken)
   data.Set("profile_id", profileID)

   if profilePin != "" {
      data.Set("profile_pin", profilePin)
   } else {
      data.Set("profile_pin", "")
   }

   req, err := http.NewRequest("POST", endpoint, strings.NewReader(data.Encode()))
   if err != nil {
      return nil, err
   }

   req.Header.Set("Authorization", BasicAuth)
   req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
   req.Header.Set("User-Agent", UserAgent)
   req.Header.Set("Origin", "https://www.crave.ca")

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

func (t *TokenResponse) content() (*http.Response, error) {
   var req http.Request
   req.Header = http.Header{}
   req.URL = &url.URL{}
   req.URL.Host = "stream.video.9c9media.com"
   req.URL.Path = "/meta/content/938361/contentpackage/8143402/destination/1880/platform/1"
   value := url.Values{}
   value["format"] = []string{"mpd"}
   req.Header.Add("Authorization", "Bearer " + t.AccessToken)
   req.URL.RawQuery = value.Encode()
   req.URL.Scheme = "https"
   return http.DefaultClient.Do(&req)
}

type TokenResponse struct {
   AccessToken  string `json:"access_token"`
   RefreshToken string `json:"refresh_token"`
   AccountID    string `json:"account_id,omitempty"`
   ExpiresIn    int    `json:"expires_in"`
}

const BaseURL = "https://account.bellmedia.ca"

// Basic base64("crave-web:default")
const BasicAuth = "Basic Y3JhdmUtd2ViOmRlZmF1bHQ="

const UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:140.0) Gecko/20100101 Firefox/140.0"

///

// MagicLinkLogin consumes the magic link token to obtain new session tokens.
func MagicLinkLogin(magicLinkToken string) (*TokenResponse, error) {
   endpoint := fmt.Sprintf("%s/api/login/v2.2", BaseURL)

   data := url.Values{}
   data.Set("grant_type", "magic_link_token")
   data.Set("magic_link_token", magicLinkToken)

   req, err := http.NewRequest("POST", endpoint, strings.NewReader(data.Encode()))
   if err != nil {
      return nil, err
   }

   req.Header.Set("Authorization", BasicAuth)
   req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
   req.Header.Set("User-Agent", UserAgent)
   req.Header.Set("Origin", "https://www.crave.ca")

   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()

   if resp.StatusCode < 200 || resp.StatusCode >= 300 {
      body, _ := io.ReadAll(resp.Body)
      return nil, fmt.Errorf("magic link login failed with status %d: %s", resp.StatusCode, string(body))
   }

   var tokenResp TokenResponse
   if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
      return nil, err
   }

   return &tokenResp, nil
}
