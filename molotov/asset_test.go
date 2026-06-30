// asset_test.go
package molotov

import (
   "encoding/json"
   "os"
   "testing"
)

func TestGetAsset(t *testing.T) {
   authFile, err := os.ReadFile("auth_test.json")
   if err != nil {
      t.Fatalf("Failed to read auth_test.json. Please run TestSignin first: %v", err)
   }
   var signinResp SigninResponse
   if err := json.Unmarshal(authFile, &signinResp); err != nil {
      t.Fatalf("Failed to unmarshal auth data: %v", err)
   }

   userFile, err := os.ReadFile("user_test.json")
   if err != nil {
      t.Fatalf("Failed to read user_test.json. Please run TestGetUser first: %v", err)
   }
   var userResp UserResponse
   if err := json.Unmarshal(userFile, &userResp); err != nil {
      t.Fatalf("Failed to unmarshal user data: %v", err)
   }

   assetID := "VOD_314017"

   // Using the structs from the previous requests
   assetResp, err := GetAsset(assetID, &signinResp, &userResp)
   if err != nil {
      t.Fatalf("GetAsset failed: %v", err)
   }

   if assetResp.Stream.URL == "" || assetResp.DRM.LicenseURL == "" || assetResp.DRM.Token == "" {
      t.Fatalf("Expected non-empty values in AssetResponse. Got URL: %s", assetResp.Stream.URL)
   }

   // Directly save the AssetResponse struct for the next test
   data, err := json.MarshalIndent(assetResp, "", "  ")
   if err != nil {
      t.Fatalf("Failed to marshal asset data: %v", err)
   }

   err = os.WriteFile("asset_test.json", data, 0600)
   if err != nil {
      t.Fatalf("Failed to write asset data to file: %v", err)
   }

   t.Logf("Successfully retrieved AssetResponse and saved to asset_test.json\nMPD URL: %s", assetResp.Stream.URL)
}
