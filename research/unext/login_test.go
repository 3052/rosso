// login_test.go
package unext

import (
   "encoding/json"
   "fmt"
   "net/http"
   "net/http/cookiejar"
   "os"
   "os/exec"
   "testing"
)

const tokensFile = "tokens.json"

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

func TestLogin(t *testing.T) {
   // --- Fetch credentials from credential.exe ---
   creds, err := GetCredentials("unext.jp")
   if err != nil {
      t.Fatalf("GetCredentials: %v", err)
   }
   t.Logf("found %d credential(s) for unext.jp", len(creds))

   cred := creds[0]
   t.Logf("using username: %s", cred.Username)

   // --- Generate PKCE pair and random state/nonce ---
   codeVerifier, codeChallenge, err := pkcePair()
   if err != nil {
      t.Fatalf("pkcePair: %v", err)
   }

   state, err := generateRandomString(43)
   if err != nil {
      t.Fatalf("generateRandomString (state): %v", err)
   }

   nonce, err := generateRandomString(43)
   if err != nil {
      t.Fatalf("generateRandomString (nonce): %v", err)
   }

   t.Logf("code_verifier  = %s", codeVerifier)
   t.Logf("code_challenge = %s", codeChallenge)
   t.Logf("state          = %s", state)
   t.Logf("nonce          = %s", nonce)

   // --- HTTP client with cookie jar (needed for oauth_session_id) ---
   jar, err := cookiejar.New(nil)
   if err != nil {
      t.Fatalf("cookiejar.New: %v", err)
   }

   client := &http.Client{
      Jar: jar,
      // Disable auto-redirect so we can capture 302 Location headers.
      CheckRedirect: func(req *http.Request, via []*http.Request) error {
         return http.ErrUseLastResponse
      },
   }

   // --- Step 1: GET /oauth2/auth → challenge_id ---
   challengeID, err := Step1GetChallenge(client, state, nonce)
   if err != nil {
      t.Fatalf("Step1GetChallenge: %v", err)
   }
   t.Logf("challenge_id = %s", challengeID)

   // --- Step 2: POST /oauth2/login → post_auth_endpoint (+ cookie) ---
   postAuthEndpoint, err := Step2Login(client, cred.Username, cred.Password, challengeID)
   if err != nil {
      t.Fatalf("Step2Login: %v", err)
   }
   t.Logf("post_auth_endpoint = %s", postAuthEndpoint)

   // --- Step 3: POST /oauth2/auth (with code_challenge) → auth code ---
   authCode, err := Step3GetAuthCode(client, postAuthEndpoint, codeChallenge)
   if err != nil {
      t.Fatalf("Step3GetAuthCode: %v", err)
   }
   t.Logf("auth code = %s", authCode)

   // --- Step 4: POST /oauth2/token (with code_verifier) → tokens ---
   tokens, err := Step4GetToken(client, authCode, codeVerifier)
   if err != nil {
      t.Fatalf("Step4GetToken: %v", err)
   }

   // --- Assert we got meaningful tokens ---
   if tokens.AccessToken == "" {
      t.Fatal("access_token is empty")
   }
   if tokens.RefreshToken == "" {
      t.Fatal("refresh_token is empty")
   }

   // --- Print final result ---
   out, _ := json.MarshalIndent(tokens, "", "  ")
   t.Logf("tokens:\n%s", string(out))

   // --- Save tokens to file for future tests ---
   if err := SaveTokens(tokensFile, tokens); err != nil {
      t.Fatalf("SaveTokens: %v", err)
   }
   t.Logf("tokens saved to %s", tokensFile)
}

// CredentialEntry represents one entry from credential.exe output.
type CredentialEntry struct {
   Date     string `json:"date"`
   Host     string `json:"host"`
   Password string `json:"password"`
   Username string `json:"username"`
}
