package hboMax

import (
   "encoding/json"
   "fmt"
   "net/http"
   "net/url"
   "strings"
)

func (e *Error) Error() string {
   var data strings.Builder
   // 1. print code
   data.WriteString("code = ")
   data.WriteString(e.Code)
   // 2, 3, 4. if detail print detail, if message print message, if both print
   // one
   if e.Detail != "" {
      data.WriteString("\ndetail = ")
      data.WriteString(e.Detail)
   } else if e.Message != "" {
      data.WriteString("\nmessage = ")
      data.WriteString(e.Message)
   }
   return data.String()
}

type Error struct {
   Code    string // 2026-04-10
   Detail  string // 2026-04-10
   Message string // 2026-04-10
}

func (l Login) entity_request(endpoint *url.URL) ([]*Entity, error) {
   // Scheme
   endpoint.Scheme = "https"
   // Host
   endpoint.Host = "default.prd.api.hbomax.com"
   // RawQuery
   queryParams := endpoint.Query()
   queryParams.Set("include", "default")
   endpoint.RawQuery = queryParams.Encode()
   req := http.Request{
      URL:    endpoint,
      Header: http.Header{},
   }
   req.Header.Set("authorization", "Bearer "+l.Token)
   resp, err := http.DefaultClient.Do(&req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   var result struct {
      Errors   []Error
      Included []*Entity `json:"included"`
   }
   err = json.NewDecoder(resp.Body).Decode(&result)
   if err != nil {
      return nil, err
   }
   if len(result.Errors) >= 1 {
      return nil, &result.Errors[0]
   }
   return result.Included, nil
}

type Login struct {
   Token string
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

// Resource represents a relationship pointer in the JSON:API graph
type Resource struct {
   Id   string
   Type string
}

// Entity represents a single unified node in the Max API response
type Entity struct {
   Attributes struct {
      Name          string
      Alias         string
      ShowType      string
      VideoType     string
      MaterialType  string
      Description   string
      SeasonNumber  int
      EpisodeNumber int
      AirDate       string
   }
   Id            string
   Relationships struct {
      Edit struct {
         Data Resource
      }
      Items struct {
         Data []Resource
      }
      Show struct {
         Data Resource
      }
      Video struct {
         Data Resource
      }
   }
   Type string
}
