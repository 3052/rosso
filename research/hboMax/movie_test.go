package hboMax

import (
   "testing"
)

func TestGetMovieLive(t *testing.T) {
   client := NewClient(testToken)

   movieRouteID := "bebe611d-8178-481a-a4f2-de743b5b135a"

   entities, err := client.GetMovie(movieRouteID)
   if err != nil {
      t.Fatalf("Live API request failed: %v", err)
   }

   movies := GetMovies(entities)

   if len(movies) == 0 {
      t.Fatalf("Expected at least one movie entity, got 0")
   }

   t.Log("==================================================")
   t.Logf(" Found %d Movie Entities", len(movies))
   t.Log("==================================================")
   for i, movie := range movies {
      t.Log(movie.String())
      t.Log("--------------------------------------------------")

      if movie.Relationships.Edit.Data.ID == "" {
         t.Errorf("Movie %d is missing an Edit ID", i+1)
      }
   }
}
