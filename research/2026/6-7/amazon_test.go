// amazon_auth_test.go
package amazon

import (
   "encoding/json"
   "net/http"
   "net/http/cookiejar"
   "net/url"
   "os"
   "path/filepath"
   "strings"
   "testing"
)

type authState struct {
   Cookies         []*http.Cookie    `json:"cookies"`
   OtpActionUrl    string            `json:"otp_action_url"`
   OtpHiddenParams map[string]string `json:"otp_hidden_params"`
   CodeVerifier    string            `json:"code_verifier"`
}

func getStatePath() string {
   return filepath.Join(os.TempDir(), "amazon_test_state.json")
}

func Test01_TriggerOTP(t *testing.T) {
   jar, _ := cookiejar.New(nil)
   client := &http.Client{
      Jar: jar,
      CheckRedirect: func(req *http.Request, via []*http.Request) error {
         return http.ErrUseLastResponse
      },
   }

   deviceSerial := "ad5e1b330b2d4e5eac8a31dd694bed17"
   authDeviceType := "A1MPSLFC7L5AFK"
   clientId := GenerateClientID(deviceSerial, authDeviceType)

   verifier, challenge, err := GeneratePKCE()
   if err != nil {
      t.Fatalf("Failed to generate PKCE: %v", err)
   }
   t.Logf("Generated PKCE Challenge: %s", challenge)

   phoneNumber, err := GetPhoneNumber()
   if err != nil {
      t.Fatalf("Failed to fetch phone number from credential.exe: %v", err)
   }
   t.Logf("Using phone number: %s", phoneNumber)

   t.Log("Step 1: Initializing Signin...")
   actionUrl, hidden, err := InitSignin(client, clientId, challenge)
   if err != nil {
      t.Fatalf("InitSignin failed: %v", err)
   }

   t.Log("Step 2: Submitting SMS Login...")
   otpRedirectUrl, err := SubmitSMS(client, actionUrl, hidden, phoneNumber)
   if err != nil {
      t.Fatalf("SubmitSMS failed: %v", err)
   }

   t.Log("Step 3: Triggering OTP...")
   otpActionUrl, otpHidden, err := TriggerOTP(client, otpRedirectUrl)
   if err != nil {
      t.Fatalf("TriggerOTP failed: %v", err)
   }

   // Extract cookies to save state
   amazonURL, _ := url.Parse("https://www.amazon.com")
   cookies := jar.Cookies(amazonURL)

   state := authState{
      Cookies:         cookies,
      OtpActionUrl:    otpActionUrl,
      OtpHiddenParams: otpHidden,
      CodeVerifier:    verifier,
   }

   stateBytes, err := json.Marshal(state)
   if err != nil {
      t.Fatalf("Failed to marshal state: %v", err)
   }

   statePath := getStatePath()
   if err := os.WriteFile(statePath, stateBytes, 0644); err != nil {
      t.Fatalf("Failed to write state file to %s: %v", statePath, err)
   }
   t.Logf("State saved to %s", statePath)

   err = os.WriteFile("otp.txt", []byte(""), 0644)
   if err != nil {
      t.Fatalf("Failed to create otp.txt: %v", err)
   }
   t.Log("OTP triggered successfully! Please manually enter the 6-digit code into 'otp.txt' before running Test02.")
}

func Test02_SubmitOTPAndRegister(t *testing.T) {
   statePath := getStatePath()
   stateBytes, err := os.ReadFile(statePath)
   if err != nil {
      t.Fatalf("Failed to read state file (did you run Test01 first?): %v", err)
   }

   var state authState
   if err := json.Unmarshal(stateBytes, &state); err != nil {
      t.Fatalf("Failed to unmarshal state: %v", err)
   }

   jar, _ := cookiejar.New(nil)
   amazonURL, _ := url.Parse("https://www.amazon.com")
   jar.SetCookies(amazonURL, state.Cookies)

   client := &http.Client{
      Jar: jar,
      CheckRedirect: func(req *http.Request, via []*http.Request) error {
         return http.ErrUseLastResponse
      },
   }

   otpBytes, err := os.ReadFile("otp.txt")
   if err != nil {
      t.Fatalf("Failed to read otp.txt. Please create it and add the 6-digit code: %v", err)
   }

   otpCode := strings.TrimSpace(string(otpBytes))
   if len(otpCode) != 6 {
      t.Fatalf("Expected a 6-digit OTP code in otp.txt, got: '%s'", otpCode)
   }

   t.Logf("Read OTP: %s. Submitting...", otpCode)
   claimTokenUrl, err := SubmitOTP(client, state.OtpActionUrl, state.OtpHiddenParams, otpCode)
   if err != nil {
      t.Fatalf("SubmitOTP failed: %v", err)
   }

   t.Log("Exchanging Claim Token for Authorization Code...")
   authCode, err := ExchangeClaimToken(client, claimTokenUrl)
   if err != nil {
      t.Fatalf("ExchangeClaimToken failed: %v", err)
   }
   t.Logf("Obtained Auth Code: %s", authCode)

   t.Log("Registering Device via API...")
   deviceSerial := "ad5e1b330b2d4e5eac8a31dd694bed17"

   resp, err := RegisterDevice(client, authCode, state.CodeVerifier, deviceSerial)
   if err != nil {
      t.Fatalf("RegisterDevice failed: %v", err)
   }

   t.Logf("Success! Access Token: %s...", resp.Response.Success.Tokens.Bearer.AccessToken[:10])
   t.Logf("Success! Refresh Token: %s...", resp.Response.Success.Tokens.Bearer.RefreshToken[:10])

   // Cleaning up the state file, but deliberately leaving otp.txt intact
   _ = os.Remove(statePath)
}
