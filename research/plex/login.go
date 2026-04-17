package plex

import (
   "encoding/json"
   "fmt"
   "net/http"
)

type AnonymousUserResponse struct {
   ID        int    `json:"id"`
   UUID      string `json:"uuid"`
   AuthToken string `json:"authToken"`
}

// CreateAnonymousUser creates a new session and returns the decoded JSON response
// so you can extract the AuthToken for subsequent requests.
func CreateAnonymousUser(clientID string) (*AnonymousUserResponse, error) {
   req, err := http.NewRequest("POST", "https://plex.tv/api/v2/users/anonymous", nil)
   if err != nil {
      return nil, err
   }

   // Explicitly set 0 since Plex expects it for this empty POST
   req.ContentLength = 0

   req.Header.Set("Accept", "application/json")
   req.Header.Set("X-Plex-Client-Identifier", clientID)
   req.Header.Set("X-Plex-Product", "Plex Mediaverse")

   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()

   if resp.StatusCode < 200 || resp.StatusCode >= 300 {
      return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
   }

   var result AnonymousUserResponse
   if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
      return nil, err
   }

   return &result, nil
}
