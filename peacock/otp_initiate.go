// otp_initiate.go
package peacock

import (
   "encoding/json"
   "fmt"
   "io"
   "net/http"
   "net/url"
   "strings"
)

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
