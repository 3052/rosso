package unext

import (
   "encoding/json"
   "fmt"
   "io"
   "net/http"
   "net/url"
   "strings"
)

func min(a, b int) int {
   if a < b {
      return a
   }
   return b
}

// TokenResponse is the JSON returned by /oauth2/token.
type TokenResponse struct {
   AccessToken  string `json:"access_token"`
   ExpiresIn    int    `json:"expires_in"`
   RefreshToken string `json:"refresh_token"`
   Scope        string `json:"scope"`
   TokenType    string `json:"token_type"`
}

// Step4GetToken exchanges the authorization code for access and refresh tokens.
func Step4GetToken(client *http.Client, authCode, codeVerifier string) (*TokenResponse, error) {
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

   resp, err := client.Do(req)
   if err != nil {
      return nil, fmt.Errorf("step4: sending request: %w", err)
   }
   defer resp.Body.Close()

   respBody, _ := io.ReadAll(resp.Body)

   if resp.StatusCode != http.StatusOK {
      return nil, fmt.Errorf("step4: expected 200, got %d: %s", resp.StatusCode, string(respBody))
   }

   var tokenResp TokenResponse
   if err := json.Unmarshal(respBody, &tokenResp); err != nil {
      return nil, fmt.Errorf("step4: parsing response: %w (body starts with: %q)", err, string(respBody[:min(len(respBody), 50)]))
   }

   return &tokenResp, nil
}
