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
