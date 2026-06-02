package amazon

import (
   "41.neocities.org/diana/widevine"
   "encoding/hex"
   "os"
   "strings"
   "testing"
)

// Run this to test requesting a Widevine License with a real CDM
func TestStep4_GetLicense(t *testing.T) {
   // 1. Read Amazon Access Token
   tokenFile := getTempTokenPath()
   tokenBytes, err := os.ReadFile(tokenFile)
   if err != nil || len(tokenBytes) == 0 {
      t.Fatalf("Failed to read token from %s. Please run TestStep1 and TestStep2 first.", tokenFile)
   }
   accessToken := strings.TrimSpace(string(tokenBytes))

   // 2. Load CDM Files
   t.Log("Loading CDM files...")
   clientIDPath := `C:\Users\Steven\AppData\Local\L3\client_id.bin`
   privateKeyPath := `C:\Users\Steven\AppData\Local\L3\private_key.pem`

   clientIDBytes, err := os.ReadFile(clientIDPath)
   if err != nil {
      t.Fatalf("Failed to read client_id.bin: %v", err)
   }

   privKeyBytes, err := os.ReadFile(privateKeyPath)
   if err != nil {
      t.Fatalf("Failed to read private_key.pem: %v", err)
   }

   privateKey, err := widevine.DecodePrivateKey(privKeyBytes)
   if err != nil {
      t.Fatalf("Failed to decode private key: %v", err)
   }

   // 3. Setup playback options
   asin := "B075RND57T"
   // DOCUMENTATION: Use "HD" or "UHD" here if you have an L1 CDM.
   DefaultPlaybackOptions.VideoQuality = "SD" // Swapped to SD to match the MPD requested in Step 3!

   // 4. Prepare the PSSH/Challenge data
   targetKeyIDs := []string{
      // HD
      //"AE7133F4-B2D9-4F32-96E0-F7C0089493BC",
      // SD
      "F661444B-15A3-45F8-B06D-13541B98B2E5",
   }

   var keyIDBytes [][]byte
   for _, kid := range targetKeyIDs {
      cleanKeyID := strings.ReplaceAll(kid, "-", "")
      kb, err := hex.DecodeString(cleanKeyID)
      if err != nil {
         t.Fatalf("Failed to decode Key ID hex: %v", err)
      }
      keyIDBytes = append(keyIDBytes, kb)
   }

   psshData := &widevine.PsshData{
      KeyIds: keyIDBytes,
   }

   requestData, err := psshData.EncodeLicenseRequest(clientIDBytes)
   if err != nil {
      t.Fatalf("Failed to encode license request: %v", err)
   }

   signedChallenge, err := widevine.EncodeSignedMessage(requestData, privateKey)
   if err != nil {
      t.Fatalf("Failed to sign challenge: %v", err)
   }

   // 5. Fetch the Manifest to get the Customer ID
   t.Logf("Fetching manifest for ASIN %s to retrieve customerID...", asin)
   manifestResp, err := GetPlaybackResources(
      accessToken,
      asin,
      marketplaceIDUS,
   )
   if err != nil {
      t.Fatalf("Failed to fetch playback resources: %v", err)
   }

   customerID := manifestResp.ReturnedTitleRendition.SelectedEntitlement.GrantedByCustomerId
   if customerID == "" {
      t.Fatalf("grantedByCustomerId not found in entitlement data")
   }
   t.Logf("Got customerID: %s", customerID)

   // 6. Request the Widevine License from Amazon
   t.Log("Requesting Widevine License from Amazon...")
   amazonLicenseBytes, err := GetWidevineLicense(
      accessToken,
      asin,
      marketplaceIDUS,
      customerID,
      signedChallenge,
   )
   if err != nil {
      t.Fatalf("Amazon License request failed: %v", err)
   }
   t.Logf("Received Amazon License response! (%d bytes)", len(amazonLicenseBytes))

   // 7. Decode the License Response to get the Content Keys
   t.Log("Decrypting content keys...")
   keyContainers, err := widevine.DecodeLicenseResponse(amazonLicenseBytes, requestData, privateKey)
   if err != nil {
      t.Fatalf("Failed to decode license response: %v", err)
   }

   if len(keyContainers) == 0 {
      t.Fatalf("License response parsed successfully, but no keys were returned!")
   }

   // 8. Print the extracted keys
   t.Log("=====================================================")
   t.Log("SUCCESS! Decrypted Content Keys:")
   for _, kc := range keyContainers {
      kidHex := hex.EncodeToString(kc.Id)
      keyHex := hex.EncodeToString(kc.Key)
      t.Logf("KID: %s\nKEY: %s\n", kidHex, keyHex)
   }
   t.Log("=====================================================")
}
