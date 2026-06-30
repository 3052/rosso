// signin_test.go
package molotov

import (
   "encoding/json"
   "os"
   "os/exec"
   "testing"
)

func TestSignin(t *testing.T) {
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

   resp, err := Signin(creds[0].Username, creds[0].Password)
   if err != nil {
      t.Fatalf("Signin failed: %v", err)
   }

   // Accessing the unwrapped field directly
   if resp.AccessToken == "" {
      t.Fatal("Expected an access token in the response, but got an empty string")
   }

   // Directly save the SigninResponse struct for the next test
   data, err := json.MarshalIndent(resp, "", "  ")
   if err != nil {
      t.Fatalf("Failed to marshal auth data: %v", err)
   }

   err = os.WriteFile("auth_test.json", data, 0600)
   if err != nil {
      t.Fatalf("Failed to write auth data to file: %v", err)
   }

   t.Log("Successfully signed in and saved SigninResponse to auth_test.json")
}

type CredentialItem struct {
   Username string `json:"username"`
   Password string `json:"password"`
}
