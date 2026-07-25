// id_session.go
package peacock

import (
   "net/http"
   "strings"
)

var Territory = "US"

// IdSession holds the idsession cookie obtained via the OTP sign-in flow.
type IdSession struct {
   Cookie *http.Cookie
}

func (*IdSession) CachePath() string {
   return "rosso/peacock/IdSession"
}

// InitiateOTP starts the OTP sign-in journey by sending a one-time passcode
// to the given email address. It returns the opaque token that must be
// passed to VerifyOTP.
func (*IdSession) InitiateOTP(email string) (string, error) {
   return RequestInitiateOTP(email)
}

// VerifyOTP completes the OTP sign-in by submitting the token (from
// InitiateOTP) and the 6-digit code the user received by email. On success
// the idsession cookie is stored.
func (s *IdSession) VerifyOTP(token, otp string) error {
   _, cookie, err := RequestVerifyOTP(token, otp)
   if err != nil {
      return err
   }
   s.Cookie = cookie
   return nil
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
