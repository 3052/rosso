// step2_login.go
package unext

import (
   "bytes"
   "encoding/json"
   "fmt"
   "net/http"
)

// Step2Login authenticates with email/password and returns the post_auth_endpoint.
func Step2Login(email, password, challengeID string) (string, error) {
   loginURL := "https://oauth.unext.jp/oauth2/login"
   body := map[string]any{
      "id":           email,
      "password":     password,
      "challenge_id": challengeID,
      "scope": []string{
         "offline",
         "unext",
      },
   }
   jsonBody, err := json.Marshal(body)
   if err != nil {
      return "", fmt.Errorf("step2: marshalling body: %w", err)
   }
   req, err := http.NewRequest("POST", loginURL, bytes.NewReader(jsonBody))
   if err != nil {
      return "", fmt.Errorf("step2: creating request: %w", err)
   }
   req.Header.Set("x-forwarded-for", "159.26.119.122")
   resp, err := clientDo(req)
   if err != nil {
      return "", fmt.Errorf("step2: sending request: %w", err)
   }
   defer resp.Body.Close()

   if resp.StatusCode != http.StatusOK {
      return "", fmt.Errorf("step2: expected 200, got %d", resp.StatusCode)
   }

   var loginResp LoginResponse
   if err := json.NewDecoder(resp.Body).Decode(&loginResp); err != nil {
      return "", fmt.Errorf("step2: parsing response: %w", err)
   }

   if loginResp.PostAuthEndpoint == "" {
      return "", fmt.Errorf("step2: post_auth_endpoint is empty")
   }

   return loginResp.PostAuthEndpoint, nil
}

// LoginResponse is the JSON returned by /oauth2/login.
type LoginResponse struct {
   PostAuthEndpoint string `json:"post_auth_endpoint"`
}
