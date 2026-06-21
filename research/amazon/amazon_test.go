package amazon

import (
   "fmt"
   "net/http"
   "testing"
)

func TestGetManifestAndLicense(t *testing.T) {
   // -------------------------------------------------------------------------
   // Setup: Fill these with your actual scraped/extracted values to run a live test
   // -------------------------------------------------------------------------
   deviceID := ""      // e.g. "uuidb43bee409bd448cfb5ba3337bd241645"
   authBearer := ""    // e.g. "Atna|EwMDICIPxLGAmnVlZgnFhnKMSRVvjHua..."
   titleID := ""       // e.g. "amzn1.dv.gti.af991753-e4cf-4d28-880d-dfca3d1e8d24"
   marketplaceID := "" // e.g. "ATVPDKIKX0DER"
   playbackEnv := ""   // e.g. "MDJ8Cm0KBHBlbnYSJGI1YWQ0MjdhLTIyY2MtN..."

   if authBearer == "" || titleID == "" {
      t.Skip("Skipping live Amazon API test: missing credentials or title info")
   }

   // 1. Initialize the client
   client := NewClient(&http.Client{})

   // 2. Select the desired profile
   profile := WidevineL3Profile
   profile.DeviceID = deviceID
   profile.AuthBearer = authBearer

   // 3. Request the Manifest
   mpdURL, handoffToken, err := client.GetManifest(profile, titleID, marketplaceID, playbackEnv)
   if err != nil {
      t.Fatalf("Failed to get manifest: %v", err)
   }

   if mpdURL == "" {
      t.Fatal("Received empty MPD URL")
   }
   if handoffToken == "" {
      t.Fatal("Received empty Handoff Token")
   }

   fmt.Printf("MPD URL: %s\n", mpdURL)
   fmt.Printf("Handoff Token: %s\n", handoffToken)

   // 4. In a real scenario, you would download the MPD, parse it, find the lowest
   // quality representation, and extract its PSSH. Your CDM would then generate
   // the challenge.
   //
   // For Widevine: challengeBytes is the raw byte array from the CDM.
   // For PlayReady: challengeBytes is the UTF-8 encoded XML SOAP envelope string cast to []byte.

   /* Uncomment to test actual license retrieval once PSSH is parsed:
   var mockChallenge []byte = []byte("mock_challenge_from_cdm")
   licenseB64, err := client.GetLicense(profile, titleID, marketplaceID, playbackEnv, handoffToken, mockChallenge)
   if err != nil {
      t.Fatalf("Failed to get license: %v", err)
   }

   if licenseB64 == "" {
      t.Fatal("Received empty license base64 string")
   }

   fmt.Printf("License successfully retrieved! Base64 Length: %d\n", len(licenseB64))
   */
}
