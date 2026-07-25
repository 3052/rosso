package peacock

import (
   "net/http"
   "net/http/cookiejar"
)

const (
   BaseID  = "https://rango.id.peacocktv.com"
   BaseSAS = "https://sas.peacocktv.com"
)

// NewClient creates an *http.Client with a cookie jar so that
// Set-Cookie response headers from the verify step are carried
// forward to subsequent requests automatically.
func NewClient() (*http.Client, error) {
   jar, err := cookiejar.New(nil)
   if err != nil {
      return nil, err
   }
   return &http.Client{Jar: jar}, nil
}

// SkyHeaders returns the common x-skyott-* headers used across all
// requests in the Peacock sign-in flow.
func SkyHeaders() http.Header {
   h := http.Header{}
   h.Set("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:140.0) Gecko/20100101 Firefox/140.0")
   h.Set("accept", "application/vnd.siren+json")
   h.Set("accept-language", "en-US,en;q=0.5")
   h.Set("accept-encoding", "gzip, deflate, br, zstd")
   h.Set("referer", "https://www.peacocktv.com/")
   h.Set("origin", "https://www.peacocktv.com")
   h.Set("x-skyott-platform", "PC")
   h.Set("x-skyott-device", "COMPUTER")
   h.Set("x-skyott-provider", "NBCU")
   h.Set("x-skyott-proposition", "NBCUOTT")
   h.Set("x-skyott-territory", "US")
   return h
}
