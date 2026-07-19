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
func Step3GetAuthCode(postAuthEndpoint, codeChallenge string) (string, error) {
   // postAuthEndpoint is a path like:
   // /oauth2/auth?challenge_id=...&client_id=...&nonce=...&redirect_uri=...&response_type=code&scope=offline+unext&state=...
   fullURL := "https://oauth.unext.jp" + postAuthEndpoint

   form := url.Values{}
   form.Set("code_challenge", codeChallenge)
   form.Set("code_challenge_method", "S256")

   req, err := http.NewRequest("POST", fullURL, strings.NewReader(form.Encode()))
   if err != nil {
      return "", fmt.Errorf("step3: creating request: %w", err)
   }

   req.Header.Set("user-agent", "U-NEXT Phone App Android12 5.71.0 sdk_gphone64_x86_64")
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
