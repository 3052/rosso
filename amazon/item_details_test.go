package amazon

import (
   "encoding/json"
   "os"
   "testing"
)

func TestGetItemDetails(t *testing.T) {
   actorFilePath := getActorTokensFilePath()
   data, err := os.ReadFile(actorFilePath)
   if err != nil {
      t.Fatalf("Failed to read actor tokens from temp dir (did you run TestGetActorToken first?): %v", err)
   }

   var actorState actorTokenState
   if err := json.Unmarshal(data, &actorState); err != nil {
      t.Fatalf("Failed to unmarshal actor state: %v", err)
   }

   // Using the title ID from the provided dump
   titleId := "amzn1.dv.gti.28b85d90-1338-720b-4be7-3247683a7624"

   playbackEnvelope, err := GetItemDetails(actorState.AccessToken, titleId)
   if err != nil {
      t.Fatalf("Failed to get item details (playback envelope): %v", err)
   }

   t.Log("Successfully retrieved playback envelope!")

   envState := envelopeState{
      PlaybackEnvelope: playbackEnvelope,
   }

   envData, err := json.Marshal(envState)
   if err != nil {
      t.Fatalf("Failed to marshal envelope state: %v", err)
   }

   envFilePath := getEnvelopeFilePath()
   err = os.WriteFile(envFilePath, envData, 0600)
   if err != nil {
      t.Fatalf("Failed to write envelope state to temp dir: %v", err)
   }

   t.Logf("Playback envelope saved to: %s", envFilePath)
}
