package hboMax

import (
   "encoding/json"
   "fmt"
   "io"
   "net/http"
)

// MovieResponse represents the parsed API response for a movie.
type MovieResponse struct {
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

// GetEditID extracts the underlying 'Edit' ID from the movie response.
func (mr *MovieResponse) GetEditID() (string, error) {
   for _, item := range mr.Included {
      // Identify the primary video entity for the movie
      if item.Type == "video" && item.Attributes.VideoType == "MOVIE" {
         return item.Relationships.Edit.Data.ID, nil
      }
   }

   return "", fmt.Errorf("edit ID not found in the movie response")
}

// GetMovie fetches the CMS data for a movie ID and returns the parsed response.
func (c *Client) GetMovie(movieRouteID string) (*MovieResponse, error) {
   endpoint := fmt.Sprintf("/cms/routes/movie/%s?include=default&decorators=viewingHistory,isFavorite,contentAction,badges&page[items.size]=10", movieRouteID)

   req, err := c.newRequest("GET", endpoint)
   if err != nil {
      return nil, err
   }

   resp, err := c.HTTPClient.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()

   if resp.StatusCode != http.StatusOK {
      return nil, fmt.Errorf("API returned non-200 status code: %d", resp.StatusCode)
   }

   bodyBytes, err := io.ReadAll(resp.Body)
   if err != nil {
      return nil, err
   }

   var movieResp MovieResponse
   if err := json.Unmarshal(bodyBytes, &movieResp); err != nil {
      return nil, err
   }

   return &movieResp, nil
}
