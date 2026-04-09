package hboMax

import (
   "encoding/json"
   "fmt"
   "io"
   "net/http"
   "net/url"
)

// SearchResult represents a normalized search item.
type SearchResult struct {
   Name      string
   MediaType string
}

// SearchResource represents a relationship pointer in the JSON:API graph.
type SearchResource struct {
   ID   string `json:"id"`
   Type string `json:"type"`
}

// SearchEntity represents a single node in the API response.
type SearchEntity struct {
   ID         string `json:"id"`
   Type       string `json:"type"`
   Attributes struct {
      Alias     string `json:"alias"`
      Name      string `json:"name"`
      ShowType  string `json:"showType"`
      VideoType string `json:"videoType"`
   } `json:"attributes"`
   Relationships struct {
      Items struct {
         Data []SearchResource `json:"data"`
      } `json:"items"`
      Show struct {
         Data SearchResource `json:"data"`
      } `json:"show"`
      Video struct {
         Data SearchResource `json:"data"`
      } `json:"video"`
   } `json:"relationships"`
}

// SearchResponse represents the root JSON structure returned by the API.
type SearchResponse struct {
   Included []SearchEntity `json:"included"`
}

// GetResults parses the JSON graph into an ordered list of SearchResults.
func (sr *SearchResponse) GetResults() ([]SearchResult, error) {
   entitiesMap := make(map[string]SearchEntity)
   for _, e := range sr.Included {
      entitiesMap[e.ID] = e
   }

   var searchResultsCollection SearchEntity
   found := false
   for _, e := range sr.Included {
      if e.Type == "collection" && e.Attributes.Alias == "search-page-rail-results" {
         searchResultsCollection = e
         found = true
         break
      }
   }

   if !found {
      return nil, fmt.Errorf("could not find the search results collection in the response payload")
   }

   var results []SearchResult
   for _, itemRes := range searchResultsCollection.Relationships.Items.Data {
      colItem, exists := entitiesMap[itemRes.ID]
      if !exists {
         continue
      }

      targetID := colItem.Relationships.Show.Data.ID
      if targetID == "" {
         targetID = colItem.Relationships.Video.Data.ID
      }

      if targetID == "" {
         continue
      }

      mediaEntity, exists := entitiesMap[targetID]
      if !exists {
         continue
      }

      mediaType := mediaEntity.Attributes.ShowType
      if mediaType == "" {
         mediaType = mediaEntity.Attributes.VideoType
      }

      results = append(results, SearchResult{
         Name:      mediaEntity.Attributes.Name,
         MediaType: mediaType,
      })
   }

   return results, nil
}

// Search queries the API and returns the parsed JSON response.
func (c *Client) Search(query string) (*SearchResponse, error) {
   endpoint := fmt.Sprintf("/cms/routes/search/result?include=default&decorators=viewingHistory,isFavorite,contentAction,badges&page[items.size]=10&contentFilter[query]=%s", url.QueryEscape(query))

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

   var apiResp SearchResponse
   if err := json.Unmarshal(bodyBytes, &apiResp); err != nil {
      return nil, err
   }

   return &apiResp, nil
}
