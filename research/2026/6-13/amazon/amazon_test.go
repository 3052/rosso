package amazon

import (
   "encoding/json"
   "os"
   "path/filepath"
   "testing"
)

type authState struct {
   PublicCode  string `json:"public_code"`
   PrivateCode string `json:"private_code"`
}

type tokenState struct {
   AccessToken  string `json:"access_token"`
   RefreshToken string `json:"refresh_token"`
}

type actorTokenState struct {
   ActorId     string `json:"actor_id"`
   AccessToken string `json:"access_token"`
}

func getStateFilePath() string {
   return filepath.Join(os.TempDir(), "amazon_auth_state.json")
}

func getTokensFilePath() string {
   return filepath.Join(os.TempDir(), "amazon_tokens.json")
}

func getActorTokensFilePath() string {
   return filepath.Join(os.TempDir(), "amazon_actor_tokens.json")
}

func TestInitiateLogin(t *testing.T) {
   publicCode, privateCode, err := CreateCodePair()
   if err != nil {
      t.Fatalf("Failed to create code pair: %v", err)
   }

   err = InitiateMDSO(publicCode)
   if err != nil {
      t.Fatalf("Failed to initiate MDSO: %v", err)
   }

   t.Logf("\n=== AMAZON LOGIN ===\nPlease navigate to https://www.amazon.com/us/code\nEnter the following code: %s\n====================\n", publicCode)

   state := authState{
      PublicCode:  publicCode,
      PrivateCode: privateCode,
   }

   data, err := json.Marshal(state)
   if err != nil {
      t.Fatalf("Failed to marshal auth state: %v", err)
   }

   filePath := getStateFilePath()
   err = os.WriteFile(filePath, data, 0600)
   if err != nil {
      t.Fatalf("Failed to write state to temp dir: %v", err)
   }

   t.Logf("State saved to: %s", filePath)
}

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

func TestGetActorToken(t *testing.T) {
   tokensFilePath := getTokensFilePath()
   data, err := os.ReadFile(tokensFilePath)
   if err != nil {
      t.Fatalf("Failed to read tokens from temp dir (did you run TestCompleteLogin first?): %v", err)
   }

   var tokens tokenState
   if err := json.Unmarshal(data, &tokens); err != nil {
      t.Fatalf("Failed to unmarshal tokens: %v", err)
   }

   actorId, err := GetPrimaryProfile(tokens.AccessToken)
   if err != nil {
      t.Fatalf("Failed to get primary profile: %v", err)
   }

   actorAccessToken, err := GetActorToken(tokens.RefreshToken, actorId)
   if err != nil {
      t.Fatalf("Failed to get actor token: %v", err)
   }

   t.Log("Successfully retrieved actor token!")

   actorState := actorTokenState{
      ActorId:     actorId,
      AccessToken: actorAccessToken,
   }

   actorData, err := json.Marshal(actorState)
   if err != nil {
      t.Fatalf("Failed to marshal actor state: %v", err)
   }

   actorFilePath := getActorTokensFilePath()
   err = os.WriteFile(actorFilePath, actorData, 0600)
   if err != nil {
      t.Fatalf("Failed to write actor state to temp dir: %v", err)
   }

   t.Logf("Actor state saved to: %s", actorFilePath)
}
