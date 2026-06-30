// user_test.go
package molotov

import (
   "encoding/json"
   "os"
   "testing"
)

func TestGetUser(t *testing.T) {
   authFile, err := os.ReadFile("auth_test.json")
   if err != nil {
      t.Fatalf("Failed to read auth_test.json. Please run TestSignin first: %v", err)
   }

   var signinResp SigninResponse
   if err := json.Unmarshal(authFile, &signinResp); err != nil {
      t.Fatalf("Failed to unmarshal SigninResponse: %v", err)
   }

   userResp, err := GetUser(&signinResp)
   if err != nil {
      t.Fatalf("GetUser failed: %v", err)
   }

   // Accessing the unwrapped fields directly
   if userResp.ID == "" || userResp.Profiles[0].ID == "" {
      t.Fatalf("Expected non-empty user and profile IDs")
   }

   data, err := json.MarshalIndent(userResp, "", "  ")
   if err != nil {
      t.Fatalf("Failed to marshal user data: %v", err)
   }

   err = os.WriteFile("user_test.json", data, 0600)
   if err != nil {
      t.Fatalf("Failed to write user data to file: %v", err)
   }

   t.Logf("Successfully retrieved UserResponse and saved to user_test.json")
}
