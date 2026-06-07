package amazon

import (
   "encoding/json"
   "os"
   "path/filepath"
   "testing"
)

func getStoredTokensPath() string {
   return filepath.Join(os.TempDir(), "amazon_tokens.json")
}

func TestPlaybackFlow(t *testing.T) {
   deviceID := "ad5e1b330b2d4e5eac8a31dd694bed17"

   // 1. Load the MAP tokens saved from the auth flow
   tokenData, err := os.ReadFile(getStoredTokensPath())
   if err != nil {
      t.Fatalf("Failed to read tokens from disk. Have you successfully run TestAuthFlow_Part2_VerifyOTP? Error: %v", err)
   }

   var tokens SavedTokens
   if err := json.Unmarshal(tokenData, &tokens); err != nil {
      t.Fatalf("Failed to unmarshal tokens: %v", err)
   }

   if tokens.AdpToken == "" || tokens.DevicePrivateKey == "" {
      t.Fatal("Missing ADP Token or Private Key on disk. You must re-run TestAuthFlow_Part2_VerifyOTP.")
   }

   t.Log("Successfully loaded MAP tokens and RSA key from disk")

   // 2. Exchange MAP tokens for ATV (Video) tokens via ADP signature
   t.Log("--- Executing GetVideoDeviceToken ---")
   atvAccessToken, atvRefreshToken, err := GetVideoDeviceToken(deviceID, tokens.AdpToken, tokens.DevicePrivateKey)
   if err != nil {
      t.Fatalf("GetVideoDeviceToken failed: %v", err)
   }
   if atvAccessToken == "" || atvRefreshToken == "" {
      t.Fatal("GetVideoDeviceToken returned empty tokens")
   }

   t.Log("========================================")
   t.Log("SUCCESS! VIDEO DEVICE TOKEN FETCH COMPLETE!")
   t.Logf("ATV Access Token: %s... (length: %d)", atvAccessToken[:15], len(atvAccessToken))
   t.Logf("ATV Refresh Token: %s... (length: %d)", atvRefreshToken[:15], len(atvRefreshToken))
   t.Log("========================================")
}
