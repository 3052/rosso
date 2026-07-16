package main

import (
   "encoding/json"
   "fmt"
   "io"
   "net/http"
   "net/url"
   "strings"
)

// TokenResponse is the JSON returned by /oauth2/token.
type TokenResponse struct {
   AccessToken  string `json:"access_token"`
   ExpiresIn    int    `json:"expires_in"`
   RefreshToken string `json:"refresh_token"`
   Scope        string `json:"scope"`
   TokenType    string `json:"token_type"`
}

// step4GetToken exchanges the authorization code for access and refresh tokens.
func step4GetToken(client *http.Client, authCode, codeVerifier string) (*TokenResponse, error) {
   tokenURL := "https://oauth.unext.jp/oauth2/token"

   form := url.Values{}
   form.Set("code", authCode)
   form.Set("grant_type", "authorization_code")
   form.Set("client_id", "unextAndroidApp")
   form.Set("client_secret", "unextAndroidApp")
   form.Set("code_verifier", codeVerifier)
   form.Set("redirect_uri", "jp.unext://page=oauth_callback")

   req, err := http.NewRequest("POST", tokenURL, strings.NewReader(form.Encode()))
   if err != nil {
      return nil, fmt.Errorf("step4: creating request: %w", err)
   }

   req.Header.Set("user-agent", "U-NEXT Phone App Android12 5.71.0 sdk_gphone64_x86_64")
   req.Header.Set("content-type", "application/x-www-form-urlencoded")
   req.Header.Set("accept-encoding", "gzip")

   resp, err := client.Do(req)
   if err != nil {
      return nil, fmt.Errorf("step4: sending request: %w", err)
   }
   defer resp.Body.Close()

   respBody, _ := io.ReadAll(resp.Body)

   if resp.StatusCode != http.StatusOK {
      return nil, fmt.Errorf("step4: expected 200, got %d: %s", resp.StatusCode, string(respBody))
   }

   // Decompress if gzipped (Go handles this automatically via Transport)

   var tokenResp TokenResponse
   if err := json.Unmarshal(respBody, &tokenResp); err != nil {
      return nil, fmt.Errorf("step4: parsing response: %w", err)
   }

   fmt.Printf("[step4] access_token  = %s...\n", tokenResp.AccessToken[:50])
   fmt.Printf("[step4] refresh_token = %s\n", tokenResp.RefreshToken)
   fmt.Printf("[step4] expires_in    = %d\n", tokenResp.ExpiresIn)

   return &tokenResp, nil
}
