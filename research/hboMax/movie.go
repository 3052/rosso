// movie.go
package hboMax

import (
   "encoding/json"
   "fmt"
   "io"
   "net/http"
)

// Internal structures for parsing the Movie API response
type movieResponse struct {
   Included []struct {
      ID         string `json:"id"`
      Type       string `json:"type"`
      Attributes struct {
         VideoType string `json:"videoType"`
      } `json:"attributes"`
      Relationships struct {
         Edit struct {
            Data struct {
               ID string `json:"id"`
            } `json:"data"`
         } `json:"edit"`
      } `json:"relationships"`
   } `json:"included"`
}

// GetMovieEditID fetches the CMS data for a movie ID and extracts the underlying 'Edit' ID.
func (c *Client) GetMovieEditID(movieRouteID string) (string, error) {
   endpoint := fmt.Sprintf("/cms/routes/movie/%s?include=default&decorators=viewingHistory,isFavorite,contentAction,badges&page[items.size]=10", movieRouteID)

   req, err := c.newRequest("GET", endpoint)
   if err != nil {
      return "", err
   }

   resp, err := c.HTTPClient.Do(req)
   if err != nil {
      return "", err
   }
   defer resp.Body.Close()

   if resp.StatusCode != http.StatusOK {
      return "", fmt.Errorf("API returned non-200 status code: %d", resp.StatusCode)
   }

   bodyBytes, err := io.ReadAll(resp.Body)
   if err != nil {
      return "", err
   }

   var maxResp movieResponse
   if err := json.Unmarshal(bodyBytes, &maxResp); err != nil {
      return "", err
   }

   for _, item := range maxResp.Included {
      // Identify the primary video entity for the movie
      if item.Type == "video" && item.Attributes.VideoType == "MOVIE" {
         return item.Relationships.Edit.Data.ID, nil
      }
   }

   return "", fmt.Errorf("edit ID not found for movie route ID: %s", movieRouteID)
}
