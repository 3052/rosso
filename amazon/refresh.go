package amazon

import (
   "encoding/json"
   "fmt"
   "net/http"
)

// RefreshAccountToken exchanges an existing refresh token for a new TokenPair
// (new access token and potentially a new refresh token).
func RefreshAccountToken(refreshToken string) (*TokenPair, error) {
   url := "https://api.amazon.com/auth/o2/token"

   // The HAR capture shows the payload is passed via query parameters rather than the POST body.
   req, err := http.NewRequest("POST", url, nil)
   if err != nil {
      return nil, err
   }

   q := req.URL.Query()
   q.Add("refresh_token", refreshToken)
   q.Add("grant_type", "refresh_token")
   // Extracted from your HAR file
   q.Add("client_id", "amzn1.application-oa2-client.176ed4f81bb24970b60ae95ce8a7a9ac")
   req.URL.RawQuery = q.Encode()

   // Keeping the User-Agent consistent with your other endpoints
   req.Header.Set("User-Agent", "Android/google/sdk_gphone_x86/generic_x86_arm:11/RSR1.240422.006/12134477:userdebug/dev-keys, Ignition X/15.5.2026042820-android, Google")
   req.Header.Set("Content-Type", "application/json; charset=utf-8")
   req.Header.Set("Accept", "application/json")

   client := &http.Client{}
   resp, err := client.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()

   if resp.StatusCode != http.StatusOK {
      return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
   }

   // Because the JSON returned is flat: {"access_token": "...", "refresh_token": "...", ...}
   // We can decode directly into your existing TokenPair struct.
   var result TokenPair
   if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
      return nil, err
   }

   return &result, nil
}
