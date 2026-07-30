// journey_start_test.go
package peacock

import (
   "encoding/json"
   "os"
   "testing"
)

// TestStartJourneyLive makes a real request to the Peacock API to start a journey.
func TestStartJourneyLive(t *testing.T) {
   client := NewClient("")

   resp, err := client.StartJourney()
   if err != nil {
      t.Fatalf("StartJourney failed: %v", err)
   }

   if resp.JourneyID == "" {
      t.Error("expected non-empty JourneyID")
   }
   if len(resp.OneTimePassword) != 6 {
      t.Errorf("expected 6-character OneTimePassword, got %d characters: %q", len(resp.OneTimePassword), resp.OneTimePassword)
   }
   if resp.PollingPeriodSecs <= 0 {
      t.Errorf("expected positive PollingPeriodSecs, got %d", resp.PollingPeriodSecs)
   }
   if resp.RemainingOTPTTLSecs <= 0 {
      t.Errorf("expected positive RemainingOTPTTLSecs, got %d", resp.RemainingOTPTTLSecs)
   }
   if resp.Status != "NOT_STARTED" {
      t.Errorf("expected status NOT_STARTED, got %s", resp.Status)
   }

   t.Logf("DeviceID:           %s", client.DeviceID)
   t.Logf("JourneyID:          %s", resp.JourneyID)
   t.Logf("OneTimePassword:    %s", resp.OneTimePassword)
   t.Logf("PollingPeriodSecs:  %d", resp.PollingPeriodSecs)
   t.Logf("TTL Secs:           %d", resp.RemainingOTPTTLSecs)
   t.Logf("Status:             %s", resp.Status)

   t.Logf(`
====================================
 MANUAL ACTION REQUIRED 
====================================
1. Visit: https://peacocktv.com/activate
2. Enter the code: %s
3. Complete the login within %d seconds
====================================`, resp.OneTimePassword, resp.RemainingOTPTTLSecs)

   state := journeyState{
      DeviceID:            client.DeviceID,
      JourneyID:           resp.JourneyID,
      OneTimePassword:     resp.OneTimePassword,
      PollingPeriodSecs:   resp.PollingPeriodSecs,
      RemainingOTPTTLSecs: resp.RemainingOTPTTLSecs,
      Status:              resp.Status,
   }

   data, err := json.MarshalIndent(state, "", "  ")
   if err != nil {
      t.Fatalf("failed to marshal journey state: %v", err)
   }

   if err := os.WriteFile("journey_state.json", data, 0644); err != nil {
      t.Fatalf("failed to write journey_state.json: %v", err)
   }
}

type journeyState struct {
   DeviceID            string `json:"deviceId"`
   JourneyID           string `json:"journeyId"`
   OneTimePassword     string `json:"oneTimePassword"`
   PollingPeriodSecs   int    `json:"pollingPeriodSecs"`
   RemainingOTPTTLSecs int    `json:"remainingOtpTtlSecs"`
   Status              string `json:"status"`
   HouseholdID         string `json:"householdId,omitempty"`
   Token               string `json:"token,omitempty"`
   UserToken           string `json:"userToken,omitempty"`
   ContentID           string `json:"contentId,omitempty"`
   ProviderVariantID   string `json:"providerVariantId,omitempty"`
}
