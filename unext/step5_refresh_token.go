// step5_refresh_token.go
package unext

import (
   "encoding/json"
   "fmt"
   "io"
   "net/http"
   "net/url"
   "strings"
)

// Refresh exchanges the receiver's RefreshToken for a new set of tokens
// and writes the result back into the receiver.
func (t *TokenResponse) Refresh() error {
   tokenURL := "https://oauth.unext.jp/oauth2/token"

   form := url.Values{}
   form.Set("grant_type", "refresh_token")
   form.Set("client_id", "unextAndroidApp")
   form.Set("client_secret", "unextAndroidApp")
   form.Set("refresh_token", t.RefreshToken)

   req, err := http.NewRequest("POST", tokenURL, strings.NewReader(form.Encode()))
   if err != nil {
      return fmt.Errorf("refresh: creating request: %w", err)
   }

   req.Header.Set("user-agent", "U-NEXT Phone App Android12 5.71.0 sdk_gphone64_x86_64")
   req.Header.Set("content-type", "application/x-www-form-urlencoded")

   resp, err := clientDo(req)
   if err != nil {
      return fmt.Errorf("refresh: sending request: %w", err)
   }
   defer resp.Body.Close()

   respBody, err := io.ReadAll(resp.Body)
   if err != nil {
      return fmt.Errorf("refresh: reading response body: %w", err)
   }

   if resp.StatusCode != http.StatusOK {
      return fmt.Errorf("refresh: expected 200, got %d: %s", resp.StatusCode, string(respBody))
   }

   var newToken TokenResponse
   if err := json.Unmarshal(respBody, &newToken); err != nil {
      return fmt.Errorf("refresh: parsing response: %w (body starts with: %q)", err, string(respBody[:min(len(respBody), 50)]))
   }

   // Write the new tokens back into the receiver.
   *t = newToken
   return nil
}
