package hboMax

import (
   "fmt"
)

// GetMovie fetches the CMS data for a movie ID and returns the parsed entities.
func (c *Client) GetMovie(movieRouteID string) ([]*Entity, error) {
   endpoint := fmt.Sprintf("/cms/routes/movie/%s?include=default&decorators=viewingHistory,isFavorite,contentAction,badges&page[items.size]=10", movieRouteID)
   return c.getEntities(endpoint)
}
