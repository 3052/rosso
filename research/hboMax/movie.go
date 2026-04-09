package hboMax

import (
   "fmt"
   "net/url"
)

// GetMovies filters the entity slice for primary movie video entities.
func GetMovies(entities []*Entity) []*Entity {
   var movies []*Entity
   for _, item := range entities {
      // Identify the primary video entity for the movie
      if item.Type == "video" && item.Attributes.VideoType == "MOVIE" {
         movies = append(movies, item)
      }
   }
   return movies
}

// GetMovie fetches the CMS data for a movie ID and returns the parsed entities.
func (c *Client) GetMovie(movieRouteID string) ([]*Entity, error) {
   u, err := url.Parse(fmt.Sprintf("/cms/routes/movie/%s", movieRouteID))
   if err != nil {
      return nil, err
   }

   q := u.Query()
   q.Set("page[items.size]", "10")
   u.RawQuery = q.Encode()

   return c.getEntities(u)
}
