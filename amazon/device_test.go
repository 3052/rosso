package amazon

import (
   "encoding/json"
   "fmt"
   "os"
   "path/filepath"
   "testing"
)

// 1. Run this to get the code
func TestStep1_StartProcess(t *testing.T) {
   t.Log("Fetching Code Pair...")
   codePair, err := GetCodePair()
   if err != nil {
      t.Fatalf("Failed to get Code Pair: %v", err)
   }

   data, _ := json.MarshalIndent(codePair, "", "  ")
   stateFile := getTempStatePath()
   if err := os.WriteFile(stateFile, data, 0644); err != nil {
      t.Fatalf("Failed to write state file to temp dir: %v", err)
   }

   fmt.Println()
   fmt.Println("========================================================================")
   fmt.Println("1. Go to: https://amazon.com/mytv")
   fmt.Println("2. Sign in to your Amazon account normally using your email & password.")
   fmt.Println("3. Once signed in, you will be on the 'Register Your Device' screen.")
   fmt.Printf("4. Enter this code: %s\n", codePair.PublicCode)
   fmt.Println("5. Click 'Register Device'.")
   fmt.Println("6. Finally, come back here and run TestStep2_CompleteProcess")
   fmt.Println("========================================================================")
   fmt.Println()
}

// 2. Run this after approving the code in your browser
func TestStep2_CompleteProcess(t *testing.T) {
   stateFile := getTempStatePath()
   data, err := os.ReadFile(stateFile)
   if err != nil {
      t.Fatalf("Failed to read state file from temp dir. Run TestStep1 first. Error: %v", err)
   }

   var codePair CodePairResponse

   if err := json.Unmarshal(data, &codePair); err != nil {
      t.Fatalf("Failed to parse state file: %v", err)
   }

   t.Log("Registering device (checking if you entered the code)...")
   regResponse, err := RegisterDevice(&codePair)
   if err != nil {
      t.Fatalf("Failed to register device: %v\n(Did you forget to enter the code at amazon.com/mytv?)", err)
   }

   bearer := regResponse.Response.Success.Tokens.Bearer

   // Write token to temp dir for Step 3 to use
   tokenFile := getTempTokenPath()
   if err := os.WriteFile(tokenFile, []byte(bearer.AccessToken), 0644); err != nil {
      t.Fatalf("Failed to save access token to temp dir: %v", err)
   }

   fmt.Println()
   fmt.Println("=====================================================")
   fmt.Println("SUCCESS! Final Credentials Generated:")
   fmt.Printf("Access Token: %s\n", bearer.AccessToken)
   if bearer.RefreshToken != "" {
      fmt.Printf("Refresh Token: %s\n", bearer.RefreshToken)
   }
   fmt.Printf("Expires In: %s seconds\n", bearer.ExpiresIn)
   fmt.Printf("Token saved to: %s\n", tokenFile)
   fmt.Println("=====================================================")
   fmt.Println()

   // Clean up the intermediate state file
   _ = os.Remove(stateFile)
}

//type AuthState struct {
//   CodePair *CodePairResponse `json:"code_pair"`
//   Device   map[string]string `json:"device"`
//}

// Helper functions for temp files
func getTempStatePath() string {
   return filepath.Join(os.TempDir(), "amazon_auth_state.json")
}

func getTempTokenPath() string {
   return filepath.Join(os.TempDir(), "amazon_token.txt")
}
