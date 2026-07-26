// id_session.go
package peacock

import (
   "net/http"
   "strings"
)

const (
   BaseID  = "https://rango.id.peacocktv.com"
   BaseSAS = "https://sas.peacocktv.com"
)

var Territory = "US"

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

// IdSession holds the idsession cookie obtained via the OTP sign-in flow.
type IdSession struct {
   Cookie *http.Cookie
}

// VerifyOTP completes the OTP sign-in by submitting the token (from
// InitiateOTP) and the 6-digit code the user received by email. On success
// the idsession cookie is stored.
func VerifyOTP(token, otp string) (*IdSession, error) {
   _, cookie, err := RequestVerifyOTP(token, otp)
   if err != nil {
      return nil, err
   }
   return &IdSession{cookie}, nil
}

func (*IdSession) CachePath() string {
   return "rosso/peacock/IdSession"
}

// SignInErrors is retained for compatibility but is no longer populated
// by the OTP flow.
type SignInErrors struct {
   Code   string
   Errors struct {
      CategoryErrors []struct {
         Code    string
         Message string
      }
   }
}

func (e *SignInErrors) Error() string {
   if e.Code != "" {
      return e.Code
   }
   var parts []string
   for _, ce := range e.Errors.CategoryErrors {
      parts = append(parts, ce.Code+": "+ce.Message)
   }
   return strings.Join(parts, "; ")
}
