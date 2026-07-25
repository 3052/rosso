package peacock

import (
   "encoding/json"
   "fmt"
   "net/http"
   "net/http/cookiejar"
   "os/exec"
)

const (
   BaseID  = "https://rango.id.peacocktv.com"
   BaseSAS = "https://sas.peacocktv.com"
)

// GetFirstEmail returns the username of the first credential entry.
func GetFirstEmail(host string) (string, error) {
   creds, err := GetCredentials(host)
   if err != nil {
      return "", err
   }
   return creds[0].Username, nil
}

// NewClient creates an *http.Client with a cookie jar.
func NewClient() (*http.Client, error) {
   jar, err := cookiejar.New(nil)
   if err != nil {
      return nil, err
   }
   return &http.Client{Jar: jar}, nil
}

// SkyHeaders returns the common headers used across all requests.
func SkyHeaders() http.Header {
   h := http.Header{}
   h.Set("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:140.0) Gecko/20100101 Firefox/140.0")
   h.Set("accept", "application/vnd.siren+json")
   h.Set("accept-language", "en-US,en;q=0.5")
   // Deliberately omit "accept-encoding": Go's transport adds gzip
   // automatically and transparently decompresses it for us.
   h.Set("referer", "https://www.peacocktv.com/")
   h.Set("origin", "https://www.peacocktv.com")
   h.Set("x-skyott-platform", "PC")
   h.Set("x-skyott-device", "COMPUTER")
   h.Set("x-skyott-provider", "NBCU")
   h.Set("x-skyott-proposition", "NBCUOTT")
   h.Set("x-skyott-territory", "US")
   return h
}

// Credential represents one entry returned by credential.exe.
type Credential struct {
   Date     string `json:"date"`
   Host     string `json:"host"`
   Password string `json:"password"`
   Username string `json:"username"`
}

// GetCredentials invokes credential.exe with the given host.
func GetCredentials(host string) ([]Credential, error) {
   cmd := exec.Command("credential.exe", fmt.Sprintf("-j=%s", host))
   out, err := cmd.Output()
   if err != nil {
      return nil, fmt.Errorf("running credential.exe: %w", err)
   }

   var creds []Credential
   if err := json.Unmarshal(out, &creds); err != nil {
      return nil, fmt.Errorf("parsing credential.exe output: %w", err)
   }
   if len(creds) == 0 {
      return nil, fmt.Errorf("no credentials returned for host %q", host)
   }
   return creds, nil
}
