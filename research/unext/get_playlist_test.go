package unext

import (
   "net/http"
   "testing"
)

func TestGetPlaylistURL(t *testing.T) {
   // --- Load saved tokens from previous test ---
   tokens, err := LoadTokens("tokens.json")
   if err != nil {
      t.Skipf("no saved tokens found: %v", err)
   }

   t.Logf("loaded access_token: %s...", tokens.AccessToken[:50])

   // --- HTTP client ---
   client := &http.Client{}

   // --- Call GetPlaylistURL with the loaded access token ---
   playlist, err := GetPlaylistURL(client, tokens.AccessToken)
   if err != nil {
      t.Fatalf("GetPlaylistURL: %v", err)
   }

   // --- Get the updated DASH MPD URL with the play_token appended ---
   mpdURL, err := playlist.GetDASHPlaylistURL()
   if err != nil {
      t.Fatalf("GetDASHPlaylistURL: %v", err)
   }

   t.Logf("MPD URL with Token: %s", mpdURL.String())
}
