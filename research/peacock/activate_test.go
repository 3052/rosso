package peacock

import (
   "encoding/json"
   "os"
   "testing"
)

// TestActivateLive makes a real request to the Peacock API to activate the device
// and retrieve the OAuth2 token. It relies on journey_state.json being generated
// and updated by the previous tests.
func TestActivateLive(t *testing.T) {
   // Read the state from the previous test
   data, err := os.ReadFile("journey_state.json")
   if err != nil {
      t.Fatalf("missing prerequisites: could not read journey_state.json: %v", err)
   }

   var state journeyState
   if err := json.Unmarshal(data, &state); err != nil {
      t.Fatalf("failed to unmarshal journey_state.json: %v", err)
   }

   if state.JourneyID == "" {
      t.Fatal("journey_state.json is missing JourneyID")
   }

   client := NewClient(state.DeviceID)

   t.Logf("Activating JourneyID: %s", state.JourneyID)

   token, err := client.Activate(state.JourneyID)
   if err != nil {
      t.Fatalf("Activate failed: %v", err)
   }

   if token == "" {
      t.Fatal("expected non-empty OAuth2 token")
   }

   t.Logf("Token: %s", token)
}
