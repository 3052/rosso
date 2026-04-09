package hboMax

import (
   "fmt"
   "net/url"
)

// SearchResult represents a normalized search item.
type SearchResult struct {
   Name      string
   MediaType string
}

// GetSearchResults parses an entity slice into an ordered list of SearchResults.
func GetSearchResults(entities []Entity) ([]SearchResult, error) {
   entitiesMap := make(map[string]Entity)
   for _, e := range entities {
      entitiesMap[e.ID] = e
   }

   var searchResultsCollection Entity
   found := false
   for _, e := range entities {
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

// Search queries the API and returns the root entity slice.
func (c *Client) Search(query string) ([]Entity, error) {
   endpoint := fmt.Sprintf("/cms/routes/search/result?include=default&decorators=viewingHistory,isFavorite,contentAction,badges&page[items.size]=10&contentFilter[query]=%s", url.QueryEscape(query))
   return c.getEntities(endpoint)
}
