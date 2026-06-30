// license_test.go
package molotov

import (
   "encoding/json"
   "os"
   "testing"
)

func TestGetLicense(t *testing.T) {
   assetFile, err := os.ReadFile("asset_test.json")
   if err != nil {
      t.Fatalf("Failed to read asset_test.json. Please run TestGetAsset first: %v", err)
   }

   var assetResp AssetResponse
   if err := json.Unmarshal(assetFile, &assetResp); err != nil {
      t.Fatalf("Failed to unmarshal AssetResponse: %v", err)
   }

   t.Logf("Requesting license from: %s", assetResp.DRM.LicenseURL)

   // Use the method on the struct
   license, err := assetResp.GetLicense(nil)
   if err != nil {
      t.Fatalf("GetLicense failed (this is expected with a nil challenge!): %v", err)
   }

   if len(license) == 0 {
      t.Fatal("Expected license data, got empty response")
   }

   t.Logf("Successfully retrieved license data, length: %d bytes", len(license))
}
