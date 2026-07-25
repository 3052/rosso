package peacock

import (
   "encoding/json"
   "fmt"
   "os"
   "os/exec"
   "path/filepath"
   "testing"
)

// GetFirstEmail returns the username of the first credential entry.
func GetFirstEmail(host string) (string, error) {
   creds, err := GetCredentials(host)
   if err != nil {
      return "", err
   }
   return creds[0].Username, nil
}

// GetCredentials invokes credential.exe with the given host.
func GetCredentials(host string) ([]Credential, error) {
   cmd := exec.Command("credential.exe", fmt.Sprintf("-j=%s", host))
   out, err := cmd.Output()
   if err != nil {
      return nil, fmt.Errorf("running credential.exe: %w", err)
   }

   var creds []Credential
   if err := json.Unmarshal(out, &creds); err != nil {
      return nil, fmt.Errorf("parsing credential.exe output: %w", err)
   }
   if len(creds) == 0 {
      return nil, fmt.Errorf("no credentials returned for host %q", host)
   }
   return creds, nil
}

// TestRequestInitiateOTP hits the live Peacock identity endpoint to
// initiate an OTP sign-in. The email is fetched from credential.exe
// for the host "peacocktv.com". On success the returned token is
// written to token.txt so the verify step can pick it up.
func TestRequestInitiateOTP(t *testing.T) {
   email, err := GetFirstEmail("peacocktv.com")
   if err != nil {
      t.Fatalf("GetFirstEmail: %v", err)
   }
   t.Logf("using email: %s", email)

   client, err := NewClient()
   if err != nil {
      t.Fatalf("NewClient: %v", err)
   }

   token, err := RequestInitiateOTP(client, email)
   if err != nil {
      t.Fatalf("RequestInitiateOTP: %v", err)
   }
   if token == "" {
      t.Fatal("token is empty")
   }
   t.Logf("token received: %s", token)

   if err := os.WriteFile(filepath.Join(".", "token.txt"), []byte(token), 0o600); err != nil {
      t.Fatalf("writing token file: %v", err)
   }
   t.Log("token written to token.txt")
}

// Credential represents one entry returned by credential.exe.
type Credential struct {
   Date     string `json:"date"`
   Host     string `json:"host"`
   Password string `json:"password"`
   Username string `json:"username"`
}
