package unext

import (
   "encoding/json"
   "fmt"
   "os"
)

// SaveTokens writes the TokenResponse to a JSON file.
func SaveTokens(path string, tokens *TokenResponse) error {
   data, err := json.MarshalIndent(tokens, "", "  ")
   if err != nil {
      return fmt.Errorf("marshalling tokens: %w", err)
   }

   if err := os.WriteFile(path, data, 0600); err != nil {
      return fmt.Errorf("writing tokens file: %w", err)
   }

   return nil
}

// LoadTokens reads the TokenResponse from a JSON file.
func LoadTokens(path string) (*TokenResponse, error) {
   data, err := os.ReadFile(path)
   if err != nil {
      return nil, fmt.Errorf("reading tokens file: %w", err)
   }

   var tokens TokenResponse
   if err := json.Unmarshal(data, &tokens); err != nil {
      return nil, fmt.Errorf("parsing tokens file: %w", err)
   }

   return &tokens, nil
}
