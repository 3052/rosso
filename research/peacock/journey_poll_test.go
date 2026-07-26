package peacock

import (
   "encoding/json"
   "os"
   "testing"
)

// TestPollJourneyLive makes a real request to the Peacock API to poll a journey.
// It expects the user to have already visited peacocktv.com/activate and entered the code.
func TestPollJourneyLive(t *testing.T) {
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

   t.Logf("Polling JourneyID: %s", state.JourneyID)

   // Poll exactly once. We expect the user to have already entered the code.
   status, err := client.PollJourney(state.JourneyID)
   if err != nil {
      t.Fatalf("PollJourney failed: %v", err)
   }

   t.Logf("Status: %s", status.Status)

   if status.Status != JourneyStatusCompleted {
      t.Fatalf("expected status COMPLETED, got %s (did you enter the code in time?)", status.Status)
   }

   if status.HouseholdID == "" {
      t.Fatal("expected non-empty HouseholdID after completion")
   }

   t.Logf("HouseholdID: %s", status.HouseholdID)

   // Update the state file with the completed status and household ID for the next test
   state.Status = status.Status
   state.HouseholdID = status.HouseholdID

   updatedData, err := json.MarshalIndent(state, "", "  ")
   if err != nil {
      t.Fatalf("failed to marshal updated journey state: %v", err)
   }

   if err := os.WriteFile("journey_state.json", updatedData, 0644); err != nil {
      t.Fatalf("failed to write updated journey_state.json: %v", err)
   }
}
