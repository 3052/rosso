// otp_verify.go
package peacock

import (
   "encoding/json"
   "fmt"
   "io"
   "net/http"
   "net/url"
   "strings"
)

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
