package hboMax

import (
   "testing"
)

func TestGetSeasonEpisodesLive(t *testing.T) {
   client := NewClient(testToken)

   showID := "4ffd33c9-e0d6-4cd6-bd13-34c266c79be0"
   seasonNumber := 2

   entities, err := client.GetSeasonEpisodes(showID, seasonNumber)
   if err != nil {
      t.Fatalf("Live API request failed: %v", err)
   }

   episodes := GetEpisodes(entities)

   if len(episodes) == 0 {
      t.Fatalf("Expected episodes for Show ID %s Season %d, got 0", showID, seasonNumber)
   }

   if len(episodes) > 1 && episodes[0].Attributes.EpisodeNumber > episodes[1].Attributes.EpisodeNumber {
      t.Errorf("Episodes were not sorted properly")
   }

   t.Log("==================================================")
   t.Logf(" Found %d Episodes for Season %d", len(episodes), seasonNumber)
   t.Log("==================================================")
   for _, ep := range episodes {
      t.Log(ep.String())
      t.Log("--------------------------------------------------")
   }
}
