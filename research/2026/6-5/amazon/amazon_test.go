package amazon

import (
   "encoding/json"
   "os/exec"
   "testing"
)

type Credential struct {
   Date     string `json:"date"`
   Host     string `json:"host"`
   Password string `json:"password"`
   Trial    string `json:"trial"`
   Username string `json:"username"`
}

func TestAuthFlow(t *testing.T) {
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
   } else {
      t.Logf("Generated codeVerifier: %s", codeVerifier)
   }

   // Verify that the critical hidden fields were successfully extracted
   csrfToken := formValues.Get("anti-csrftoken-a2z")
   if csrfToken == "" {
      t.Error("Expected 'anti-csrftoken-a2z' to be populated in form values")
   } else {
      t.Logf("Extracted anti-csrftoken-a2z: %s", csrfToken)
   }

   appActionToken := formValues.Get("appActionToken")
   if appActionToken == "" {
      t.Error("Expected 'appActionToken' to be populated in form values")
   } else {
      t.Logf("Extracted appActionToken: %s", appActionToken)
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
   t.Logf("Extracted session-id: %s", sessionID)

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

   redirectURL, updatedCookies, err := SubmitCredentials(sessionID, formValues, cookies)
   if err != nil {
      t.Fatalf("SubmitCredentials failed: %v\n(Did you use a valid phone number? Invalid numbers return 200 OK instead of the expected 302 redirect)", err)
   }

   if redirectURL == "" {
      t.Error("Expected a redirect URL to the CVF (OTP) page, got empty string")
   } else {
      t.Logf("Success! Redirect URL for OTP step: %s", redirectURL)
   }

   if len(updatedCookies) > 0 {
      t.Logf("Received %d new/updated cookies", len(updatedCookies))
      for _, cookie := range updatedCookies {
         t.Logf("  %s: %s", cookie.Name, cookie.Value)
      }
   }
}
