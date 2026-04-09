package hboMax

import (
   "testing"
)

func TestSearchLive(t *testing.T) {
   client := NewClient(testToken)
   query := "marnie"

   // Get raw entities
   entities, err := client.Search(query)
   if err != nil {
      t.Fatalf("Live API request failed: %v", err)
   }

   // Extract formatted results
   results, err := entities.GetSearchResults()
   if err != nil {
      t.Fatalf("Failed to extract results from search entities: %v", err)
   }

   if len(results) == 0 {
      t.Fatalf("Expected results for query '%s', got 0", query)
   }

   t.Log("---------------------------------------------------------")
   t.Logf("Search Results for '%s':", query)
   t.Log("---------------------------------------------------------")
   for i, res := range results {
      t.Logf("%2d. %s [%s]", i+1, res.Name, res.MediaType)
   }
   t.Log("---------------------------------------------------------")
}
