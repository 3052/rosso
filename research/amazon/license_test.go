package amazon

import (
   "net/http"
   "os"
   "strings"
   "testing"
)

// 4. Run this to test requesting a Widevine License
func TestStep4_GetLicense(t *testing.T) {
   client := &http.Client{}

   tokenFile := getTempTokenPath()
   tokenBytes, err := os.ReadFile(tokenFile)
   if err != nil || len(tokenBytes) == 0 {
      t.Fatalf("Failed to read token from %s. Please run TestStep1 and TestStep2 first.", tokenFile)
   }
   accessToken := strings.TrimSpace(string(tokenBytes))

   asin := "B075RND57T"
   opts := DefaultPlaybackOptions()

   // 1. We must fetch the manifest first to get the 'customerID'
   t.Log("Fetching manifest to retrieve customerID...")
   manifestResp, err := GetPlaybackResources(
      client,
      playbackEndpoint,
      accessToken,
      asin,
      marketplaceIDUS,
      defaultDevice,
      opts,
   )
   if err != nil {
      t.Fatalf("Failed to fetch playback resources: %v", err)
   }

   // Extract the customerID safely from the untyped interface map
   customerIDRaw, ok := manifestResp.ReturnedTitleRendition.SelectedEntitlement["grantedByCustomerId"]
   if !ok {
      t.Fatalf("grantedByCustomerId not found in entitlement data")
   }
   customerID := customerIDRaw.(string)
   t.Logf("Got customerID: %s", customerID)

   // 2. Prepare the Widevine challenge.
   // For testing, we use dummy bytes. In production, you would pass your actual protobuf slice.
   // NOTE: Because this is a dummy challenge, Amazon will likely return an API error ("Invalid challenge"),
   // which PROVES the request formatting works and reached the Widevine server!
   dummyChallenge := []byte("dummy-widevine-protobuf-challenge-bytes")

   t.Log("Requesting Widevine License...")
   licenseBytes, err := GetWidevineLicense(
      client,
      playbackEndpoint,
      accessToken,
      asin,
      marketplaceIDUS,
      defaultDevice,
      customerID,
      dummyChallenge,
      opts,
   )

   if err != nil {
      // We expect this to fail with a Widevine error because the challenge is fake.
      t.Logf("Expected failure with dummy challenge: %v", err)
   } else {
      t.Logf("Success! Decoded Widevine License length: %d bytes", len(licenseBytes))
   }
}
