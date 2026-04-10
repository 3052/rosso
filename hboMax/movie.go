package hboMax

import (
   "net/url"
)

// GetMovie fetches the CMS data for a movie ID and returns the parsed entities
func (l Login) GetMovie(movieRouteID string) ([]*Entity, error) {
   queryParams := url.Values{}
   queryParams.Set("page[items.size]", "1")
   parsedURL := &url.URL{
      Path:     "/cms/routes/movie/" + movieRouteID,
      RawQuery: queryParams.Encode(),
   }
   return l.getEntities(parsedURL)
}

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
