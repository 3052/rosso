package amazon

import (
   "encoding/json"
   "os"
   "testing"
)

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

   // Updated to receive a *Profile
   profile, err := GetPrimaryProfile(tokens.AccessToken)
   if err != nil {
      t.Fatalf("Failed to get primary profile: %v", err)
   }

   // Pass the extracted string to GetActorToken and receive an *ActorToken
   actorToken, err := GetActorToken(tokens.RefreshToken, profile.ProfileID)
   if err != nil {
      t.Fatalf("Failed to get actor token: %v", err)
   }

   t.Log("Successfully retrieved actor token!")

   // Map the properties of the returned structs into your local test struct
   actorState := actorTokenState{
      ActorId:     profile.ProfileID,
      AccessToken: actorToken.Token,
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
