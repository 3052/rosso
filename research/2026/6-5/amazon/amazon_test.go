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

// SimpleCookie is used to safely serialize cookie state to JSON
type SimpleCookie struct {
   Name  string `json:"name"`
   Value string `json:"value"`
}

// AuthState holds the data we need to persist between Part 1 and Part 2 of the test
type AuthState struct {
   Cookies      []SimpleCookie `json:"cookies"`
   FormValues   url.Values     `json:"form_values"`
   CodeVerifier string         `json:"code_verifier"`
}

// mergeCookies safely overwrites older cookies with newer ones of the same name
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

func TestAuthFlow_Part1_RequestOTP(t *testing.T) {
   // Dummy 32-character hex string for the device ID
   deviceID := "ad5e1b330b2d4e5eac8a31dd694bed17"

   // ==========================================
   // STEP 1: Fetch Sign-In Page
   // ==========================================
   t.Log("--- Executing FetchSignInPage ---")
   formValues, cookies, codeVerifier, err := FetchSignInPage(deviceID)
   if err != nil {
      t.Fatalf("FetchSignInPage failed: %v", err)
   }

   if codeVerifier == "" {
      t.Error("Expected codeVerifier to not be empty")
   }

   // Extract the session-id cookie required for the next request's URL
   var sessionID string
   for _, cookie := range cookies {
      if cookie.Name == "session-id" {
         sessionID = cookie.Value
         break
      }
   }
   if sessionID == "" {
      t.Fatal("Expected 'session-id' cookie to be returned, but it was missing")
   }

   // ==========================================
   // STEP 2: Fetch Credentials from external exe
   // ==========================================
   t.Log("--- Executing credential.exe ---")
   cmd := exec.Command("credential.exe", "-j=amazon.com")
   output, err := cmd.Output()
   if err != nil {
      t.Fatalf("Failed to execute credential.exe: %v", err)
   }

   var creds []Credential
   if err := json.Unmarshal(output, &creds); err != nil {
      t.Fatalf("Failed to parse JSON output from credential.exe: %v\nOutput: %s", err, string(output))
   }

   if len(creds) == 0 {
      t.Fatal("No credentials returned by credential.exe")
   }

   testPhone := creds[0].Username
   t.Logf("Using phone number from credential.exe: %s", testPhone)

   // ==========================================
   // STEP 3: Submit Phone Number (Passwordless)
   // ==========================================
   t.Log("--- Executing SubmitCredentials (SMS Login) ---")
   formValues.Set("email", testPhone)

   redirectURL, newCookies, err := SubmitCredentials(sessionID, formValues, cookies)
   if err != nil {
      t.Fatalf("SubmitCredentials failed: %v", err)
   }

   if !strings.Contains(redirectURL, "/ap/cvf/request") {
      t.Fatalf("Unexpected redirect URL. Expected OTP challenge page ('/ap/cvf/request'), but got: %s", redirectURL)
   }

   // Update our cookie jar with the new cookies returned from the POST
   cookies = mergeCookies(cookies, newCookies)

   // ==========================================
   // STEP 4: Fetch CVF (OTP) Page
   // ==========================================
   t.Log("--- Executing FetchCVFPage (Triggering SMS) ---")
   cvfFormValues, cvfNewCookies, err := FetchCVFPage(redirectURL, cookies)
   if err != nil {
      t.Fatalf("FetchCVFPage failed: %v", err)
   }

   // Update our cookie jar again
   cookies = mergeCookies(cookies, cvfNewCookies)

   // EXPLICIT CHECK: Ensure Amazon actually returned the OTP form
   csrfToken := cvfFormValues.Get("anti-csrftoken-a2z")
   if csrfToken == "" {
      t.Error("Expected 'anti-csrftoken-a2z' in CVF form values, but it was missing (Amazon might have blocked the request)")
   }

   actionVal := cvfFormValues.Get("action")
   if actionVal != "code" && actionVal != "verify" {
      t.Errorf("Expected 'action' field to be 'code' or 'verify', got: '%s'", actionVal)
   }

   // ==========================================
   // SAVE STATE TO DISK
   // ==========================================
   t.Log("--- Saving State to Disk ---")
   state := AuthState{
      FormValues:   cvfFormValues,
      CodeVerifier: codeVerifier,
   }

   for _, c := range cookies {
      state.Cookies = append(state.Cookies, SimpleCookie{Name: c.Name, Value: c.Value})
   }

   stateData, err := json.MarshalIndent(state, "", "  ")
   if err != nil {
      t.Fatalf("Failed to marshal auth state: %v", err)
   }

   statePath := getTempStatePath()
   if err := os.WriteFile(statePath, stateData, 0644); err != nil {
      t.Fatalf("Failed to write state file to %s: %v", statePath, err)
   }

   t.Logf("Successfully requested OTP! State saved to %s", statePath)
   t.Log("Please retrieve the SMS code and proceed to TestAuthFlow_Part2_VerifyOTP")
}

func TestAuthFlow_Part2_VerifyOTP(t *testing.T) {
   statePath := getTempStatePath()
   stateData, err := os.ReadFile(statePath)
   if err != nil {
      t.Fatalf("Failed to read state file at %s (Did you run Part 1 first?): %v", statePath, err)
   }

   var state AuthState
   if err := json.Unmarshal(stateData, &state); err != nil {
      t.Fatalf("Failed to unmarshal state data: %v", err)
   }

   // Reconstruct the http.Cookies
   var cookies []*http.Cookie
   for _, sc := range state.Cookies {
      cookies = append(cookies, &http.Cookie{Name: sc.Name, Value: sc.Value})
   }

   t.Log("--- Successfully Loaded State ---")
   t.Logf("Code Verifier: %s", state.CodeVerifier)
   t.Logf("Form fields loaded: %d", len(state.FormValues))
   t.Logf("Cookies loaded: %d", len(cookies))

   // Next step: Prompt for or read the OTP, insert it into state.FormValues, and call VerifyOTP
}
