package main

import (
   "encoding/json"
   "fmt"
   "io"
   "log"
   "net/http"
)

// MaxResponse defines the structure needed to extract the edit ID from the JSON payload.
type MaxResponse struct {
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

func main() {
   reqURL := "https://default.any-emea.prd.api.hbomax.com/cms/routes/movie/bebe611d-8178-481a-a4f2-de743b5b135a?include=default&decorators=viewingHistory,isFavorite,contentAction,badges&page[items.size]=10"

   req, err := http.NewRequest("GET", reqURL, nil)
   if err != nil {
      log.Fatalf("Failed to create request: %v", err)
   }

   // 1. Set required static headers
   req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:140.0) Gecko/20100101 Firefox/140.0")
   req.Header.Set("Accept", "*/*")
   req.Header.Set("Accept-Language", "en-US,en;q=0.5")
   req.Header.Set("Referer", "https://play.hbomax.com/")
   req.Header.Set("X-Device-Info", "hbomax/6.17.1 (desktop/desktop; Windows/NT 10.0; f681564c-1be5-4495-882b-6efc06cd8a9d/da0cdd94-5a39-42ef-aa68-54cbc1b852c3)")
   req.Header.Set("X-Disco-Client", "WEB:NT 10.0:hbomax:6.17.1")
   req.Header.Set("X-Disco-Params", "realm=bolt,bid=beam,features=ar")

   // 2. Set Authentication headers
   stToken := "eyJhbGciOiJSUzI1NiJ9.eyJqdGkiOiJ0b2tlbi1hMmNlZThjMy0zNGNhLTQ0YTEtYjM4NC04YzIzOWNkZmQxZWQiLCJpc3MiOiJmcGEtaXNzdWVyIiwic3ViIjoiVVNFUklEOmJvbHQ6MGQ0NWNjZjgtYjRhMi00MTQ3LWJiZWItYzdiY2IxNDBmMzgyIiwiaWF0IjoxNzc1NjE5MTE5LCJleHAiOjIwOTA5NzkxMTksInR5cGUiOiJBQ0NFU1NfVE9LRU4iLCJzdWJkaXZpc2lvbiI6ImJlYW1fZW1lYSIsInNjb3BlIjoiZGVmYXVsdCIsImlpZCI6ImJlMzI5MzdhLTU3MWEtNDAzMC1hZWIyLTQ1MWViZjI3M2M5YiIsInZlcnNpb24iOiJ2NCIsImFub255bW91cyI6ZmFsc2UsImRldmljZUlkIjoiZjY4MTU2NGMtMWJlNS00NDk1LTg4MmItNmVmYzA2Y2Q4YTlkIn0.kkxM9-egjkpxnz2fSft9G1cQMdfFh9qK8_DHTk2D7Zb43FpORAkUbU92X7o-AMZxPl9pQfDlsE4KWmJHIB3vQUAC5WJmJHUDC2jc7nFYvKhDJfFLDcZD7Jc6TvpNrIYkbhP0gfF_lAxImYfoUFAQx9XzGWFiVfGe1Sy8lalVMwF-nQBdNSPxGijg1IAp-8Nt4xIScM3RScJDaJ7LqQzpNc4p9vK1l68oVUXA-NsE1RpB6caS7AucluygtjVSIGqtLE2HNDMQhJijPdCvYjRmNrQq30Ke_6tC6ezGIj5OD3Z2Sm4lJ0gFdzMZu_MggPUyadEbbK2LDI9nTU5qch1RYw"
   req.Header.Set("Cookie", fmt.Sprintf("st=%s", stToken))

   // Execute the HTTP Request
   client := &http.Client{}
   resp, err := client.Do(req)
   if err != nil {
      log.Fatalf("Failed to execute request: %v", err)
   }
   defer resp.Body.Close()

   if resp.StatusCode != 200 {
      log.Printf("Warning: Received non-200 status code: %d", resp.StatusCode)
   }

   // Read the JSON response body
   bodyBytes, err := io.ReadAll(resp.Body)
   if err != nil {
      log.Fatalf("Failed to read response body: %v", err)
   }

   // Unmarshal into our Go struct
   var maxResp MaxResponse
   if err := json.Unmarshal(bodyBytes, &maxResp); err != nil {
      log.Fatalf("Failed to parse JSON: %v", err)
   }

   // Extract the Edit ID
   var editID string
   for _, item := range maxResp.Included {
      // Identify the primary video entity for the movie
      if item.Type == "video" && item.Attributes.VideoType == "MOVIE" {
         editID = item.Relationships.Edit.Data.ID
         break
      }
   }

   // Output result
   if editID != "" {
      fmt.Printf("Successfully found Edit ID: %s\n", editID)
   } else {
      fmt.Println("Edit ID not found in the response.")
   }
}
