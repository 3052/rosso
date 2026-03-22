package crave

import (
   "encoding/json"
   "fmt"
   "io"
   "net/http"
   "net/url"
   "strings"
)

// ProfileLogin exchanges a refresh token for a fully authorized profile-specific Bearer token.
func (c *Client) ProfileLogin(refreshToken, profileID, profilePin string) (*TokenResponse, error) {
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

   resp, err := c.HTTPClient.Do(req)
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
