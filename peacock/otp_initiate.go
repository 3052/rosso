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

// RequestInitiateOTP sends the email to Peacock's identity service to
// start the OTP sign-in journey. It returns the opaque token that must
// be passed to RequestVerifyOTP.
func RequestInitiateOTP(email string) (string, error) {
   form := url.Values{}
   form.Set("userIdentifier", email)
   form.Set("journeyContext", "web-signin")

   body := form.Encode()
   req, err := http.NewRequest(http.MethodPost, BaseID+"/signin/otp", nil)
   if err != nil {
      return "", fmt.Errorf("building request: %w", err)
   }
   req.Header = SkyHeaders()
   req.Header.Set("content-type", "application/x-www-form-urlencoded")
   req.Body = io.NopCloser(strings.NewReader(body))
   req.ContentLength = int64(len(body))

   resp, err := doRequest(req)
   if err != nil {
      return "", fmt.Errorf("sending request: %w", err)
   }
   defer resp.Body.Close()

   respBody, err := io.ReadAll(resp.Body)
   if err != nil {
      return "", fmt.Errorf("reading body: %w", err)
   }
   if resp.StatusCode != http.StatusCreated {
      return "", fmt.Errorf("unexpected status %d: %s", resp.StatusCode, respBody)
   }

   var result OTPInitiateResponse
   if err := json.Unmarshal(respBody, &result); err != nil {
      return "", fmt.Errorf("parsing JSON: %w", err)
   }

   if result.Properties.Data.Token == "" {
      return "", fmt.Errorf("no token in response: %s", respBody)
   }
   return result.Properties.Data.Token, nil
}

// OTPInitiateResponse is the JSON returned by POST /signin/otp.
type OTPInitiateResponse struct {
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
