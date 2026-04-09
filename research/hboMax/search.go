package main

import (
   "encoding/json"
   "fmt"
   "io"
   "log"
   "net/http"
)

// Resource represents a relationship pointer in the JSON:API graph.
type Resource struct {
   ID   string `json:"id"`
   Type string `json:"type"`
}

// Entity represents a single node (Show, Video, Collection, CollectionItem) in the API response.
type Entity struct {
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
         Data []Resource `json:"data"`
      } `json:"items"`
      Show struct {
         Data Resource `json:"data"`
      } `json:"show"`
      Video struct {
         Data Resource `json:"data"`
      } `json:"video"`
   } `json:"relationships"`
}

// MaxResponse represents the root JSON structure returned by the Max API.
type MaxResponse struct {
   Included []Entity `json:"included"`
}

func main() {
   reqURL := "https://default.any-emea.prd.api.hbomax.com/cms/routes/search/result?include=default&decorators=viewingHistory,isFavorite,contentAction,badges&page[items.size]=10&contentFilter[query]=marnie"

   req, err := http.NewRequest("GET", reqURL, nil)
   if err != nil {
      log.Fatalf("Failed to create request: %v", err)
   }

   // 1. Set required Headers
   req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:140.0) Gecko/20100101 Firefox/140.0")
   req.Header.Set("Accept", "*/*")
   req.Header.Set("Accept-Language", "en-US,en;q=0.5")
   req.Header.Set("Referer", "https://play.hbomax.com/")
   req.Header.Set("X-Device-Info", "hbomax/6.17.1 (desktop/desktop; Windows/NT 10.0; f681564c-1be5-4495-882b-6efc06cd8a9d/da0cdd94-5a39-42ef-aa68-54cbc1b852c3)")
   req.Header.Set("X-Disco-Client", "WEB:NT 10.0:hbomax:6.17.1")
   req.Header.Set("X-Disco-Params", "realm=bolt,bid=beam,features=ar")
   req.Header.Set("X-Wbd-Device-Consent", "gpc=0")
   req.Header.Set("X-Wbd-Preferred-Language", "en-US,en")

   // Set the Authentication Cookie (Only the "st" cookie is needed)
   stCookie := "st=eyJhbGciOiJSUzI1NiJ9.eyJqdGkiOiJ0b2tlbi1hMmNlZThjMy0zNGNhLTQ0YTEtYjM4NC04YzIzOWNkZmQxZWQiLCJpc3MiOiJmcGEtaXNzdWVyIiwic3ViIjoiVVNFUklEOmJvbHQ6MGQ0NWNjZjgtYjRhMi00MTQ3LWJiZWItYzdiY2IxNDBmMzgyIiwiaWF0IjoxNzc1NjE5MTE5LCJleHAiOjIwOTA5NzkxMTksInR5cGUiOiJBQ0NFU1NfVE9LRU4iLCJzdWJkaXZpc2lvbiI6ImJlYW1fZW1lYSIsInNjb3BlIjoiZGVmYXVsdCIsImlpZCI6ImJlMzI5MzdhLTU3MWEtNDAzMC1hZWIyLTQ1MWViZjI3M2M5YiIsInZlcnNpb24iOiJ2NCIsImFub255bW91cyI6ZmFsc2UsImRldmljZUlkIjoiZjY4MTU2NGMtMWJlNS00NDk1LTg4MmItNmVmYzA2Y2Q4YTlkIn0.kkxM9-egjkpxnz2fSft9G1cQMdfFh9qK8_DHTk2D7Zb43FpORAkUbU92X7o-AMZxPl9pQfDlsE4KWmJHIB3vQUAC5WJmJHUDC2jc7nFYvKhDJfFLDcZD7Jc6TvpNrIYkbhP0gfF_lAxImYfoUFAQx9XzGWFiVfGe1Sy8lalVMwF-nQBdNSPxGijg1IAp-8Nt4xIScM3RScJDaJ7LqQzpNc4p9vK1l68oVUXA-NsE1RpB6caS7AucluygtjVSIGqtLE2HNDMQhJijPdCvYjRmNrQq30Ke_6tC6ezGIj5OD3Z2Sm4lJ0gFdzMZu_MggPUyadEbbK2LDI9nTU5qch1RYw"
   req.Header.Set("Cookie", stCookie)

   // 2. Execute Request
   client := &http.Client{}
   resp, err := client.Do(req)
   if err != nil {
      log.Fatalf("Failed to execute request: %v", err)
   }
   defer resp.Body.Close()

   if resp.StatusCode != 200 {
      log.Printf("Warning: Received non-200 status code: %d", resp.StatusCode)
   }

   bodyBytes, err := io.ReadAll(resp.Body)
   if err != nil {
      log.Fatalf("Failed to read response body: %v", err)
   }

   // 3. Parse JSON
   var apiResp MaxResponse
   if err := json.Unmarshal(bodyBytes, &apiResp); err != nil {
      log.Fatalf("Failed to decode JSON: %v", err)
   }

   // 4. Create a map of all entities for O(1) lookups by ID
   entitiesMap := make(map[string]Entity)
   for _, e := range apiResp.Included {
      entitiesMap[e.ID] = e
   }

   // 5. Find the Collection that holds the Search Results
   var searchResultsCollection Entity
   found := false
   for _, e := range apiResp.Included {
      if e.Type == "collection" && e.Attributes.Alias == "search-page-rail-results" {
         searchResultsCollection = e
         found = true
         break
      }
   }

   if !found {
      log.Fatal("Could not find the search results collection in the response payload.")
   }

   // 6. Traverse the graph and print the Search Results in exact order
   fmt.Println("---------------------------------------------------------")
   fmt.Println("Search Results for 'marnie':")
   fmt.Println("---------------------------------------------------------")

   count := 1
   // The `items` array on the Collection holds pointers to CollectionItems
   for _, itemRes := range searchResultsCollection.Relationships.Items.Data {

      // 6a. Get the CollectionItem
      colItem, exists := entitiesMap[itemRes.ID]
      if !exists {
         continue
      }

      // 6b. A CollectionItem points to either a Show or a Video entity. Determine which one.
      targetID := colItem.Relationships.Show.Data.ID
      if targetID == "" {
         targetID = colItem.Relationships.Video.Data.ID
      }

      if targetID == "" {
         continue
      }

      // 6c. Retrieve the actual Media Entity (Show or Video)
      mediaEntity, exists := entitiesMap[targetID]
      if !exists {
         continue
      }

      // 6d. Extract the Media Type (MOVIE, SERIES, STANDALONE_EVENT, etc.)
      mediaType := mediaEntity.Attributes.ShowType
      if mediaType == "" {
         mediaType = mediaEntity.Attributes.VideoType
      }

      fmt.Printf("%2d. %s [%s]\n", count, mediaEntity.Attributes.Name, mediaType)
      count++
   }
   fmt.Println("---------------------------------------------------------")
}
