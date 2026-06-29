// signin_test.go
package molotov

import (
   "encoding/json"
   "os"
   "os/exec"
   "testing"
)

func TestSignin(t *testing.T) {
   // Fetch credentials using the local credential executable
   cmd := exec.Command("credential.exe", "-j", "molotov.tv")
   output, err := cmd.Output()
   if err != nil {
      t.Fatalf("Failed to execute credential.exe: %v", err)
   }

   var creds []CredentialItem
   if err := json.Unmarshal(output, &creds); err != nil {
      t.Fatalf("Failed to parse credentials JSON: %v\nOutput: %s", err, string(output))
   }

   if len(creds) == 0 {
      t.Fatal("No credentials returned by credential.exe")
   }

   username := creds[0].Username
   password := creds[0].Password

   token, err := Signin(username, password)
   if err != nil {
      t.Fatalf("Signin failed: %v", err)
   }

   if token == "" {
      t.Fatal("Expected an access token, but got an empty string")
   }

   // Prepare the data to be written for future tests (e.g., GetUser, GetAsset)
   authData := TestAuthData{
      AccessToken: token,
   }

   data, err := json.MarshalIndent(authData, "", "  ")
   if err != nil {
      t.Fatalf("Failed to marshal auth data: %v", err)
   }

   // Write the resulting token to a file
   err = os.WriteFile("auth_test.json", data, 0600)
   if err != nil {
      t.Fatalf("Failed to write auth data to file: %v", err)
   }

   t.Log("Successfully signed in and saved credentials to auth_test.json")
}

type CredentialItem struct {
   Username string `json:"username"`
   Password string `json:"password"`
}

// TestAuthData is used to serialize the authentication result for future tests.
type TestAuthData struct {
   AccessToken string `json:"access_token"`
}
