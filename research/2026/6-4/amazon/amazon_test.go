// --- amazon_test.go ---
package amazon

import (
   "fmt"
   "testing"
)

// TestActualAmazonRequests makes REAL network requests to the Amazon endpoints.
// Note: The Bearer token provided here is taken from your trace.
// Since these tokens expire, you will need to replace it with a fresh, valid token
// when running this test, otherwise you will receive an HTTP 401/403 Unauthorized error.
func TestActualAmazonRequests(t *testing.T) {
   // Values extracted from the provided trace
   titleID := "amzn1.dv.gti.28b85d90-1338-720b-4be7-3247683a7624"
   deviceID := "7f93fd3fd18646658ccd555423c8e5b8"

   // Replace this with a fresh token before running
   bearerToken := "Atna|EwMDIJzl6tGByLJxk7iiHnOJIinXrQnKk9Bigdbl69Og1AD0-KI_RL6SuPL3In83ZBEwmm2sNIC6HywQ1JdCGuKUVPKL7TbRIQlc2u0J-F_xmmPkJ84cOfptXU7IHJEVYTmyod1ll1am7D6Y09JgIHrBxsNwt7BfIWt4EyNJJgHsCmHqFrlPrO3RI_qFfvhFy0UhJTioPPhw5PtnVFOu_q8gNpkgdgjgdvLGXrgNWXbqLfwn6ay4GNz89KyNqyPFgmUcE7s-PMdv3cQSbftP-flH60WtUqsjVfTKhCSYsdP8CBj1RsMZGmYaMn5_hvlvBn9ke7w"

   fmt.Println("--------------------------------------------------")
   fmt.Println("1. Fetching Playback Envelope...")
   envelope, err := GetPlaybackEnvelope(titleID, deviceID, bearerToken)
   if err != nil {
      t.Fatalf("Failed to get playback envelope: %v\n(Check if your Bearer token is still valid)", err)
   }

   if envelope == "" {
      t.Fatalf("Envelope returned empty")
   }

   fmt.Println("Success! Envelope extracted:")
   fmt.Printf("%s...\n", envelope[:50]) // Print just the beginning to avoid terminal spam

   fmt.Println("--------------------------------------------------")
   fmt.Println("2. Fetching MPD URL...")
   mpdURL, err := GetMPDUrl(titleID, deviceID, bearerToken, envelope)
   if err != nil {
      t.Fatalf("Failed to get MPD URL: %v", err)
   }

   if mpdURL == "" {
      t.Fatalf("MPD URL returned empty")
   }

   fmt.Println("Success! MPD URL extracted:")
   fmt.Println(mpdURL)
   fmt.Println("--------------------------------------------------")
}
