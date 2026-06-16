package amazon

import (
   "bytes"
   "encoding/json"
   "fmt"
   "net/http"
)

// RefreshToken exchanges an existing refresh token for a new access token
// using the /auth/token endpoint.
func RefreshToken(refreshToken string) (*TokenPair, error) {
   url := "https://api.amazon.com/auth/token"

   // Mimicking the Python payload structure:
   // "app_name": device["app_name"],
   // "app_version": device["app_version"],
   // "source_token_type": "refresh_token",
   // "source_token": refresh_token,
   // "requested_token_type": "access_token"
   payload := map[string]interface{}{
      "app_name":             "AIV",    // Kept consistent with your other Go files
      "app_version":          "3.12.0", // Kept consistent with your other Go files
      "source_token_type":    "refresh_token",
      "source_token":         refreshToken,
      "requested_token_type": "access_token",
   }

   body, err := json.Marshal(payload)
   if err != nil {
      return nil, err
   }

   req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
   if err != nil {
      return nil, err
   }

   // Keeping headers consistent with your Go codebase
   req.Header.Set("User-Agent", "Android/google/sdk_gphone_x86/generic_x86_arm:11/RSR1.240422.006/12134477:userdebug/dev-keys, Ignition X/15.5.2026042820-android, Google")
   req.Header.Set("Content-Type", "application/json")
   req.Header.Set("Accept", "application/json")

   client := &http.Client{}
   resp, err := client.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()

   // Decode into an anonymous struct handling the expected Python response keys
   var result struct {
      AccessToken string `json:"access_token"`
      TokenType   string `json:"token_type"`
      Error       string `json:"error"`
      ErrorDesc   string `json:"error_description"`
   }

   if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
      return nil, err
   }

   // Handle API errors as seen in the Python code
   if result.Error != "" {
      return nil, fmt.Errorf("failed to refresh device token: %s [%s]", result.ErrorDesc, result.Error)
   }

   if result.TokenType != "bearer" {
      return nil, fmt.Errorf("unexpected returned refreshed token type: %s", result.TokenType)
   }

   // The refresh endpoint typically only returns a new access_token.
   // We return your TokenPair carrying forward the original refresh_token
   // (just like the Python script does with: refreshed_tokens["refresh_token"] = cache["refresh_token"])
   return &TokenPair{
      AccessToken:  result.AccessToken,
      RefreshToken: refreshToken,
   }, nil
}
