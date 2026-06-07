package amazon

import (
   "encoding/json"
   "net/http"
   "net/url"
   "os"
   "os/exec"
   "path/filepath"
   "strings"
   "testing"
)

type Credential struct {
   Date     string `json:"date"`
   Host     string `json:"host"`
   Password string `json:"password"`
   Trial    string `json:"trial"`
   Username string `json:"username"`
}

type SimpleCookie struct {
   Name  string `json:"name"`
   Value string `json:"value"`
}

type AuthState struct {
   SessionID    string         `json:"session_id"`
   Cookies      []SimpleCookie `json:"cookies"`
   FormValues   url.Values     `json:"form_values"`
   CodeVerifier string         `json:"code_verifier"`
}

type SavedTokens struct {
   AccessToken  string `json:"access_token"`
   RefreshToken string `json:"refresh_token"`
}

func mergeCookies(existing []*http.Cookie, newCookies []*http.Cookie) []*http.Cookie {
   cookieMap := make(map[string]*http.Cookie)
   for _, c := range existing {
      cookieMap[c.Name] = c
   }
   for _, c := range newCookies {
      cookieMap[c.Name] = c
   }
   var merged []*http.Cookie
   for _, c := range cookieMap {
      merged = append(merged, c)
   }
   return merged
}

func getTempStatePath() string {
   return filepath.Join(os.TempDir(), "amazon_auth_state.json")
}

func getTempTokensPath() string {
   return filepath.Join(os.TempDir(), "amazon_tokens.json")
}

func TestAuthFlow_Part1_RequestOTP(t *testing.T) {
   deviceID := "ad5e1b330b2d4e5eac8a31dd694bed17"

   t.Log("--- Executing FetchSignInPage ---")
   formValues, cookies, codeVerifier, err := FetchSignInPage(deviceID)
   if err != nil {
      t.Fatalf("FetchSignInPage failed: %v", err)
   }

   var sessionID string
   for _, cookie := range cookies {
      if cookie.Name == "session-id" {
         sessionID = cookie.Value
         break
      }
   }

   t.Log("--- Executing credential.exe ---")
   cmd := exec.Command("credential.exe", "-j=amazon.com")
   output, err := cmd.Output()
   if err != nil {
      t.Fatalf("Failed to execute credential.exe: %v", err)
   }

   var creds []Credential
   if err := json.Unmarshal(output, &creds); err != nil {
      t.Fatalf("Failed to parse JSON output: %v", err)
   }
   testPhone := creds[0].Username

   t.Log("--- Executing SubmitCredentials (SMS Login) ---")
   formValues.Set("email", testPhone)

   redirectURL, newCookies, err := SubmitCredentials(sessionID, formValues, cookies)
   if err != nil {
      t.Fatalf("SubmitCredentials failed: %v", err)
   }
   cookies = mergeCookies(cookies, newCookies)

   t.Log("--- Executing FetchCVFPage (Triggering SMS) ---")
   cvfFormValues, cvfNewCookies, err := FetchCVFPage(redirectURL, cookies)
   if err != nil {
      t.Fatalf("FetchCVFPage failed: %v", err)
   }
   cookies = mergeCookies(cookies, cvfNewCookies)

   t.Log("--- Saving State to Disk ---")
   state := AuthState{
      SessionID:    sessionID,
      FormValues:   cvfFormValues,
      CodeVerifier: codeVerifier,
   }
   for _, c := range cookies {
      state.Cookies = append(state.Cookies, SimpleCookie{Name: c.Name, Value: c.Value})
   }
   stateData, _ := json.MarshalIndent(state, "", "  ")
   os.WriteFile(getTempStatePath(), stateData, 0644)

   t.Log("Successfully requested OTP! Please create 'otp.txt' and proceed to Part 2.")
}

func TestAuthFlow_Part2_VerifyOTP(t *testing.T) {
   deviceID := "ad5e1b330b2d4e5eac8a31dd694bed17"

   stateData, err := os.ReadFile(getTempStatePath())
   if err != nil {
      t.Fatalf("Failed to read state file: %v", err)
   }
   var state AuthState
   json.Unmarshal(stateData, &state)

   var cookies []*http.Cookie
   for _, sc := range state.Cookies {
      cookies = append(cookies, &http.Cookie{Name: sc.Name, Value: sc.Value})
   }

   otpData, err := os.ReadFile("otp.txt")
   if err != nil {
      t.Fatalf("Failed to read 'otp.txt': %v", err)
   }
   state.FormValues.Set("code", strings.TrimSpace(string(otpData)))

   t.Log("--- Executing VerifyOTP ---")
   claimRedirectURL, newCookies, err := VerifyOTP(state.FormValues, cookies)
   if err != nil {
      t.Fatalf("VerifyOTP failed: %v", err)
   }
   cookies = mergeCookies(cookies, newCookies)

   if !strings.Contains(claimRedirectURL, "claimToken=") {
      t.Fatalf("Expected redirect URL to contain 'claimToken=', got: %s", claimRedirectURL)
   }

   t.Log("--- Executing FetchClaimSignInPage ---")
   finalRedirectURL, _, err := FetchClaimSignInPage(claimRedirectURL, cookies)
   if err != nil {
      t.Fatalf("FetchClaimSignInPage failed: %v", err)
   }

   parsedUrl, err := url.Parse(finalRedirectURL)
   if err != nil {
      t.Fatalf("Failed to parse final redirect URL: %v", err)
   }

   authCode := parsedUrl.Query().Get("openid.oa2.authorization_code")
   if authCode == "" {
      t.Fatalf("Expected 'openid.oa2.authorization_code' in final redirect URL, got: %s", finalRedirectURL)
   }

   t.Log("--- Executing RegisterDevice ---")
   accessToken, refreshToken, err := RegisterDevice(authCode, state.CodeVerifier, deviceID)
   if err != nil {
      t.Fatalf("RegisterDevice failed: %v", err)
   }

   if accessToken == "" || refreshToken == "" {
      t.Fatal("RegisterDevice returned empty tokens")
   }

   tokens := SavedTokens{
      AccessToken:  accessToken,
      RefreshToken: refreshToken,
   }
   tokenData, _ := json.MarshalIndent(tokens, "", "  ")
   tokenPath := getTempTokensPath()
   err = os.WriteFile(tokenPath, tokenData, 0644)
   if err != nil {
      t.Fatalf("Failed to save tokens to disk: %v", err)
   }

   t.Log("========================================")
   t.Log("SUCCESS! WEB AUTHENTICATION COMPLETE!")
   t.Logf("Saved access and refresh tokens to %s", tokenPath)
   t.Log("You can now run the playback flow tests.")
   t.Log("========================================")
}
