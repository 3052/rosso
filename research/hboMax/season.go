package hboMax

import (
   "cmp"
   "fmt"
   "net/url"
   "slices"
)

// GetEpisodes filters the entity slice for episodes and sorts them chronologically.
func GetEpisodes(entities []*Entity) []*Entity {
   var episodes []*Entity
   for _, item := range entities {
      if item.Type == "video" && item.Attributes.MaterialType == "EPISODE" {
         episodes = append(episodes, item)
      }
   }

   // Sort episodes by EpisodeNumber using the modern slices.SortFunc
   slices.SortFunc(episodes, func(a, b *Entity) int {
      return cmp.Compare(a.Attributes.EpisodeNumber, b.Attributes.EpisodeNumber)
   })

   return episodes
}

// GetSeasonEpisodes fetches all entities for a given show ID and season number.
func (c *Client) GetSeasonEpisodes(showID string, seasonNumber int) ([]*Entity, error) {
   u, err := url.Parse("/cms/collections/generic-show-page-rail-episodes-tabbed-content")
   if err != nil {
      return nil, err
   }

   q := u.Query()
   q.Set("pf[show.id]", showID)
   q.Set("pf[seasonNumber]", fmt.Sprintf("%d", seasonNumber))
   u.RawQuery = q.Encode()

   return c.getEntities(u)
}
