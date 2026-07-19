// step5_refresh_token.go
package unext

import (
   "encoding/json"
   "fmt"
   "net/http"
   "net/url"
   "strings"
)

// Refresh exchanges the receiver's RefreshToken for a new set of tokens
// and writes the result back into the receiver.
func (t *TokenResponse) Refresh() error {
   tokenURL := "https://oauth.unext.jp/oauth2/token"
   form := url.Values{}
   form.Set("refresh_token", t.RefreshToken)
   form.Set("grant_type", "refresh_token")
   form.Set("client_id", "unextAndroidApp")
   req, err := http.NewRequest("POST", tokenURL, strings.NewReader(form.Encode()))
   if err != nil {
      return fmt.Errorf("refresh: creating request: %w", err)
   }
   req.Header.Set("content-type", "application/x-www-form-urlencoded")
   resp, err := clientDo(req)
   if err != nil {
      return fmt.Errorf("refresh: sending request: %w", err)
   }
   defer resp.Body.Close()

   if resp.StatusCode != http.StatusOK {
      return fmt.Errorf("refresh: expected 200, got %d", resp.StatusCode)
   }

   var newToken TokenResponse
   if err := json.NewDecoder(resp.Body).Decode(&newToken); err != nil {
      return fmt.Errorf("refresh: parsing response: %w", err)
   }

   // Write the new tokens back into the receiver.
   *t = newToken
   return nil
}
