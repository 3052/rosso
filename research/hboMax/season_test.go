package hboMax

import (
   "testing"
)

func TestGetSeasonEpisodesLive(t *testing.T) {
   client := NewClient(testToken)

   showID := "4ffd33c9-e0d6-4cd6-bd13-34c266c79be0"
   seasonNumber := 2

   // Get raw entities
   entities, err := client.GetSeasonEpisodes(showID, seasonNumber)
   if err != nil {
      t.Fatalf("Live API request failed: %v", err)
   }

   // Extract formatted episodes
   episodes := GetVideos(entities)

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
      t.Logf("Episode %d: %s", ep.Attributes.EpisodeNumber, ep.Attributes.Name)
      t.Logf("Edit ID:   %s", ep.Relationships.Edit.Data.ID)
      t.Logf("Air Date:  %s", ep.Attributes.AirDate)
      t.Logf("Summary:   %s", ep.Attributes.Description)
      t.Log("--------------------------------------------------")
   }
}
