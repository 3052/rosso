// playout_test.go
package peacock

import (
   "encoding/json"
   "os"
   "testing"
)

// TestPlayoutVodLive makes a real request to the Peacock API to request a VOD
// playout URL using mTLS.
// It relies on journey_state.json being generated and updated by the previous tests.
func TestPlayoutVodLive(t *testing.T) {
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

   t.Logf("Requesting playout for DeviceID: %s", client.DeviceID)

   resp, err := client.PlayoutVod(PlayoutVodParams{
      UserToken:         state.UserToken,
      ContentID:         "GMO_00000000158234_02_HDSDR",
      ProviderVariantID: "1cba422b-3533-33a4-84af-d57cb97bbfa1",
   })
   if err != nil {
      t.Fatalf("PlayoutVod failed: %v", err)
   }

   if len(resp.Asset.Endpoints) == 0 {
      t.Error("expected non-empty Asset.Endpoints")
   }

   t.Logf("LicenceAcquisitionUrl: %s", resp.Protection.LicenceAcquisitionUrl)
   for i, endpoint := range resp.Asset.Endpoints {
      t.Logf("Endpoint[%d]: Cdn=%s Url=%s", i, endpoint.Cdn, endpoint.Url)
   }
}
