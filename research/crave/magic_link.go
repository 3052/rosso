package crave

import (
   "encoding/json"
   "fmt"
   "io"
   "net/http"
   "net/url"
   "strings"
)

// GenerateMagicLink generates a magic link token used for SSO across Bell Media domains.
func GenerateMagicLink(accessToken string) (string, error) {
   endpoint := fmt.Sprintf("%s/api/magic-link/v2.1/generate", BaseURL)

   req, err := http.NewRequest("POST", endpoint, nil)
   if err != nil {
      return "", err
   }

   req.Header.Set("Authorization", "Bearer "+accessToken)
   req.Header.Set("User-Agent", UserAgent)
   req.Header.Set("Accept", "application/json, text/plain, */*")

   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return "", err
   }
   defer resp.Body.Close()

   if resp.StatusCode < 200 || resp.StatusCode >= 300 {
      body, _ := io.ReadAll(resp.Body)
      return "", fmt.Errorf("magic link generation failed with status %d: %s", resp.StatusCode, string(body))
   }

   body, err := io.ReadAll(resp.Body)
   if err != nil {
      return "", err
   }

   // The response is a raw JSON string like "MgZ_TtZxjd..."
   return strings.Trim(string(body), `"`), nil
}

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
