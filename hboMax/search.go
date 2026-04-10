package hboMax

import (
   "fmt"
   "net/url"
)

// Search queries the API and returns the root entity slice
func (l Login) Search(query string) ([]*Entity, error) {
   queryParams := url.Values{}
   queryParams.Set("page[items.size]", "10")
   queryParams.Set("contentFilter[query]", query)
   parsedURL := &url.URL{
      Path:     "/cms/routes/search/result",
      RawQuery: queryParams.Encode(),
   }
   return l.getEntities(parsedURL)
}

// GetSearchResults parses an entity slice into an ordered list of matching
// media entities
func GetSearchResults(entities []*Entity) ([]*Entity, error) {
   entitiesMap := make(map[string]*Entity)
   for _, entity := range entities {
      entitiesMap[entity.ID] = entity
   }

   var searchResultsCollection *Entity
   for _, entity := range entities {
      if entity.Type == "collection" && entity.Attributes.Alias == "search-page-rail-results" {
         searchResultsCollection = entity
         break
      }
   }

   if searchResultsCollection == nil {
      return nil, fmt.Errorf("could not find the search results collection in the response payload")
   }

   var results []*Entity
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

      // Append the actual show/movie entity
      results = append(results, mediaEntity)
   }

   return results, nil
}
