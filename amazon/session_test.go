package amazon

import (
   "encoding/json"
   "os"
   "testing"
)

func TestStartSession(t *testing.T) {
   actorFilePath := getActorTokensFilePath()
   actorData, err := os.ReadFile(actorFilePath)
   if err != nil {
      t.Fatalf("Failed to read actor tokens from temp dir (did you run TestGetActorToken first?): %v", err)
   }

   var actorState actorTokenState
   if err := json.Unmarshal(actorData, &actorState); err != nil {
      t.Fatalf("Failed to unmarshal actor state: %v", err)
   }

   envFilePath := getEnvelopeFilePath()
   envData, err := os.ReadFile(envFilePath)
   if err != nil {
      t.Fatalf("Failed to read envelope state from temp dir (did you run TestGetItemDetails first?): %v", err)
   }

   var envState envelopeState
   if err := json.Unmarshal(envData, &envState); err != nil {
      t.Fatalf("Failed to unmarshal envelope state: %v", err)
   }

   err = StartSession(actorState.AccessToken, envState.PlaybackEnvelope)
   if err != nil {
      t.Fatalf("Failed to start playback session: %v", err)
   }

   t.Log("Successfully started playback session!")
}
