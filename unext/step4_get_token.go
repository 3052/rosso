// step4_get_token.go
package unext

import (
   "encoding/json"
   "fmt"
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

// Step4GetToken exchanges the authorization code for access and refresh tokens.
func Step4GetToken(authCode string, auth *AuthState) (*TokenResponse, error) {
   tokenURL := "https://oauth.unext.jp/oauth2/token"
   form := url.Values{}
   form.Set("code", authCode)
   form.Set("code_verifier", auth.CodeVerifier)
   form.Set("grant_type", "authorization_code")
   form.Set("client_id", "unextAndroidApp")
   form.Set("redirect_uri", "jp.unext://page=oauth_callback")
   req, err := http.NewRequest("POST", tokenURL, strings.NewReader(form.Encode()))
   if err != nil {
      return nil, fmt.Errorf("step4: creating request: %w", err)
   }
   req.Header.Set("content-type", "application/x-www-form-urlencoded")
   resp, err := clientDo(req)
   if err != nil {
      return nil, fmt.Errorf("step4: sending request: %w", err)
   }
   defer resp.Body.Close()

   if resp.StatusCode != http.StatusOK {
      return nil, fmt.Errorf("step4: expected 200, got %d", resp.StatusCode)
   }

   var tokenResp TokenResponse
   if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
      return nil, fmt.Errorf("step4: parsing response: %w", err)
   }

   return &tokenResp, nil
}

func (*TokenResponse) CachePath() string {
   return "rosso/unext/TokenResponse"
}
