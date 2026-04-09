package main

import (
   "encoding/json"
   "fmt"
   "io"
   "log"
   "net/http"
   "sort"
)

// Define the structs needed to parse the JSON:API response format
type HBOApiResponse struct {
   Included []IncludedData `json:"included"`
}

type IncludedData struct {
   Id            string        `json:"id"`
   Type          string        `json:"type"`
   Attributes    Attributes    `json:"attributes"`
   Relationships Relationships `json:"relationships"`
}

type Attributes struct {
   MaterialType  string `json:"materialType"`
   Name          string `json:"name"`
   Description   string `json:"description"`
   SeasonNumber  int    `json:"seasonNumber"`
   EpisodeNumber int    `json:"episodeNumber"`
   AirDate       string `json:"airDate"`
}

type Relationships struct {
   Edit RelationshipEdit `json:"edit"`
}

type RelationshipEdit struct {
   Data RelationshipData `json:"data"`
}

type RelationshipData struct {
   Id   string `json:"id"`
   Type string `json:"type"`
}

func main() {
   // The collection ID '227084608563650952176059252419027445293' represents the "Season Tabbed Content" UI rail.
   // You can change pf[seasonNumber]=2 to target different seasons.
   url := "https://default.any-emea.prd.api.hbomax.com/cms/collections/227084608563650952176059252419027445293?include=default&decorators=viewingHistory,isFavorite,contentAction,badges&pf[show.id]=4ffd33c9-e0d6-4cd6-bd13-34c266c79be0&pf[seasonNumber]=2"

   req, err := http.NewRequest("GET", url, nil)
   if err != nil {
      log.Fatalf("Error creating request: %v", err)
   }

   // Standard Headers
   req.Header.Set("accept", "application/json")
   req.Header.Set("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:140.0) Gecko/20100101 Firefox/140.0")

   // Max/Discovery Specific Headers
   req.Header.Set("x-device-info", "hbomax/6.17.1 (desktop/desktop; Windows/NT 10.0; f681564c-1be5-4495-882b-6efc06cd8a9d/da0cdd94-5a39-42ef-aa68-54cbc1b852c3)")
   req.Header.Set("x-disco-client", "WEB:NT 10.0:hbomax:6.17.1")
   req.Header.Set("x-disco-params", "realm=bolt,bid=beam,features=ar")

   // ONLY the 'st' cookie is provided for authentication
   stToken := "eyJhbGciOiJSUzI1NiJ9.eyJqdGkiOiJ0b2tlbi1hMmNlZThjMy0zNGNhLTQ0YTEtYjM4NC04YzIzOWNkZmQxZWQiLCJpc3MiOiJmcGEtaXNzdWVyIiwic3ViIjoiVVNFUklEOmJvbHQ6MGQ0NWNjZjgtYjRhMi00MTQ3LWJiZWItYzdiY2IxNDBmMzgyIiwiaWF0IjoxNzc1NjE5MTE5LCJleHAiOjIwOTA5NzkxMTksInR5cGUiOiJBQ0NFU1NfVE9LRU4iLCJzdWJkaXZpc2lvbiI6ImJlYW1fZW1lYSIsInNjb3BlIjoiZGVmYXVsdCIsImlpZCI6ImJlMzI5MzdhLTU3MWEtNDAzMC1hZWIyLTQ1MWViZjI3M2M5YiIsInZlcnNpb24iOiJ2NCIsImFub255bW91cyI6ZmFsc2UsImRldmljZUlkIjoiZjY4MTU2NGMtMWJlNS00NDk1LTg4MmItNmVmYzA2Y2Q4YTlkIn0.kkxM9-egjkpxnz2fSft9G1cQMdfFh9qK8_DHTk2D7Zb43FpORAkUbU92X7o-AMZxPl9pQfDlsE4KWmJHIB3vQUAC5WJmJHUDC2jc7nFYvKhDJfFLDcZD7Jc6TvpNrIYkbhP0gfF_lAxImYfoUFAQx9XzGWFiVfGe1Sy8lalVMwF-nQBdNSPxGijg1IAp-8Nt4xIScM3RScJDaJ7LqQzpNc4p9vK1l68oVUXA-NsE1RpB6caS7AucluygtjVSIGqtLE2HNDMQhJijPdCvYjRmNrQq30Ke_6tC6ezGIj5OD3Z2Sm4lJ0gFdzMZu_MggPUyadEbbK2LDI9nTU5qch1RYw"
   req.Header.Set("Cookie", fmt.Sprintf("st=%s", stToken))

   client := &http.Client{}
   resp, err := client.Do(req)
   if err != nil {
      log.Fatalf("Request failed: %v", err)
   }
   defer resp.Body.Close()

   if resp.StatusCode != http.StatusOK {
      log.Fatalf("API returned non-200 status code: %d", resp.StatusCode)
   }

   body, err := io.ReadAll(resp.Body)
   if err != nil {
      log.Fatalf("Failed to read body: %v", err)
   }

   // Parse the JSON Response
   var apiResponse HBOApiResponse
   if err := json.Unmarshal(body, &apiResponse); err != nil {
      log.Fatalf("Failed to unmarshal JSON: %v", err)
   }

   // Extract episodes from the "included" array
   var episodes []IncludedData
   for _, item := range apiResponse.Included {
      if item.Type == "video" && item.Attributes.MaterialType == "EPISODE" {
         episodes = append(episodes, item)
      }
   }

   // Sort episodes by EpisodeNumber just in case the API returned them out of order
   sort.Slice(episodes, func(i, j int) bool {
      return episodes[i].Attributes.EpisodeNumber < episodes[j].Attributes.EpisodeNumber
   })

   // Print the output
   fmt.Println("==================================================")
   fmt.Printf(" Found %d Episodes for Season %d\n", len(episodes), 2)
   fmt.Println("==================================================")

   for _, ep := range episodes {
      fmt.Printf("Episode %d: %s\n", ep.Attributes.EpisodeNumber, ep.Attributes.Name)
      fmt.Printf("Video ID:  %s\n", ep.Id)
      fmt.Printf("Edit ID:   %s\n", ep.Relationships.Edit.Data.Id)
      fmt.Printf("Air Date:  %s\n", ep.Attributes.AirDate)
      fmt.Printf("Summary:   %s\n", ep.Attributes.Description)
      fmt.Println("--------------------------------------------------")
   }
}
