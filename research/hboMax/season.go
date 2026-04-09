package hboMax

import (
   "fmt"
   "sort"
)

// Episode represents a parsed episode from a season.
type Episode struct {
   VideoID       string
   EditID        string
   Name          string
   Description   string
   SeasonNumber  int
   EpisodeNumber int
   AirDate       string
}

// GetEpisodes filters the entity slice for episodes and sorts them chronologically.
func GetEpisodes(entities []Entity) []Episode {
   var episodes []Episode
   for _, item := range entities {
      if item.Type == "video" && item.Attributes.MaterialType == "EPISODE" {
         episodes = append(episodes, Episode{
            VideoID:       item.ID,
            EditID:        item.Relationships.Edit.Data.ID,
            Name:          item.Attributes.Name,
            Description:   item.Attributes.Description,
            SeasonNumber:  item.Attributes.SeasonNumber,
            EpisodeNumber: item.Attributes.EpisodeNumber,
            AirDate:       item.Attributes.AirDate,
         })
      }
   }

   // Sort episodes by EpisodeNumber just in case the API returned them out of order
   sort.Slice(episodes, func(i, j int) bool {
      return episodes[i].EpisodeNumber < episodes[j].EpisodeNumber
   })

   return episodes
}

// GetSeasonEpisodes fetches all entities for a given show ID and season number.
func (c *Client) GetSeasonEpisodes(showID string, seasonNumber int) ([]Entity, error) {
   // The collection ID '227084608563650952176059252419027445293' represents the "Season Tabbed Content" UI rail.
   endpoint := fmt.Sprintf("/cms/collections/227084608563650952176059252419027445293?include=default&decorators=viewingHistory,isFavorite,contentAction,badges&pf[show.id]=%s&pf[seasonNumber]=%d", showID, seasonNumber)
   return c.getEntities(endpoint)
}
