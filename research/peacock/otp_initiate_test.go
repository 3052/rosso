package peacock

import (
   "os"
   "path/filepath"
   "testing"
)

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
