// otp_signin.go
package peacock

import (
   "encoding/json"
   "fmt"
   "io"
   "net/http"
   "net/url"
   "strings"
)

const idBase = "https://rango.id.peacocktv.com"

// IdSession holds the idsession cookie obtained via the OTP sign-in flow.
type IdSession struct {
   Cookie *http.Cookie
}

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

// InitiateOTP sends the email to Peacock's identity service to start the OTP
// sign-in journey. It returns the opaque token that must be passed to VerifyOTP.
func (c *Client) InitiateOTP(email string) (*OtpInitiate, error) {
   form := url.Values{}
   form.Set("userIdentifier", email)
   form.Set("journeyContext", "web-signin")

   body := form.Encode()
   req, err := http.NewRequest(http.MethodPost, idBase+"/signin/otp", nil)
   if err != nil {
      return nil, fmt.Errorf("building request: %w", err)
   }
   req.Header = c.idHeaders()
   req.Header.Set("content-type", "application/x-www-form-urlencoded")
   req.Body = io.NopCloser(strings.NewReader(body))
   req.ContentLength = int64(len(body))

   resp, err := c.HTTP.Do(req)
   if err != nil {
      return nil, fmt.Errorf("sending request: %w", err)
   }
   defer resp.Body.Close()

   if resp.StatusCode != http.StatusCreated {
      return nil, fmt.Errorf("unexpected status %d", resp.StatusCode)
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

// VerifyOTP submits the token (from InitiateOTP) along with the 6-digit OTP
// code the user received by email. On success the server sets an idsession
// cookie which is stored in the returned IdSession.
func (c *Client) VerifyOTP(token, otp string) (*IdSession, *OTPVerifyResponse, error) {
   form := url.Values{}
   form.Set("token", token)
   form.Set("otp", otp)

   body := form.Encode()
   req, err := http.NewRequest(http.MethodPost, idBase+"/signin/otp/verify", nil)
   if err != nil {
      return nil, nil, fmt.Errorf("building request: %w", err)
   }
   req.Header = c.idHeaders()
   req.Header.Set("content-type", "application/x-www-form-urlencoded")
   req.Body = io.NopCloser(strings.NewReader(body))
   req.ContentLength = int64(len(body))

   resp, err := c.HTTP.Do(req)
   if err != nil {
      return nil, nil, fmt.Errorf("sending request: %w", err)
   }
   defer resp.Body.Close()

   if resp.StatusCode != http.StatusCreated {
      return nil, nil, fmt.Errorf("unexpected status %d", resp.StatusCode)
   }

   var result OTPVerifyResponse
   respBody, err := io.ReadAll(resp.Body)
   if err != nil {
      return nil, nil, fmt.Errorf("reading body: %w", err)
   }
   if err := json.Unmarshal(respBody, &result); err != nil {
      return nil, nil, fmt.Errorf("parsing JSON: %w", err)
   }

   var sessionCookie *http.Cookie
   for _, cookie := range resp.Cookies() {
      if cookie.Name == "idsession" {
         sessionCookie = cookie
         break
      }
   }
   if sessionCookie == nil {
      return nil, nil, fmt.Errorf("idsession cookie not present")
   }

   return &IdSession{sessionCookie}, &result, nil
}

// idHeaders returns the headers used by the identity/OTP endpoints.
func (c *Client) idHeaders() http.Header {
   h := http.Header{}
   h.Set("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:140.0) Gecko/20100101 Firefox/140.0")
   h.Set("accept", "application/vnd.siren+json")
   h.Set("accept-language", "en-US,en;q=0.5")
   h.Set("referer", "https://www.peacocktv.com/")
   h.Set("origin", "https://www.peacocktv.com")
   h.Set("x-skyott-platform", "PC")
   h.Set("x-skyott-device", "COMPUTER")
   h.Set("x-skyott-provider", "NBCU")
   h.Set("x-skyott-proposition", "NBCUOTT")
   h.Set("x-skyott-territory", "US")
   return h
}
