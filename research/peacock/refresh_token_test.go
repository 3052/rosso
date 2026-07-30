// refresh_token_test.go
package peacock

import (
   "encoding/json"
   "os"
   "testing"
)

// TestRefreshTokenLive makes a real request to the Peacock API to refresh
// an existing user token using mTLS.
// It relies on journey_state.json being generated and updated by the previous tests.
func TestRefreshTokenLive(t *testing.T) {
   data, err := os.ReadFile("journey_state.json")
   if err != nil {
      t.Fatalf("missing prerequisites: could not read journey_state.json: %v", err)
   }

   var state journeyState
   if err := json.Unmarshal(data, &state); err != nil {
      t.Fatalf("failed to unmarshal journey_state.json: %v", err)
   }

   if state.UserToken == "" {
      t.Fatal("journey_state.json is missing UserToken; run TestExchangeTokenLive first")
   }

   client := NewClient(state.DeviceID)

   t.Logf("Refreshing token for DeviceID: %s", client.DeviceID)

   resp, err := client.RefreshToken(state.UserToken)
   if err != nil {
      t.Fatalf("RefreshToken failed: %v", err)
   }

   if resp.UserToken == "" {
      t.Error("expected non-empty UserToken")
   }

   t.Logf("UserToken: %s", resp.UserToken)
   t.Logf("TokenExpiryTime: %s", resp.TokenExpiryTime)
   t.Logf("RecommendedTokenReacquireTime: %s", resp.RecommendedTokenReacquireTime)

   state.UserToken = resp.UserToken

   updatedData, err := json.MarshalIndent(state, "", "  ")
   if err != nil {
      t.Fatalf("failed to marshal updated journey state: %v", err)
   }

   if err := os.WriteFile("journey_state.json", updatedData, 0644); err != nil {
      t.Fatalf("failed to write updated journey_state.json: %v", err)
   }
}
