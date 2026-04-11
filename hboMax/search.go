package hboMax

import (
   "fmt"
   "net/url"
)

// Entity represents a single unified node in the Max API response.
type Entity struct {
   Attributes struct {
      Name          string `json:"name"`
      Alias         string `json:"alias"`
      ShowType      string `json:"showType"`
      VideoType     string `json:"videoType"`
      MaterialType  string `json:"materialType"`
      Description   string `json:"description"`
      SeasonNumber  int    `json:"seasonNumber"`
      EpisodeNumber int    `json:"episodeNumber"`
      AirDate       string `json:"airDate"`
   }
   Id            string `json:"id"`
   Relationships struct {
      Edit struct {
         Data Resource `json:"data"`
      } `json:"edit"`
      Items struct {
         Data []Resource `json:"data"`
      } `json:"items"`
      Show struct {
         Data Resource `json:"data"`
      } `json:"show"`
      Video struct {
         Data Resource `json:"data"`
      } `json:"video"`
   } `json:"relationships"`
   Type string `json:"type"`
}

func SearchResults(entities []*Entity) ([]*Entity, error) {
   entitiesMap := make(map[string]*Entity)
   for _, entity := range entities {
      entitiesMap[entity.Id] = entity
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
      colItem, exists := entitiesMap[itemRes.Id]
      if !exists {
         continue
      }

      targetId := colItem.Relationships.Show.Data.Id
      if targetId == "" {
         targetId = colItem.Relationships.Video.Data.Id
      }

      if targetId == "" {
         continue
      }

      mediaEntity, exists := entitiesMap[targetId]
      if !exists {
         continue
      }

      // Append the actual show/movie entity
      results = append(results, mediaEntity)
   }
   return results, nil
}

func (l Login) SearchRequest(query string) ([]*Entity, error) {
   queryParams := url.Values{}
   queryParams.Set("contentFilter[query]", query)
   parsedUrl := &url.URL{
      Path:     "/cms/routes/search/result",
      RawQuery: queryParams.Encode(),
   }
   return l.entity_request(parsedUrl)
}
