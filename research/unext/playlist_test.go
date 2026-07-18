// playlist_test.go
package unext

import (
   "encoding/json"
   "io"
   "net/http"
   "net/http/cookiejar"
   "testing"
)

func TestGetPlaylist(t *testing.T) {
   // --- Load tokens from file (saved by TestLogin) ---
   tokens, err := LoadTokens(tokensFile)
   if err != nil {
      t.Fatalf("LoadTokens: %v", err)
   }
   t.Logf("loaded access_token (first 20 chars): %s...", tokens.AccessToken[:20])

   // --- HTTP client with cookie jar ---
   jar, err := cookiejar.New(nil)
   if err != nil {
      t.Fatalf("cookiejar.New: %v", err)
   }

   client := &http.Client{
      Jar: jar,
      // Disable auto-redirect so we can capture 302 Location headers if any.
      CheckRedirect: func(req *http.Request, via []*http.Request) error {
         return http.ErrUseLastResponse
      },
   }

   // --- Step 5: GET playlist via GraphQL ---
   playlist, err := Step5GetPlaylist(client, tokens.AccessToken)
   if err != nil {
      t.Fatalf("Step5GetPlaylist: %v", err)
   }

   // --- Assert we got meaningful data ---
   if playlist.PlayToken == "" {
      t.Fatal("playToken is empty")
   }
   if len(playlist.UrlInfo) == 0 {
      t.Fatal("urlInfo is empty")
   }

   // --- Print playlist result ---
   out, _ := json.MarshalIndent(playlist, "", "  ")
   t.Logf("playlist:\n%s", string(out))

   // --- Get the MPD URL ---
   mpdURL, err := playlist.MPDURL()
   if err != nil {
      t.Fatalf("MPDURL: %v", err)
   }
   t.Logf("MPD URL: %s", mpdURL.String())

   // --- Fetch the MPD ---
   req, err := http.NewRequest("GET", mpdURL.String(), nil)
   if err != nil {
      t.Fatalf("creating MPD request: %v", err)
   }

   resp, err := client.Do(req)
   if err != nil {
      t.Fatalf("fetching MPD: %v", err)
   }
   defer resp.Body.Close()

   mpdBody, err := io.ReadAll(resp.Body)
   if err != nil {
      t.Fatalf("reading MPD body: %v", err)
   }

   if resp.StatusCode != http.StatusOK {
      t.Fatalf("MPD request: expected 200, got %d: %s", resp.StatusCode, string(mpdBody))
   }

   t.Logf("MPD response body (%d bytes):\n%s", len(mpdBody), string(mpdBody))
}
