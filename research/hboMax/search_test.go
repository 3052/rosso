package hboMax

import (
   "testing"
)

func TestSearchLive(t *testing.T) {
   client := NewClient(testToken)
   query := "marnie"

   entities, err := client.Search(query)
   if err != nil {
      t.Fatalf("Live API request failed: %v", err)
   }

   results, err := GetSearchResults(entities)
   if err != nil {
      t.Fatalf("Failed to extract results from search entities: %v", err)
   }

   if len(results) == 0 {
      t.Fatalf("Expected results for query '%s', got 0", query)
   }

   t.Log("==================================================")
   t.Logf(" Search Results for '%s'", query)
   t.Log("==================================================")
   for i, res := range results {
      t.Logf("Result %d:", i+1)
      t.Log(res.String())
      t.Log("--------------------------------------------------")
   }
}
