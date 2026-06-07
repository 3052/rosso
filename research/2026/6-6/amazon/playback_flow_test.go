package amazon

import (
   "encoding/json"
   "os"
   "path/filepath"
   "strings"
   "testing"
)

func getStoredTokensPath() string {
   return filepath.Join(os.TempDir(), "amazon_tokens.json")
}

func TestPlaybackFlow(t *testing.T) {
   deviceID := "ad5e1b330b2d4e5eac8a31dd694bed17"

   // 1. Load the tokens saved from the auth flow
   tokenData, err := os.ReadFile(getStoredTokensPath())
   if err != nil {
      t.Fatalf("Failed to read tokens from disk. Have you successfully run TestAuthFlow_Part2_VerifyOTP? Error: %v", err)
   }

   var tokens SavedTokens
   if err := json.Unmarshal(tokenData, &tokens); err != nil {
      t.Fatalf("Failed to unmarshal tokens: %v", err)
   }

   if tokens.AccessToken == "" {
      t.Fatal("Access token on disk is empty")
   }

   t.Log("Successfully loaded access token from disk")

   // 2. Get Profile ID
   t.Log("--- Executing GetPrimeVideoProfileId ---")
   profileID, err := GetPrimeVideoProfileId(tokens.AccessToken, deviceID)
   if err != nil {
      t.Fatalf("GetPrimeVideoProfileId failed: %v", err)
   }

   // Sanity Check the Profile ID
   if profileID == "" {
      t.Fatal("GetPrimeVideoProfileId returned an empty profile ID")
   }
   if !strings.HasPrefix(profileID, "amzn1.actor.person") {
      t.Fatalf("Expected profile ID to start with 'amzn1.actor.person', but got: %s", profileID)
   }

   t.Log("========================================")
   t.Log("SUCCESS! PROFILE FETCH COMPLETE!")
   t.Logf("Profile ID: %s", profileID)
   t.Log("========================================")
}
