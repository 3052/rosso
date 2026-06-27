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
   payload := map[string]string{
      "app_name":             "AIV",
      "requested_token_type": "access_token",
      "source_token":         refreshToken,
      "source_token_type":    "refresh_token",
   }
   body, err := json.Marshal(payload)
   if err != nil {
      return nil, err
   }
   req, err := http.NewRequest(
      "POST", HostAmazonAPI+"/auth/token", bytes.NewBuffer(body),
   )
   if err != nil {
      return nil, err
   }
   req.Header.Set("content-type", "application/json")
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
