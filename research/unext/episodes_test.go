// episodes_test.go
package unext

import (
   "net/http"
   "net/http/cookiejar"
   "testing"
)

func TestGetEpisodeCodesViaDetail(t *testing.T) {
   // --- Load tokens from file (saved by TestLogin) ---
   tokens, err := LoadTokens(tokensFile)
   if err != nil {
      t.Fatalf("LoadTokens: %v", err)
   }

   // --- HTTP client with cookie jar ---
   jar, err := cookiejar.New(nil)
   if err != nil {
      t.Fatalf("cookiejar.New: %v", err)
   }

   client := &http.Client{
      Jar: jar,
      CheckRedirect: func(req *http.Request, via []*http.Request) error {
         return http.ErrUseLastResponse
      },
   }

   // --- Fetch episode codes for SID0020149 ---
   codes, err := GetEpisodeCodesViaDetail(client, tokens.AccessToken, "SID0020149")
   if err != nil {
      t.Fatalf("GetEpisodeCodesViaDetail: %v", err)
   }

   if len(codes) == 0 {
      t.Fatal("no episode codes returned")
   }

   for i, code := range codes {
      t.Logf("episode %d: %s", i, code)
   }

   // --- Verify the expected episode code is present ---
   found := false
   for _, code := range codes {
      if code == "ED00092859" {
         found = true
         break
      }
   }
   if !found {
      t.Fatal("expected episode ED00092859 not found")
   }
}
