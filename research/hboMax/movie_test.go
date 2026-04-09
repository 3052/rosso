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

   // Extract formatted movie entities
   movies := GetVideos(entities)

   if len(movies) == 0 {
      t.Fatalf("Expected at least one movie entity, got 0")
   }

   t.Log("==================================================")
   t.Logf(" Found %d Movie Entities", len(movies))
   t.Log("==================================================")

   for i, movie := range movies {
      // Fallback to ID if Name is empty in the video entity
      name := movie.Attributes.Name
      if name == "" {
         name = "Unknown Name (Check Route ID)"
      }

      t.Logf("Movie %d: %s", i+1, name)
      t.Logf("Edit ID:    %s", movie.Relationships.Edit.Data.ID)
      t.Logf("Video Type: %s", movie.Attributes.VideoType)
      if movie.Attributes.Description != "" {
         t.Logf("Summary:    %s", movie.Attributes.Description)
      }
      t.Log("--------------------------------------------------")

      if movie.Relationships.Edit.Data.ID == "" {
         t.Errorf("Movie %d is missing an Edit ID", i+1)
      }
   }
}
