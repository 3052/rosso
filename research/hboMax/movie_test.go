package hboMax

import (
   "testing"
)

func TestGetMovieLive(t *testing.T) {
   client := NewClient(testToken)

   // Hit the live API for the specific movie ID
   movieRouteID := "bebe611d-8178-481a-a4f2-de743b5b135a"

   // 1. Get the Movie Response
   movieResp, err := client.GetMovie(movieRouteID)
   if err != nil {
      t.Fatalf("Live API request failed: %v", err)
   }

   // 2. Use the method to extract the Edit ID
   editID, err := movieResp.GetEditID()
   if err != nil {
      t.Fatalf("Failed to extract Edit ID: %v", err)
   }

   if editID == "" {
      t.Fatalf("Expected an Edit ID, but got an empty string")
   }

   t.Logf("Successfully retrieved Edit ID from live server: %s", editID)
}
