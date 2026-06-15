package amazon

import (
   "encoding/json"
   "os"
   "testing"
)

func TestGetVodPlaybackResources(t *testing.T) {
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

   titleId := "amzn1.dv.gti.28b85d90-1338-720b-4be7-3247683a7624"

   // Calling the updated function which returns a *PlaybackResource
   resource, err := GetVodPlaybackResources(actorState.AccessToken, titleId, envState.PlaybackEnvelope)
   if err != nil {
      t.Fatalf("Failed to get VOD playback resources: %v", err)
   }

   // Accessing the URL property from our struct via dot notation
   t.Logf("Successfully retrieved MPD URL:\n%s", resource.URL)
}
