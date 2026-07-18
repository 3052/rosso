// license_test.go
package unext

import (
   "encoding/hex"
   "net/http"
   "net/http/cookiejar"
   "os"
   "strings"
   "testing"

   "41.neocities.org/diana/widevine"
)

func TestGetLicense(t *testing.T) {
   // --- Load CDM device files ---
   privateKeyData, err := os.ReadFile(`D:\DRM\MIRC Electronics Ltd 8131 L3\device_private_key`)
   if err != nil {
      t.Fatalf("reading device_private_key: %v", err)
   }

   privateKey, err := widevine.DecodePrivateKey(privateKeyData)
   if err != nil {
      t.Fatalf("widevine.DecodePrivateKey: %v", err)
   }

   clientId, err := os.ReadFile(`D:\DRM\MIRC Electronics Ltd 8131 L3\device_client_id_blob`)
   if err != nil {
      t.Fatalf("reading device_client_id_blob: %v", err)
   }

   // --- Load tokens ---
   tokens, err := LoadTokens(tokensFile)
   if err != nil {
      t.Fatalf("LoadTokens: %v", err)
   }

   // --- HTTP client ---
   jar, err := cookiejar.New(nil)
   if err != nil {
      t.Fatalf("cookiejar.New: %v", err)
   }

   client := &http.Client{
      Jar: jar,
      CheckRedirect: func(req *http.Request, via []*http.Request) error {
         return http.ErrUseLastResponse
      },
   }

   // --- Step 5: GET playlist ---
   playlist, err := Step5GetPlaylist(client, tokens.AccessToken)
   if err != nil {
      t.Fatalf("Step5GetPlaylist: %v", err)
   }

   if playlist.PlayToken == "" {
      t.Fatal("playToken is empty")
   }

   // --- Build PSSH data ---
   keyId, err := hex.DecodeString(strings.ReplaceAll("296cc8c2-f941-11e8-897d-0c4de9cf85a0", "-", ""))
   if err != nil {
      t.Fatalf("decoding key ID: %v", err)
   }

   contentId, err := hex.DecodeString("4d455a30303030323439343538")
   if err != nil {
      t.Fatalf("decoding content ID: %v", err)
   }

   psshData := &widevine.PsshData{
      KeyIds:    [][]byte{keyId},
      ContentId: contentId,
   }

   // --- Build the license request ---
   requestData, err := psshData.EncodeLicenseRequest(clientId)
   if err != nil {
      t.Fatalf("EncodeLicenseRequest: %v", err)
   }

   // --- Sign the request ---
   challenge, err := widevine.EncodeSignedMessage(requestData, privateKey)
   if err != nil {
      t.Fatalf("EncodeSignedMessage: %v", err)
   }

   // --- Get the Widevine license URL ---
   licenseURL, err := playlist.WidevineLicenseURL()
   if err != nil {
      t.Fatalf("WidevineLicenseURL: %v", err)
   }

   // --- Step 6: POST license challenge ---
   licenseResponse, err := Step6GetLicense(client, licenseURL, playlist.PlayToken, challenge)
   if err != nil {
      t.Fatalf("Step6GetLicense: %v", err)
   }

   if len(licenseResponse) == 0 {
      t.Fatal("license response is empty")
   }

   t.Logf("license response: %d bytes", len(licenseResponse))
}
