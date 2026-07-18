// step1_get_challenge.go
package unext

import (
   "fmt"
   "io"
   "net/http"
   "net/url"
)

// Step1GetChallenge performs the initial GET to /oauth2/auth and extracts
// the challenge_id from the 302 redirect Location header.
func Step1GetChallenge(client *http.Client, state, nonce string) (string, error) {
   baseURL := "https://oauth.unext.jp/oauth2/auth"

   params := url.Values{}
   params.Set("state", state)
   params.Set("scope", "offline unext")
   params.Set("nonce", nonce)
   params.Set("response_type", "code")
   params.Set("client_id", "unextAndroidApp")
   params.Set("redirect_uri", "jp.unext://page=oauth_callback")

   req, err := http.NewRequest("GET", baseURL+"?"+params.Encode(), nil)
   if err != nil {
      return "", fmt.Errorf("step1: creating request: %w", err)
   }

   req.Header.Set("user-agent", "U-NEXT Phone App Android12 5.71.0 sdk_gphone64_x86_64")

   // Do NOT follow redirects — we need the Location header.
   resp, err := client.Do(req)
   if err != nil {
      return "", fmt.Errorf("step1: sending request: %w", err)
   }
   defer resp.Body.Close()
   io.Copy(io.Discard, resp.Body)

   if resp.StatusCode != http.StatusFound {
      return "", fmt.Errorf("step1: expected 302, got %d", resp.StatusCode)
   }

   location := resp.Header.Get("Location")
   if location == "" {
      return "", fmt.Errorf("step1: no Location header in response")
   }

   // Location looks like: https://oauth.unext.jp/login?challenge_id=cc4e1aed-...
   locURL, err := url.Parse(location)
   if err != nil {
      return "", fmt.Errorf("step1: parsing Location: %w", err)
   }

   challengeID := locURL.Query().Get("challenge_id")
   if challengeID == "" {
      return "", fmt.Errorf("step1: challenge_id not found in Location: %s", location)
   }

   return challengeID, nil
}
