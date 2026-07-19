// step1_get_challenge.go
package unext

import (
   "crypto/rand"
   "crypto/sha256"
   _ "embed"
   "encoding/base64"
   "fmt"
   "io"
   "log"
   "net/http"
   "net/url"
   "strings"
)

// Pre-minified at package init; never recomputed.
var (
   minPlaylistQuery    = gqlMinify(rawPlaylistQuery)
   minAllEpisodesQuery = gqlMinify(rawAllEpisodesQuery)
   minVideoDetailQuery = gqlMinify(rawVideoDetailQuery)
)

// DefaultClient is the http.Client used by all Step* functions.
// Set a CookieJar on it before calling Step1 if you need session
// persistence across steps (which you do).
var DefaultClient = &http.Client{}

//go:embed mad_all_episodes.graphql
var rawAllEpisodesQuery string

//go:embed mad_playlist.graphql
var rawPlaylistQuery string

//go:embed mad_video_detail.graphql
var rawVideoDetailQuery string

// generateRandomString generates a URL-safe random string of the given length.
func GenerateRandomString(length int) (string, error) {
   b := make([]byte, length)
   _, err := rand.Read(b)
   if err != nil {
      return "", err
   }
   return base64.RawURLEncoding.EncodeToString(b)[:length], nil
}

// pkcePair generates a code_verifier and its corresponding code_challenge (S256).
func PkcePair() (verifier string, challenge string, err error) {
   verifier, err = GenerateRandomString(43)
   if err != nil {
      return "", "", err
   }

   h := sha256.Sum256([]byte(verifier))
   challenge = base64.RawURLEncoding.EncodeToString(h[:])
   return verifier, challenge, nil
}

// Step1GetChallenge performs the initial GET to /oauth2/auth and extracts
// the challenge_id from the 302 redirect Location header.
func Step1GetChallenge(state, nonce string) (string, error) {
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
   resp, err := clientDo(req)
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

// clientDo wraps DefaultClient.Do with a log line so every request is visible.
func clientDo(req *http.Request) (*http.Response, error) {
   log.Println(req.Method, req.URL)
   return DefaultClient.Do(req)
}

// gqlMinify collapses insignificant whitespace in a GraphQL operation
// string. None of the embedded queries contain string literals or
// comments, so a pure whitespace-collapser is sufficient.
func gqlMinify(s string) string {
   var b strings.Builder
   b.Grow(len(s))

   prevSpace := false
   for i := 0; i < len(s); i++ {
      c := s[i]
      if c == ' ' || c == '\t' || c == '\n' || c == '\r' {
         prevSpace = true
         continue
      }
      if prevSpace {
         b.WriteByte(' ')
         prevSpace = false
      }
      b.WriteByte(c)
   }

   return b.String()
}
