// step2_login.go
package unext

import (
   "bytes"
   "encoding/json"
   "fmt"
   "io"
   "net/http"
)

// Step2Login authenticates with email/password and returns the post_auth_endpoint.
func Step2Login(email, password, challengeID string) (string, error) {
   loginURL := "https://oauth.unext.jp/oauth2/login"

   body := map[string]any{
      "id":           email,
      "password":     password,
      "challenge_id": challengeID,
      "device_code":  "920",
      "scope":        []string{"offline", "unext"},
   }

   jsonBody, err := json.Marshal(body)
   if err != nil {
      return "", fmt.Errorf("step2: marshalling body: %w", err)
   }

   req, err := http.NewRequest("POST", loginURL, bytes.NewReader(jsonBody))
   if err != nil {
      return "", fmt.Errorf("step2: creating request: %w", err)
   }

   req.Header.Set("user-agent", "U-NEXT Phone App Android12 5.71.0 sdk_gphone64_x86_64")
   req.Header.Set("content-type", "application/json; charset=utf-8")
   req.Header.Set("x-forwarded-for", "159.26.119.122")

   resp, err := clientDo(req)
   if err != nil {
      return "", fmt.Errorf("step2: sending request: %w", err)
   }
   defer resp.Body.Close()

   respBody, err := io.ReadAll(resp.Body)
   if err != nil {
      return "", fmt.Errorf("step2: reading response body: %w", err)
   }

   if resp.StatusCode != http.StatusOK {
      return "", fmt.Errorf("step2: expected 200, got %d: %s", resp.StatusCode, string(respBody))
   }

   var loginResp LoginResponse
   if err := json.Unmarshal(respBody, &loginResp); err != nil {
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
