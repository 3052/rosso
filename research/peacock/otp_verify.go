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
// server sets session cookies (idsession, skyCEsidmesso01, etc.) which
// are stored in the client's cookie jar for subsequent requests.
func RequestVerifyOTP(client *http.Client, token, otp string) (*OTPVerifyResponse, error) {
   form := url.Values{}
   form.Set("token", token)
   form.Set("otp", otp)

   req, err := http.NewRequest(http.MethodPost, BaseID+"/signin/otp/verify", nil)
   if err != nil {
      return nil, fmt.Errorf("building request: %w", err)
   }
   req.Header = SkyHeaders()
   req.Header.Set("content-type", "application/x-www-form-urlencoded")
   req.Body = io.NopCloser(strings.NewReader(form.Encode()))
   req.ContentLength = int64(len(form.Encode()))

   resp, err := client.Do(req)
   if err != nil {
      return nil, fmt.Errorf("sending request: %w", err)
   }
   defer resp.Body.Close()

   body, err := io.ReadAll(resp.Body)
   if err != nil {
      return nil, fmt.Errorf("reading body: %w", err)
   }
   if resp.StatusCode != http.StatusCreated {
      return nil, fmt.Errorf("unexpected status %d: %s", resp.StatusCode, body)
   }

   var result OTPVerifyResponse
   if err := json.Unmarshal(body, &result); err != nil {
      return nil, fmt.Errorf("parsing JSON: %w", err)
   }
   return &result, nil
}
