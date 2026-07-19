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

func GeneratePKCE() string {
   sha := sha256.Sum256([]byte(time.Now().String()))
   return base64.RawURLEncoding.EncodeToString(sha[:])
}

func clientDo(req *http.Request) (*http.Response, error) {
   log.Println(req.Method, req.URL)
   return http.DefaultClient.Do(req)
}

func clientDoNoRedirect(req *http.Request) (*http.Response, error) {
   log.Println(req.Method, req.URL)
   client := &http.Client{
      CheckRedirect: func(*http.Request, []*http.Request) error {
         return http.ErrUseLastResponse
      },
   }
   return client.Do(req)
}

// AuthState holds the challenge_id and PKCE verifier produced by Step1.
type AuthState struct {
   ChallengeID  string
   CodeVerifier string
}

// Step1GetChallenge performs the initial GET to /oauth2/auth and extracts
// the challenge_id from the 302 redirect Location header.
func Step1GetChallenge() (*AuthState, error) {
   baseURL := "https://oauth.unext.jp/oauth2/auth"
   pkce := GeneratePKCE()
   params := url.Values{}
   params.Set("nonce", pkce)
   params.Set("state", pkce)
   params.Set("scope", "offline unext")
   params.Set("response_type", "code")
   params.Set("client_id", "unextAndroidApp")
   params.Set("redirect_uri", "jp.unext://page=oauth_callback")
   req, err := http.NewRequest("GET", baseURL+"?"+params.Encode(), nil)
   if err != nil {
      return nil, fmt.Errorf("step1: creating request: %w", err)
   }
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
      ChallengeID:  challengeID,
      CodeVerifier: pkce,
   }, nil
}
