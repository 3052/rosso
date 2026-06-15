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

type envelopeState struct {
   PlaybackEnvelope string `json:"playback_envelope"`
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

func getEnvelopeFilePath() string {
   return filepath.Join(os.TempDir(), "amazon_envelope.json")
}

func TestInitiateLogin(t *testing.T) {
   // Call the updated function which now returns a *CodePair
   codes, err := CreateCodePair()
   if err != nil {
      t.Fatalf("Failed to create code pair: %v", err)
   }

   // Access the properties using dot notation
   err = InitiateMDSO(codes.PublicCode)
   if err != nil {
      t.Fatalf("Failed to initiate MDSO: %v", err)
   }

   t.Logf("\n=== AMAZON LOGIN ===\nPlease navigate to https://www.amazon.com/us/code\nEnter the following code: %s\n====================\n", codes.PublicCode)

   state := authState{
      PublicCode:  codes.PublicCode,
      PrivateCode: codes.PrivateCode,
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
