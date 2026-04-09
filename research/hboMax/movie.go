package hboMax

import (
   "fmt"
)

// GetEditID extracts the underlying 'Edit' ID from the movie entities.
func (entities Entities) GetEditID() (string, error) {
   for _, item := range entities {
      // Identify the primary video entity for the movie
      if item.Type == "video" && item.Attributes.VideoType == "MOVIE" {
         return item.Relationships.Edit.Data.ID, nil
      }
   }

   return "", fmt.Errorf("edit ID not found in the movie entities")
}

// GetMovie fetches the CMS data for a movie ID and returns the parsed entities.
func (c *Client) GetMovie(movieRouteID string) (Entities, error) {
   endpoint := fmt.Sprintf("/cms/routes/movie/%s?include=default&decorators=viewingHistory,isFavorite,contentAction,badges&page[items.size]=10", movieRouteID)
   return c.getEntities(endpoint)
}
