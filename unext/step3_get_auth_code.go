// step3_get_auth_code.go
package unext

import (
   "fmt"
   "io"
   "net/http"
   "net/url"
   "strings"
)

// Step3GetAuthCode sends the code_challenge to the post_auth_endpoint and
// extracts the authorization code from the 302 redirect Location header.
func Step3GetAuthCode(postAuthEndpoint string, auth *AuthState) (string, error) {
   fullURL := "https://oauth.unext.jp" + postAuthEndpoint
   form := url.Values{}
   form.Set("code_challenge", auth.CodeChallenge)
   form.Set("code_challenge_method", "S256")
   req, err := http.NewRequest("POST", fullURL, strings.NewReader(form.Encode()))
   if err != nil {
      return "", fmt.Errorf("step3: creating request: %w", err)
   }
   req.Header.Set("content-type", "application/x-www-form-urlencoded")
   resp, err := clientDoNoRedirect(req)
   if err != nil {
      return "", fmt.Errorf("step3: sending request: %w", err)
   }
   defer resp.Body.Close()

   if _, err := io.Copy(io.Discard, resp.Body); err != nil {
      return "", fmt.Errorf("step3: draining response body: %w", err)
   }

   if resp.StatusCode != http.StatusFound {
      return "", fmt.Errorf("step3: expected 302, got %d", resp.StatusCode)
   }

   locURL, err := resp.Location()
   if err != nil {
      return "", fmt.Errorf("step3: getting Location header: %w", err)
   }

   code := locURL.Query().Get("code")
   if code == "" {
      return "", fmt.Errorf("step3: code not found in Location: %s", locURL)
   }

   return code, nil
}
