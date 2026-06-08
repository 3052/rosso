// credentials.go
package amazon

import (
   "encoding/json"
   "fmt"
   "os/exec"
)

type Credential struct {
   Date     string `json:"date"`
   Host     string `json:"host"`
   Password string `json:"password"`
   Trial    string `json:"trial"`
   Username string `json:"username"`
}

func GetPhoneNumber() (string, error) {
   cmd := exec.Command("credential.exe", "-j=amazon.com")
   output, err := cmd.Output()
   if err != nil {
      return "", fmt.Errorf("failed to execute credential.exe: %w", err)
   }

   var creds []Credential
   if err := json.Unmarshal(output, &creds); err != nil {
      return "", fmt.Errorf("failed to parse credentials JSON: %w", err)
   }

   if len(creds) == 0 {
      return "", fmt.Errorf("no credentials found for amazon.com")
   }

   // Return the username of the first entry (acting as the phone number)
   return creds[0].Username, nil
}
