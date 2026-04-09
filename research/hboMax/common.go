package hboMax

import (
   "encoding/json"
   "fmt"
   "io"
   "net/http"
   "net/url"
   "strings"
)

// Resource represents a relationship pointer in the JSON:API graph.
type Resource struct {
   ID   string `json:"id"`
   Type string `json:"type"`
}

// Attributes holds the metadata properties for a media entity.
type Attributes struct {
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

// Entity represents a single unified node in the Max API response.
type Entity struct {
   ID            string     `json:"id"`
   Type          string     `json:"type"`
   Attributes    Attributes `json:"attributes"`
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
func (c *Client) getEntities(u *url.URL) ([]*Entity, error) {
   // Inject generic query parameters required by all endpoints
   q := u.Query()
   q.Set("include", "default")
   q.Set("decorators", "viewingHistory,isFavorite,contentAction,badges")
   u.RawQuery = q.Encode()

   // Resolve the relative URL against the base URL
   baseURL, err := url.Parse(c.BaseURL)
   if err != nil {
      return nil, err
   }
   reqURL := baseURL.ResolveReference(u).String()

   req, err := http.NewRequest("GET", reqURL, nil)
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
      Included []*Entity `json:"included"`
   }
   if err := json.Unmarshal(bodyBytes, &rootResponse); err != nil {
      return nil, err
   }

   return rootResponse.Included, nil
}

// String implements the fmt.Stringer interface to provide a clean visual output for the Entity.
func (e *Entity) String() string {
   var b strings.Builder

   // 1. print episode number if material type is episode
   if e.Attributes.MaterialType == "EPISODE" {
      fmt.Fprintf(&b, "Episode: %d\n", e.Attributes.EpisodeNumber)
   }

   // 2. print attributes name
   fmt.Fprintf(&b, "Name: %s\n", e.Attributes.Name)

   // 3 & 4. print edit ID if type is video, otherwise print ID
   if e.Type == "video" {
      fmt.Fprintf(&b, "Edit ID: %s\n", e.Relationships.Edit.Data.ID)
   } else {
      fmt.Fprintf(&b, "ID: %s\n", e.ID)
   }

   // 5. print either show type or video type
   if e.Attributes.ShowType != "" {
      fmt.Fprintf(&b, "Show Type: %s\n", e.Attributes.ShowType)
   } else if e.Attributes.VideoType != "" {
      fmt.Fprintf(&b, "Video Type: %s\n", e.Attributes.VideoType)
   }

   return strings.TrimSpace(b.String())
}
