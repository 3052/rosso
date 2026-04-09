package hboMax

import (
   "testing"
)

func TestSearchLive(t *testing.T) {
   client := NewClient(testToken)

   query := "marnie"

   // 1. Hit the live API searching for "marnie"
   searchResp, err := client.Search(query)
   if err != nil {
      t.Fatalf("Live API request failed: %v", err)
   }

   // 2. Use the method to extract and map the results
   results, err := searchResp.GetResults()
   if err != nil {
      t.Fatalf("Failed to extract results from search response: %v", err)
   }

   if len(results) == 0 {
      t.Fatalf("Expected results for query '%s', got 0", query)
   }

   // Print all search results
   t.Log("---------------------------------------------------------")
   t.Logf("Search Results for '%s':", query)
   t.Log("---------------------------------------------------------")

   for i, res := range results {
      t.Logf("%2d. %s [%s]", i+1, res.Name, res.MediaType)
   }
   t.Log("---------------------------------------------------------")
}
