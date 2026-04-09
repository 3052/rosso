package hboMax

import (
   "testing"
)

func TestGetMovieEditIDLive(t *testing.T) {
   client := NewClient(testToken)

   // Hit the live API for the specific movie ID
   movieRouteID := "bebe611d-8178-481a-a4f2-de743b5b135a"

   editID, err := client.GetMovieEditID(movieRouteID)
   if err != nil {
      t.Fatalf("Live API request failed: %v", err)
   }

   if editID == "" {
      t.Fatalf("Expected an Edit ID, but got an empty string")
   }

   t.Logf("Successfully retrieved Edit ID from live server: %s", editID)
}
