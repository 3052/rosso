package hboMax

import (
   "encoding/json"
   "fmt"
   "io"
   "net/http"
)

// Resource represents a relationship pointer in the JSON:API graph.
type Resource struct {
   ID   string `json:"id"`
   Type string `json:"type"`
}

// Entity represents a single unified node in the Max API response.
// It combines the attributes and relationships needed across Search, Season, and Movie endpoints.
type Entity struct {
   ID         string `json:"id"`
   Type       string `json:"type"`
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
      Edit struct {
         Data Resource `json:"data"`
      } `json:"edit"`
   } `json:"relationships"`
}

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

// getEntities is a shared internal method that hits an endpoint and returns the extracted JSON:API entities.
func (c *Client) getEntities(endpoint string) ([]Entity, error) {
   req, err := http.NewRequest("GET", c.BaseURL+endpoint, nil)
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

   var rootResponse struct {
      Included []Entity `json:"included"`
   }
   if err := json.Unmarshal(bodyBytes, &rootResponse); err != nil {
      return nil, err
   }

   return rootResponse.Included, nil
}
