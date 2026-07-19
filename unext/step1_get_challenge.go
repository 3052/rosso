// step1_get_challenge.go
package unext

import (
   "crypto/sha256"
   _ "embed"
   "encoding/base64"
   "fmt"
   "io"
   "log"
   "net/http"
   "net/url"
   "time"
)

//go:embed mad_all_episodes.graphql
var allEpisodesQuery string

//go:embed mad_playlist.graphql
var playlistQuery string

// clientDo wraps http.DefaultClient.Do with a log line so every request is visible.
func clientDo(req *http.Request) (*http.Response, error) {
   log.Println(req.Method, req.URL)
   return http.DefaultClient.Do(req)
}

// clientDoNoRedirect is like clientDo but does not follow redirects.
func clientDoNoRedirect(req *http.Request) (*http.Response, error) {
   log.Println(req.Method, req.URL)
   client := &http.Client{
      CheckRedirect: func(*http.Request, []*http.Request) error {
         return http.ErrUseLastResponse
      },
   }
   return client.Do(req)
}

// generateRandomString generates a 43-character string padded with leading zeros.
// 43 characters satisfies the minimum length requirement for OAuth state, nonce, and PKCE code_verifier.
func generateRandomString() string {
   return fmt.Sprintf("%043d", time.Now().UnixNano())
}

// AuthState holds the PKCE pair and challenge_id produced by Step1.
// It is passed to Step3 and Step4 so they can use the matching values.
type AuthState struct {
   ChallengeID   string
   CodeVerifier  string
   CodeChallenge string
}

// Step1GetChallenge performs the initial GET to /oauth2/auth and extracts
// the challenge_id from the 302 redirect Location header.
//
// state, nonce, and the PKCE pair are generated internally. The PKCE
// pair is returned in AuthState for use by Step3 and Step4.
func Step1GetChallenge() (*AuthState, error) {
   state := generateRandomString()
   nonce := generateRandomString()
   verifier := generateRandomString()

   h := sha256.Sum256([]byte(verifier))
   challenge := base64.RawURLEncoding.EncodeToString(h[:])

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
      return nil, fmt.Errorf("step1: creating request: %w", err)
   }

   req.Header.Set("user-agent", "U-NEXT Phone App Android12 5.71.0 sdk_gphone64_x86_64")

   resp, err := clientDoNoRedirect(req)
   if err != nil {
      return nil, fmt.Errorf("step1: sending request: %w", err)
   }
   defer resp.Body.Close()

   if _, err := io.Copy(io.Discard, resp.Body); err != nil {
      return nil, fmt.Errorf("step1: draining response body: %w", err)
   }

   if resp.StatusCode != http.StatusFound {
      return nil, fmt.Errorf("step1: expected 302, got %d", resp.StatusCode)
   }

   locURL, err := resp.Location()
   if err != nil {
      return nil, fmt.Errorf("step1: getting Location header: %w", err)
   }

   challengeID := locURL.Query().Get("challenge_id")
   if challengeID == "" {
      return nil, fmt.Errorf("step1: challenge_id not found in Location: %s", locURL)
   }

   return &AuthState{
      ChallengeID:   challengeID,
      CodeVerifier:  verifier,
      CodeChallenge: challenge,
   }, nil
}
