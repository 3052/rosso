// asset_test.go
package molotov

import (
   "encoding/json"
   "os"
   "testing"
)

func TestGetAsset(t *testing.T) {
   // Read the auth data
   authFile, err := os.ReadFile("auth_test.json")
   if err != nil {
      t.Fatalf("Failed to read auth_test.json. Please run TestSignin first: %v", err)
   }
   var authData TestAuthData // defined in signin_test.go
   if err := json.Unmarshal(authFile, &authData); err != nil {
      t.Fatalf("Failed to unmarshal auth data: %v", err)
   }

   // Read the user data
   userFile, err := os.ReadFile("user_test.json")
   if err != nil {
      t.Fatalf("Failed to read user_test.json. Please run TestGetUser first: %v", err)
   }
   var userData TestUserData // defined in user_test.go
   if err := json.Unmarshal(userFile, &userData); err != nil {
      t.Fatalf("Failed to unmarshal user data: %v", err)
   }

   // The asset ID observed in the HAR file
   assetID := "VOD_314017"
   sessionID := "test-session-id" // Using a placeholder as session tracking isn't strictly necessary for playback URLs

   mpdURL, licenseURL, dtAuthToken, err := GetAsset(
      assetID,
      authData.AccessToken,
      userData.UserID,
      userData.ProfileID,
      sessionID,
   )
   if err != nil {
      t.Fatalf("GetAsset failed: %v", err)
   }

   if mpdURL == "" || licenseURL == "" || dtAuthToken == "" {
      t.Fatalf("Expected non-empty mpdURL, licenseURL, and dtAuthToken.\nGot mpdURL: %s\nlicenseURL: %s\ndtAuthToken: %s", mpdURL, licenseURL, dtAuthToken)
   }

   // Prepare the data to be written for the final license test
   assetData := TestAssetData{
      MpdURL:      mpdURL,
      LicenseURL:  licenseURL,
      DtAuthToken: dtAuthToken,
   }

   data, err := json.MarshalIndent(assetData, "", "  ")
   if err != nil {
      t.Fatalf("Failed to marshal asset data: %v", err)
   }

   // Write the resulting asset data to a file
   err = os.WriteFile("asset_test.json", data, 0600)
   if err != nil {
      t.Fatalf("Failed to write asset data to file: %v", err)
   }

   t.Logf("Successfully retrieved Asset Data and saved to asset_test.json\nMPD URL: %s", mpdURL)
}

// TestAssetData is used to serialize the asset data for the final license test.
type TestAssetData struct {
   MpdURL      string `json:"mpd_url"`
   LicenseURL  string `json:"license_url"`
   DtAuthToken string `json:"dt_auth_token"`
}
