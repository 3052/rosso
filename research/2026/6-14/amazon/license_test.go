package amazon

import (
   "encoding/base64"
   "encoding/hex"
   "encoding/json"
   "os"
   "testing"

   "41.neocities.org/diana/widevine"
)

func TestGetWidevineLicense(t *testing.T) {
   actorFilePath := getActorTokensFilePath()
   actorData, err := os.ReadFile(actorFilePath)
   if err != nil {
      t.Fatalf("Failed to read actor tokens from temp dir (did you run TestGetActorToken first?): %v", err)
   }

   var actorState actorTokenState
   if err := json.Unmarshal(actorData, &actorState); err != nil {
      t.Fatalf("Failed to unmarshal actor state: %v", err)
   }

   envFilePath := getEnvelopeFilePath()
   envData, err := os.ReadFile(envFilePath)
   if err != nil {
      t.Fatalf("Failed to read envelope state from temp dir (did you run TestGetItemDetails first?): %v", err)
   }

   var envState envelopeState
   if err := json.Unmarshal(envData, &envState); err != nil {
      t.Fatalf("Failed to unmarshal envelope state: %v", err)
   }

   titleId := "amzn1.dv.gti.28b85d90-1338-720b-4be7-3247683a7624"

   // Load CDM files
   clientIdPath := `C:\Users\Steven\AppData\Local\L3\client_id.bin`
   privateKeyPath := `C:\Users\Steven\AppData\Local\L3\private_key.pem`

   clientId, err := os.ReadFile(clientIdPath)
   if err != nil {
      t.Fatalf("Failed to read client ID: %v", err)
   }

   privateKeyPem, err := os.ReadFile(privateKeyPath)
   if err != nil {
      t.Fatalf("Failed to read private key: %v", err)
   }

   privateKey, err := widevine.DecodePrivateKey(privateKeyPem)
   if err != nil {
      t.Fatalf("Failed to decode private key: %v", err)
   }

   // Prepare PSSH data
   psshB64 := "CAESENlROAqEW0/uqlhMioe8fWMaBmFtYXpvbiI1Y2lkOjB4YUQ4bzJPUm9XYTluZHFSVjlqRGc9PSwyVkU0Q29SYlQrNnFXRXlLaDd4OVl3PT0qAlNEMgA="
   psshBytes, err := base64.StdEncoding.DecodeString(psshB64)
   if err != nil {
      t.Fatalf("Failed to decode base64 PSSH: %v", err)
   }

   psshData, err := widevine.DecodePsshData(psshBytes)
   if err != nil {
      t.Fatalf("Failed to decode PSSH data: %v", err)
   }

   // Generate and sign the license request
   requestData, err := psshData.EncodeLicenseRequest(clientId)
   if err != nil {
      t.Fatalf("Failed to encode license request: %v", err)
   }

   signedRequest, err := widevine.EncodeSignedMessage(requestData, privateKey)
   if err != nil {
      t.Fatalf("Failed to sign license request: %v", err)
   }

   // Fetch the license from Amazon
   licenseData, err := GetWidevineLicense(actorState.AccessToken, titleId, envState.PlaybackEnvelope, signedRequest)
   if err != nil {
      t.Fatalf("Failed to get Widevine license: %v", err)
   }

   t.Logf("Successfully retrieved Widevine license! (Length: %d bytes)", len(licenseData))

   // Decode the license response to get the keys
   keys, err := widevine.DecodeLicenseResponse(licenseData, requestData, privateKey)
   if err != nil {
      t.Fatalf("Failed to decode license response: %v", err)
   }

   // Retrieve the specific key
   targetKeyIdHex := "d951380a845b4feeaa584c8a87bc7d63"
   targetKeyId, err := hex.DecodeString(targetKeyIdHex)
   if err != nil {
      t.Fatalf("Failed to decode target key ID hex: %v", err)
   }

   key, err := widevine.GetKey(keys, targetKeyId)
   if err != nil {
      t.Fatalf("Failed to get key from license: %v", err)
   }

   t.Logf("Decryption Key: %s:%s", targetKeyIdHex, hex.EncodeToString(key))
}
