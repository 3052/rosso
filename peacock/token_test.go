package peacock

import (
   "encoding/json"
   "os"
   "testing"
)

// TestExchangeTokenLive makes a real request to the Peacock API to exchange
// the activation token for a long-lived user token using mTLS.
// It relies on journey_state.json being generated and updated by the previous tests.
func TestExchangeTokenLive(t *testing.T) {
   data, err := os.ReadFile("journey_state.json")
   if err != nil {
      t.Fatalf("missing prerequisites: could not read journey_state.json: %v", err)
   }

   var state journeyState
   if err := json.Unmarshal(data, &state); err != nil {
      t.Fatalf("failed to unmarshal journey_state.json: %v", err)
   }

   if state.Token == "" {
      t.Fatal("journey_state.json is missing Token; run TestActivateLive first")
   }

   client := NewClient(state.DeviceID)

   t.Logf("Exchanging token for DeviceID: %s", client.DeviceID)

   resp, err := client.ExchangeToken(state.Token)
   if err != nil {
      t.Fatalf("ExchangeToken failed: %v", err)
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

// token_test.go
