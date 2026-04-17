package plex

import (
   "fmt"
   "io"
   "net/http"
)

// CreateAnonymousUser creates a new anonymous session and returns the JSON response
// containing the authToken.
func CreateAnonymousUser() ([]byte, error) {
   req, err := http.NewRequest("POST", "https://plex.tv/api/v2/users/anonymous", nil)
   if err != nil {
      return nil, err
   }

   req.Header.Set("Accept", "application/json")
   req.Header.Set("X-Plex-Client-Identifier", "!") // Unique ID, HAR used "!"
   req.Header.Set("X-Plex-Product", "Plex Mediaverse")

   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()

   if resp.StatusCode < 200 || resp.StatusCode >= 300 {
      return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
   }

   return io.ReadAll(resp.Body)
}
