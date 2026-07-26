// signin.go
package peacock

import (
   "encoding/json"
   "fmt"
   "io"
   "log"
   "net/http"
   "net/url"
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

func doRequest(req *http.Request) (*http.Response, error) {
   log.Println(req.Method, req.URL)
   return http.DefaultClient.Do(req)
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

// ---------------------------------------------------------------------------
// POST /signin/otp/verify — submit the OTP code and receive the cookie
// ---------------------------------------------------------------------------

// OTPVerifyResponse is the JSON returned by POST /signin/otp/verify.
type OTPVerifyResponse struct {
   Class      []string `json:"class"`
   Properties struct {
      EventType string `json:"eventType"`
      Data      struct {
         DeviceID string `json:"deviceid"`
      } `json:"data"`
   } `json:"properties"`
}

// RequestVerifyOTP submits the token (from RequestInitiateOTP) along
// with the 6-digit OTP code the user received by email. On success the
// server sets an idsession cookie which is returned alongside the response.
func RequestVerifyOTP(token, otp string) (*OTPVerifyResponse, *http.Cookie, error) {
   form := url.Values{}
   form.Set("token", token)
   form.Set("otp", otp)

   body := form.Encode()
   req, err := http.NewRequest(http.MethodPost, BaseID+"/signin/otp/verify", nil)
   if err != nil {
      return nil, nil, fmt.Errorf("building request: %w", err)
   }
   req.Header = SkyHeaders()
   req.Header.Set("content-type", "application/x-www-form-urlencoded")
   req.Body = io.NopCloser(strings.NewReader(body))
   req.ContentLength = int64(len(body))

   resp, err := doRequest(req)
   if err != nil {
      return nil, nil, fmt.Errorf("sending request: %w", err)
   }
   defer resp.Body.Close()

   respBody, err := io.ReadAll(resp.Body)
   if err != nil {
      return nil, nil, fmt.Errorf("reading body: %w", err)
   }
   if resp.StatusCode != http.StatusCreated {
      return nil, nil, fmt.Errorf("unexpected status %d: %s", resp.StatusCode, respBody)
   }

   var result OTPVerifyResponse
   if err := json.Unmarshal(respBody, &result); err != nil {
      return nil, nil, fmt.Errorf("parsing JSON: %w", err)
   }

   var sessionCookie *http.Cookie
   for _, c := range resp.Cookies() {
      if c.Name == "idsession" {
         sessionCookie = c
         break
      }
   }
   if sessionCookie == nil {
      return nil, nil, fmt.Errorf("idsession cookie not present")
   }
   return &result, sessionCookie, nil
}

// ---------------------------------------------------------------------------
// POST /signin/otp — initiate the OTP email flow
// ---------------------------------------------------------------------------

// OtpInitiate is the JSON returned by POST /signin/otp.
type OtpInitiate struct {
   Class      []string `json:"class"`
   Properties struct {
      EventType string `json:"eventType"`
      Data      struct {
         Type      string `json:"type"`
         Token     string `json:"token"`
         Timestamp string `json:"timestamp"`
      } `json:"data"`
   } `json:"properties"`
}

// RequestInitiateOTP sends the email to Peacock's identity service to
// start the OTP sign-in journey. It returns the opaque token that must
// be passed to RequestVerifyOTP.
func RequestInitiateOTP(email string) (*OtpInitiate, error) {
   form := url.Values{}
   form.Set("userIdentifier", email)
   form.Set("journeyContext", "web-signin")

   body := form.Encode()
   req, err := http.NewRequest(http.MethodPost, BaseID+"/signin/otp", nil)
   if err != nil {
      return nil, fmt.Errorf("building request: %w", err)
   }
   req.Header = SkyHeaders()
   req.Header.Set("content-type", "application/x-www-form-urlencoded")
   req.Body = io.NopCloser(strings.NewReader(body))
   req.ContentLength = int64(len(body))

   resp, err := doRequest(req)
   if err != nil {
      return nil, fmt.Errorf("sending request: %w", err)
   }
   defer resp.Body.Close()

   if resp.StatusCode != http.StatusCreated {
      return nil, fmt.Errorf("unexpected status %v", resp.StatusCode)
   }
   var result OtpInitiate
   if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
      return nil, fmt.Errorf("parsing JSON: %w", err)
   }
   if result.Properties.Data.Token == "" {
      return nil, fmt.Errorf("no token in response")
   }
   return &result, nil
}

func (*OtpInitiate) CachePath() string {
   return "rosso/peacock/OtpInitiate"
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
