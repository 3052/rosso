package hboMax

import (
   "fmt"
)

// GetSeasonEpisodes fetches all entities for a given show ID and season number.
func (c *Client) GetSeasonEpisodes(showID string, seasonNumber int) ([]*Entity, error) {
   // Use the generic collection alias for the episodes tab
   endpoint := fmt.Sprintf("/cms/collections/generic-show-page-rail-episodes-tabbed-content?include=default&decorators=viewingHistory,isFavorite,contentAction,badges&pf[show.id]=%s&pf[seasonNumber]=%d", showID, seasonNumber)
   return c.getEntities(endpoint)
}
