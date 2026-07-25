package peacock

import (
   "os"
   "strings"
   "testing"
)

// TestRequestVerifyOTP reads the token from token.txt (written by
// TestRequestInitiateOTP) and the OTP code from otp.txt (which you
// must create after receiving the email), then calls the verify
// endpoint.
//
// Run order:
//  1. go test -v -run TestRequestInitiateOTP   # writes token.txt
//  2. <check your email, write the 6-digit code to otp.txt>
//  3. go test -v -run TestRequestVerifyOTP
func TestRequestVerifyOTP(t *testing.T) {
   tokenBytes, err := os.ReadFile("token.txt")
   if err != nil {
      t.Fatalf("reading token.txt: %v", err)
   }
   token := string(tokenBytes)

   otpBytes, err := os.ReadFile("otp.txt")
   if err != nil {
      t.Fatalf("reading otp.txt: %v (create this file with the 6-digit code from your email)", err)
   }
   otp := string(otpBytes)

   // trim any stray whitespace/newlines from the OTP
   otp = strings.TrimSpace(otp)
   if len(otp) != 6 {
      t.Fatalf("otp in otp.txt doesn't look like a 6-digit code: %q", otp)
   }

   client, err := NewClient()
   if err != nil {
      t.Fatalf("NewClient: %v", err)
   }

   result, err := RequestVerifyOTP(client, token, otp)
   if err != nil {
      t.Fatalf("RequestVerifyOTP: %v", err)
   }

   t.Logf("verify success — deviceid: %s", result.Properties.Data.DeviceID)
}
