package amazon

import (
   "encoding/json"
   "os"
   "testing"
)

func TestCompleteLogin(t *testing.T) {
   filePath := getStateFilePath()
   data, err := os.ReadFile(filePath)
   if err != nil {
      t.Fatalf("Failed to read state from temp dir (did you run TestInitiateLogin first?): %v", err)
   }

   var state authState
   if err := json.Unmarshal(data, &state); err != nil {
      t.Fatalf("Failed to unmarshal auth state: %v", err)
   }

   accessToken, refreshToken, err := PollRegister(state.PublicCode, state.PrivateCode)
   if err != nil {
      t.Fatalf("Login incomplete or failed: %v", err)
   }

   t.Log("Login successful!")

   tokens := tokenState{
      AccessToken:  accessToken,
      RefreshToken: refreshToken,
   }

   tokenData, err := json.Marshal(tokens)
   if err != nil {
      t.Fatalf("Failed to marshal tokens: %v", err)
   }

   tokensFilePath := getTokensFilePath()
   err = os.WriteFile(tokensFilePath, tokenData, 0600)
   if err != nil {
      t.Fatalf("Failed to write tokens to temp dir: %v", err)
   }

   t.Logf("Tokens saved to: %s", tokensFilePath)

   // Clean up the initial state file after successful retrieval
   _ = os.Remove(filePath)
}
