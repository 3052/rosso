package hboMax

import (
   "testing"
)

// The 'st' token from your original files.
// NOTE: If the tests start failing with non-200 status codes (e.g., 401 Unauthorized),
// you will need to replace this string with a fresh token from your browser.
const testToken = "eyJhbGciOiJSUzI1NiJ9.eyJqdGkiOiJ0b2tlbi1hMmNlZThjMy0zNGNhLTQ0YTEtYjM4NC04YzIzOWNkZmQxZWQiLCJpc3MiOiJmcGEtaXNzdWVyIiwic3ViIjoiVVNFUklEOmJvbHQ6MGQ0NWNjZjgtYjRhMi00MTQ3LWJiZWItYzdiY2IxNDBmMzgyIiwiaWF0IjoxNzc1NjE5MTE5LCJleHAiOjIwOTA5NzkxMTksInR5cGUiOiJBQ0NFU1NfVE9LRU4iLCJzdWJkaXZpc2lvbiI6ImJlYW1fZW1lYSIsInNjb3BlIjoiZGVmYXVsdCIsImlpZCI6ImJlMzI5MzdhLTU3MWEtNDAzMC1hZWIyLTQ1MWViZjI3M2M5YiIsInZlcnNpb24iOiJ2NCIsImFub255bW91cyI6ZmFsc2UsImRldmljZUlkIjoiZjY4MTU2NGMtMWJlNS00NDk1LTg4MmItNmVmYzA2Y2Q4YTlkIn0.kkxM9-egjkpxnz2fSft9G1cQMdfFh9qK8_DHTk2D7Zb43FpORAkUbU92X7o-AMZxPl9pQfDlsE4KWmJHIB3vQUAC5WJmJHUDC2jc7nFYvKhDJfFLDcZD7Jc6TvpNrIYkbhP0gfF_lAxImYfoUFAQx9XzGWFiVfGe1Sy8lalVMwF-nQBdNSPxGijg1IAp-8Nt4xIScM3RScJDaJ7LqQzpNc4p9vK1l68oVUXA-NsE1RpB6caS7AucluygtjVSIGqtLE2HNDMQhJijPdCvYjRmNrQq30Ke_6tC6ezGIj5OD3Z2Sm4lJ0gFdzMZu_MggPUyadEbbK2LDI9nTU5qch1RYw"

func TestSearchLive(t *testing.T) {
   client := NewClient(testToken)

   // Hit the live API searching for "marnie"
   results, err := client.Search("marnie")
   if err != nil {
      t.Fatalf("Live API request failed: %v", err)
   }

   if len(results) == 0 {
      t.Fatalf("Expected results for query 'marnie', got 0")
   }

   t.Logf("Successfully retrieved %d search results from live server. First item: %s [%s]",
      len(results), results[0].Name, results[0].MediaType)
}

func TestGetSeasonEpisodesLive(t *testing.T) {
   client := NewClient(testToken)

   // Hit the live API for the specific show and season 2
   showID := "4ffd33c9-e0d6-4cd6-bd13-34c266c79be0"
   seasonNumber := 2

   episodes, err := client.GetSeasonEpisodes(showID, seasonNumber)
   if err != nil {
      t.Fatalf("Live API request failed: %v", err)
   }

   if len(episodes) == 0 {
      t.Fatalf("Expected episodes for Show ID %s Season %d, got 0", showID, seasonNumber)
   }

   // Verify they are sorted
   if len(episodes) > 1 && episodes[0].EpisodeNumber > episodes[1].EpisodeNumber {
      t.Errorf("Episodes were not sorted properly")
   }

   t.Logf("Successfully retrieved %d episodes from live server. First episode edit ID: %s",
      len(episodes), episodes[0].EditID)
}

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
