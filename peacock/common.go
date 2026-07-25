// common.go
package peacock

import "net/http"

const (
   BaseID  = "https://rango.id.peacocktv.com"
   BaseSAS = "https://sas.peacocktv.com"
)

// SkyHeaders returns the common headers used across all requests.
func SkyHeaders() http.Header {
   h := http.Header{}
   h.Set("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:140.0) Gecko/20100101 Firefox/140.0")
   h.Set("accept", "application/vnd.siren+json")
   h.Set("accept-language", "en-US,en;q=0.5")
   // Deliberately omit "accept-encoding": Go's transport adds gzip
   // automatically and transparently decompresses it for us.
   h.Set("referer", "https://www.peacocktv.com/")
   h.Set("origin", "https://www.peacocktv.com")
   h.Set("x-skyott-platform", "PC")
   h.Set("x-skyott-device", "COMPUTER")
   h.Set("x-skyott-provider", "NBCU")
   h.Set("x-skyott-proposition", "NBCUOTT")
   h.Set("x-skyott-territory", "US")
   return h
}
