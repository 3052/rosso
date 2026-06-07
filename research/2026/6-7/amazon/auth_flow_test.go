package amazon

import (
   "crypto/tls"
   "encoding/json"
   "net/http"
   "net/http/cookiejar"
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
   Name   string `json:"name"`
   Value  string `json:"value"`
   Domain string `json:"domain"`
   Path   string `json:"path"`
}

type AuthState struct {
   SessionID    string         `json:"session_id"`
   Cookies      []SimpleCookie `json:"cookies"`
   FormValues   url.Values     `json:"form_values"`
   CodeVerifier string         `json:"code_verifier"`
   CVFReferer   string         `json:"cvf_referer"`
}

type SavedTokens struct {
   AccessToken      string `json:"access_token"`
   RefreshToken     string `json:"refresh_token"`
   DevicePrivateKey string `json:"device_private_key"`
   AdpToken         string `json:"adp_token"`
}

func getTempStatePath() string {
   return filepath.Join(os.TempDir(), "amazon_auth_state.json")
}

func getTempTokensPath() string {
   return filepath.Join(os.TempDir(), "amazon_tokens.json")
}

// createNonBotClient creates an HTTP client that forces HTTP/1.1 to bypass AWS WAF HTTP/2 Fingerprinting
func createNonBotClient(jar *cookiejar.Jar) *http.Client {
   transport := &http.Transport{
      ForceAttemptHTTP2: false,
      TLSNextProto:      make(map[string]func(authority string, c *tls.Conn) http.RoundTripper),
      Proxy:             http.ProxyFromEnvironment,
   }
   return &http.Client{
      Transport: transport,
      Jar:       jar,
      CheckRedirect: func(req *http.Request, via []*http.Request) error {
         return http.ErrUseLastResponse
      },
   }
}

func TestAuthFlow_Part1_RequestOTP(t *testing.T) {
   deviceID := "ad5e1b330b2d4e5eac8a31dd694bed17"

   jar, _ := cookiejar.New(nil)
   client := createNonBotClient(jar)

   // Read static FRC cookie from local file
   frcData, err := os.ReadFile("frc.txt")
   var frcCookie string
   if err == nil {
      frcCookie = strings.TrimSpace(string(frcData))
      t.Log("Loaded FRC cookie from frc.txt")
   } else {
      t.Fatalf("frc.txt not found. Please create it and paste a valid frc cookie.")
   }

   t.Log("--- Executing FetchSignInPage ---")
   formValues, codeVerifier, signInUrlStr, err := FetchSignInPage(client, deviceID, frcCookie)
   if err != nil {
      t.Fatalf("FetchSignInPage failed: %v", err)
   }

   amazonURL, _ := url.Parse("https://www.amazon.com")
   var sessionID string
   for _, cookie := range client.Jar.Cookies(amazonURL) {
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
   formValues.Set("password", "") // Explicitly blank for passwordless SMS login

   redirectURL, err := SubmitCredentials(client, sessionID, formValues, signInUrlStr)
   if err != nil {
      if err.Error() == "CAPTCHA_REQUIRED" {
         t.Fatalf("AMAZON CAPTCHA DETECTED! Saved to captcha_debug.html.")
      }
      t.Fatalf("SubmitCredentials failed: %v", err)
   }

   if !strings.Contains(redirectURL, "/ap/cvf/request") {
      t.Fatalf("Expected redirect to OTP challenge (/ap/cvf/request), but got: %s.", redirectURL)
   }

   t.Log("--- Executing FetchCVFPage (Triggering SMS) ---")
   cvfFormValues, err := FetchCVFPage(client, redirectURL, signInUrlStr)
   if err != nil {
      if err.Error() == "CAPTCHA_REQUIRED" {
         t.Fatalf("AMAZON CAPTCHA DETECTED during CVF fetch! The WAF blocked the request.")
      }
      t.Fatalf("FetchCVFPage failed: %v", err)
   }

   t.Log("--- Saving State to Disk ---")
   state := AuthState{
      SessionID:    sessionID,
      FormValues:   cvfFormValues,
      CodeVerifier: codeVerifier,
      CVFReferer:   redirectURL,
   }

   for _, c := range client.Jar.Cookies(amazonURL) {
      state.Cookies = append(state.Cookies, SimpleCookie{
         Name:   c.Name,
         Value:  c.Value,
         Domain: c.Domain,
         Path:   c.Path,
      })
   }

   stateData, _ := json.MarshalIndent(state, "", "  ")
   os.WriteFile(getTempStatePath(), stateData, 0644)

   t.Log("Successfully requested OTP! Please create 'otp.txt' with your 6-digit code and proceed to Part 2.")
}

func TestAuthFlow_Part2_VerifyOTP(t *testing.T) {
   deviceID := "ad5e1b330b2d4e5eac8a31dd694bed17"

   stateData, err := os.ReadFile(getTempStatePath())
   if err != nil {
      t.Fatalf("Failed to read state file: %v", err)
   }
   var state AuthState
   json.Unmarshal(stateData, &state)

   jar, _ := cookiejar.New(nil)
   amazonURL, _ := url.Parse("https://www.amazon.com")
   var httpCookies []*http.Cookie
   for _, sc := range state.Cookies {
      domain := sc.Domain
      if domain == "" {
         domain = ".amazon.com"
      }
      path := sc.Path
      if path == "" {
         path = "/"
      }
      httpCookies = append(httpCookies, &http.Cookie{
         Name:   sc.Name,
         Value:  sc.Value,
         Domain: domain,
         Path:   path,
      })
   }
   jar.SetCookies(amazonURL, httpCookies)

   client := createNonBotClient(jar)

   otpData, err := os.ReadFile("otp.txt")
   if err != nil {
      t.Fatalf("Failed to read 'otp.txt': %v", err)
   }
   state.FormValues.Set("code", strings.TrimSpace(string(otpData)))

   t.Log("--- Executing VerifyOTP ---")
   claimRedirectURL, err := VerifyOTP(client, state.FormValues)
   if err != nil {
      t.Fatalf("VerifyOTP failed: %v", err)
   }

   if !strings.Contains(claimRedirectURL, "claimToken=") {
      t.Fatalf("Expected redirect URL to contain 'claimToken=', got: %s", claimRedirectURL)
   }

   t.Log("--- Executing FetchClaimSignInPage ---")
   finalRedirectURL, err := FetchClaimSignInPage(client, claimRedirectURL)
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
   accessToken, refreshToken, privateKey, adpToken, err := RegisterDevice(authCode, state.CodeVerifier, deviceID)
   if err != nil {
      t.Fatalf("RegisterDevice failed: %v", err)
   }

   tokens := SavedTokens{
      AccessToken:      accessToken,
      RefreshToken:     refreshToken,
      DevicePrivateKey: privateKey,
      AdpToken:         adpToken,
   }
   tokenData, _ := json.MarshalIndent(tokens, "", "  ")
   tokenPath := getTempTokensPath()
   err = os.WriteFile(tokenPath, tokenData, 0644)
   if err != nil {
      t.Fatalf("Failed to save tokens to disk: %v", err)
   }

   t.Log("========================================")
   t.Log("SUCCESS! WEB AUTHENTICATION COMPLETE!")
   t.Logf("Saved auth data to %s", tokenPath)
   t.Log("You can now run the playback flow tests.")
   t.Log("========================================")
}
