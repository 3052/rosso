package hboMax

import (
   "testing"
)

func TestGetMovieLive(t *testing.T) {
   client := NewClient(testToken)

   movieRouteID := "bebe611d-8178-481a-a4f2-de743b5b135a"

   // Get raw entities
   entities, err := client.GetMovie(movieRouteID)
   if err != nil {
      t.Fatalf("Live API request failed: %v", err)
   }

   // Extract Edit ID
   editID, err := GetEditID(entities)
   if err != nil {
      t.Fatalf("Failed to extract Edit ID: %v", err)
   }

   if editID == "" {
      t.Fatalf("Expected an Edit ID, but got an empty string")
   }

   t.Logf("Successfully retrieved Edit ID from live server: %s", editID)
}
