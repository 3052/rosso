package unext

import (
   "encoding/json"
   "fmt"
   "os/exec"
)

// GetCredentials calls `credential.exe -j <host>` and parses the JSON array.
func GetCredentials(host string) ([]CredentialEntry, error) {
   cmd := exec.Command("credential.exe", "-j", host)
   output, err := cmd.Output()
   if err != nil {
      return nil, fmt.Errorf("credential.exe: %w", err)
   }

   var entries []CredentialEntry
   if err := json.Unmarshal(output, &entries); err != nil {
      return nil, fmt.Errorf("parsing credential output: %w", err)
   }

   if len(entries) == 0 {
      return nil, fmt.Errorf("no credentials found for host %s", host)
   }

   return entries, nil
}

// CredentialEntry represents one entry from credential.exe output.
type CredentialEntry struct {
   Date     string `json:"date"`
   Host     string `json:"host"`
   Password string `json:"password"`
   Username string `json:"username"`
}
