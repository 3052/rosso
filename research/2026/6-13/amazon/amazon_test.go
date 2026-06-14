package amazon

import (
   "encoding/json"
   "fmt"
   "os"
   "path/filepath"
)

type authState struct {
   PublicCode  string `json:"public_code"`
   PrivateCode string `json:"private_code"`
}

func getStateFilePath() string {
   return filepath.Join(os.TempDir(), "amazon_auth_state.json")
}

// InitiateLogin starts the authentication process, prints the code for the user,
// and saves the necessary credentials to the OS temp directory.
func InitiateLogin() error {
   publicCode, privateCode, err := CreateCodePair()
   if err != nil {
      return fmt.Errorf("failed to create code pair: %w", err)
   }

   err = InitiateMDSO(publicCode)
   if err != nil {
      return fmt.Errorf("failed to initiate MDSO: %w", err)
   }

   fmt.Printf("\n=== AMAZON LOGIN ===\n")
   fmt.Printf("Please navigate to https://www.amazon.com/us/code\n")
   fmt.Printf("Enter the following code: %s\n", publicCode)
   fmt.Printf("====================\n\n")

   state := authState{
      PublicCode:  publicCode,
      PrivateCode: privateCode,
   }

   data, err := json.Marshal(state)
   if err != nil {
      return fmt.Errorf("failed to marshal auth state: %w", err)
   }

   filePath := getStateFilePath()
   err = os.WriteFile(filePath, data, 0600)
   if err != nil {
      return fmt.Errorf("failed to write state to temp dir: %w", err)
   }

   fmt.Printf("State saved to: %s\n", filePath)
   return nil
}

// CompleteLogin reads the state from the OS temp directory and attempts to
// complete the login once.
func CompleteLogin() (string, string, error) {
   filePath := getStateFilePath()
   data, err := os.ReadFile(filePath)
   if err != nil {
      return "", "", fmt.Errorf("failed to read state from temp dir (did you run InitiateLogin first?): %w", err)
   }

   var state authState
   if err := json.Unmarshal(data, &state); err != nil {
      return "", "", fmt.Errorf("failed to unmarshal auth state: %w", err)
   }

   // Make a single attempt to register
   accountAccessToken, accountRefreshToken, err := PollRegister(state.PublicCode, state.PrivateCode)
   if err != nil {
      return "", "", fmt.Errorf("login incomplete or failed: %w", err)
   }

   fmt.Println("Login successful!")

   // Clean up the temp file after successful retrieval
   _ = os.Remove(filePath)

   return accountAccessToken, accountRefreshToken, nil
}
