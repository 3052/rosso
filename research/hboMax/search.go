// search.go
package hboMax

import (
   "encoding/json"
   "fmt"
   "io"
   "net/http"
   "net/url"
)

// Client handles communication with the HBO Max API.
type Client struct {
   BaseURL    string
   Token      string
   HTTPClient *http.Client
}

// NewClient creates a new API client with the given st token.
func NewClient(token string) *Client {
   return &Client{
      BaseURL:    "https://default.any-emea.prd.api.hbomax.com",
      Token:      token,
      HTTPClient: &http.Client{},
   }
}

// newRequest sets up the boilerplate headers required by the API.
func (c *Client) newRequest(method, endpoint string) (*http.Request, error) {
   req, err := http.NewRequest(method, c.BaseURL+endpoint, nil)
   if err != nil {
      return nil, err
   }

   req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:140.0) Gecko/20100101 Firefox/140.0")
   req.Header.Set("Accept", "application/json")
   req.Header.Set("Referer", "https://play.hbomax.com/")
   req.Header.Set("X-Device-Info", "hbomax/6.17.1 (desktop/desktop; Windows/NT 10.0; f681564c-1be5-4495-882b-6efc06cd8a9d/da0cdd94-5a39-42ef-aa68-54cbc1b852c3)")
   req.Header.Set("X-Disco-Client", "WEB:NT 10.0:hbomax:6.17.1")
   req.Header.Set("X-Disco-Params", "realm=bolt,bid=beam,features=ar")
   req.Header.Set("Cookie", fmt.Sprintf("st=%s", c.Token))

   return req, nil
}

// SearchResult represents a normalized search item.
type SearchResult struct {
   Name      string
   MediaType string
}

// Internal structures to map the JSON:API graph for Search
type searchResource struct {
   ID   string `json:"id"`
   Type string `json:"type"`
}

type searchEntity struct {
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
         Data []searchResource `json:"data"`
      } `json:"items"`
      Show struct {
         Data searchResource `json:"data"`
      } `json:"show"`
      Video struct {
         Data searchResource `json:"data"`
      } `json:"video"`
   } `json:"relationships"`
}

type searchResponse struct {
   Included []searchEntity `json:"included"`
}

// Search queries the API and parses the JSON graph into an ordered list of SearchResults.
func (c *Client) Search(query string) ([]SearchResult, error) {
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

   var apiResp searchResponse
   if err := json.Unmarshal(bodyBytes, &apiResp); err != nil {
      return nil, err
   }

   entitiesMap := make(map[string]searchEntity)
   for _, e := range apiResp.Included {
      entitiesMap[e.ID] = e
   }

   var searchResultsCollection searchEntity
   found := false
   for _, e := range apiResp.Included {
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
