package unext

import (
   "bytes"
   "crypto/sha256"
   _ "embed"
   "encoding/base64"
   "fmt"
   "io"
   "log"
   "net/http"
   "net/url"
   "strconv"
   "strings"
   "time"
)

// Pre-minified at package init; never recomputed.
var (
   minPlaylistQuery    = gqlMinify(rawPlaylistQuery)
   minAllEpisodesQuery = gqlMinify(rawAllEpisodesQuery)
   minVideoDetailQuery = gqlMinify(rawVideoDetailQuery)
)

//go:embed mad_all_episodes.graphql
var rawAllEpisodesQuery string

//go:embed mad_playlist.graphql
var rawPlaylistQuery string

//go:embed mad_video_detail.graphql
var rawVideoDetailQuery string

// generateRandomString generates a URL-safe string of the given length,
// seeded from the current time. Different every nanosecond.
func GenerateRandomString(length int) (string, error) {
   s := strconv.FormatInt(time.Now().UnixNano(), 10)
   var buf bytes.Buffer
   for buf.Len()*4/3 < length {
      buf.WriteString(s)
   }
   return base64.RawURLEncoding.EncodeToString(buf.Bytes())[:length], nil
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

   resp, err := clientDoNoRedirect(req)
   if err != nil {
      return "", fmt.Errorf("step1: sending request: %w", err)
   }
   defer resp.Body.Close()

   if _, err := io.Copy(io.Discard, resp.Body); err != nil {
      return "", fmt.Errorf("step1: draining response body: %w", err)
   }

   if resp.StatusCode != http.StatusFound {
      return "", fmt.Errorf("step1: expected 302, got %d", resp.StatusCode)
   }

   locURL, err := resp.Location()
   if err != nil {
      return "", fmt.Errorf("step1: getting Location header: %w", err)
   }

   challengeID := locURL.Query().Get("challenge_id")
   if challengeID == "" {
      return "", fmt.Errorf("step1: challenge_id not found in Location: %s", locURL)
   }

   return challengeID, nil
}

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
