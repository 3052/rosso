package amazon

import (
   "bytes"
   "encoding/json"
   "fmt"
   "net/http"
)

// Refresh exchanges the existing refresh token for a new access token
// using the /auth/token endpoint, mutating the TokenPair in-place.
func (t *TokenPair) Refresh() error {
   if t == nil || t.RefreshToken == "" {
      return fmt.Errorf("invalid token pair or missing refresh token")
   }

   payload := map[string]string{
      "app_name":             "AIV",
      "requested_token_type": "access_token",
      "source_token":         t.RefreshToken,
      "source_token_type":    "refresh_token",
   }
   body, err := json.Marshal(payload)
   if err != nil {
      return err
   }
   req, err := http.NewRequest(
      "POST", HostAmazonAPI+"/auth/token", bytes.NewBuffer(body),
   )
   if err != nil {
      return err
   }
   req.Header.Set("content-type", "application/json")
   client := &http.Client{}
   resp, err := client.Do(req)
   if err != nil {
      return err
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
      return err
   }

   // Handle API errors as seen in the Python code
   if result.Error != "" {
      return fmt.Errorf("failed to refresh device token: %s [%s]", result.ErrorDesc, result.Error)
   }

   if result.TokenType != "bearer" {
      return fmt.Errorf("unexpected returned refreshed token type: %s", result.TokenType)
   }

   // Mutate the struct in-place with the new access token
   t.AccessToken = result.AccessToken

   return nil
}
