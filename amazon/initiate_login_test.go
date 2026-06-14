package amazon

import (
   "encoding/json"
   "os"
   "testing"
)

func TestInitiateLogin(t *testing.T) {
   // Call the updated function which now returns a *CodePair
   codes, err := CreateCodePair()
   if err != nil {
      t.Fatalf("Failed to create code pair: %v", err)
   }

   // Access the properties using dot notation
   err = InitiateMDSO(codes.PublicCode)
   if err != nil {
      t.Fatalf("Failed to initiate MDSO: %v", err)
   }

   t.Logf("\n=== AMAZON LOGIN ===\nPlease navigate to https://www.amazon.com/us/code\nEnter the following code: %s\n====================\n", codes.PublicCode)

   state := authState{
      PublicCode:  codes.PublicCode,
      PrivateCode: codes.PrivateCode,
   }

   data, err := json.Marshal(state)
   if err != nil {
      t.Fatalf("Failed to marshal auth state: %v", err)
   }

   filePath := getStateFilePath()
   err = os.WriteFile(filePath, data, 0600)
   if err != nil {
      t.Fatalf("Failed to write state to temp dir: %v", err)
   }

   t.Logf("State saved to: %s", filePath)
}
