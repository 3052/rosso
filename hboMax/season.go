package hboMax

import (
   "cmp"
   "fmt"
   "net/url"
   "slices"
)

// GetSeasonEpisodes fetches all entities for a given show ID and season number
func (l Login) GetSeasonEpisodes(showID string, seasonNumber int) ([]*Entity, error) {
   queryParams := url.Values{}
   queryParams.Set("pf[show.id]", showID)
   queryParams.Set("pf[seasonNumber]", fmt.Sprint(seasonNumber))
   parsedURL := &url.URL{
      Path:     "/cms/collections/generic-show-page-rail-episodes-tabbed-content",
      RawQuery: queryParams.Encode(),
   }
   return l.getEntities(parsedURL)
}

// GetEpisodes filters the entity slice for episodes and sorts them chronologically.
func GetEpisodes(entities []*Entity) []*Entity {
   var episodes []*Entity
   for _, item := range entities {
      if item.Type == "video" && item.Attributes.MaterialType == "EPISODE" {
         episodes = append(episodes, item)
      }
   }
   // Sort episodes by EpisodeNumber using the modern slices.SortFunc
   slices.SortFunc(episodes, func(entityA, entityB *Entity) int {
      return cmp.Compare(entityA.Attributes.EpisodeNumber, entityB.Attributes.EpisodeNumber)
   })
   return episodes
}
