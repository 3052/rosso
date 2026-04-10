package hboMax

import (
   "encoding/json"
   "fmt"
   "net/http"
   "net/url"
   "strings"
)

// getEntities is a shared internal method that hits an endpoint and returns
// the extracted JSON:API entities
func (l Login) getEntities(endpoint *url.URL) ([]*Entity, error) {
   // Scheme
   endpoint.Scheme = "https"
   // Host
   endpoint.Host = "default.prd.api.hbomax.com"
   // RawQuery
   queryParams := endpoint.Query()
   queryParams.Set("include", "default")
   endpoint.RawQuery = queryParams.Encode()
   req := http.Request{
      URL: endpoint,
      Header: http.Header{},
   }
   req.Header.Set("authorization", "Bearer "+l.Data.Attributes.Token)
   resp, err := http.DefaultClient.Do(&req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   if resp.StatusCode != http.StatusOK {
      return nil, fmt.Errorf("API returned non-200 status code: %d", resp.StatusCode)
   }
   var result struct {
      Included []*Entity `json:"included"`
   }
   err = json.NewDecoder(resp.Body).Decode(&result)
   if err != nil {
      return nil, err
   }
   return result.Included, nil
}

// String implements the fmt.Stringer interface to provide a clean visual output for the Entity.
func (e *Entity) String() string {
   var builder strings.Builder

   // 1. print episode number if material type is episode
   if e.Attributes.MaterialType == "EPISODE" {
      fmt.Fprintf(&builder, "Episode: %d\n", e.Attributes.EpisodeNumber)
   }

   // 2. print attributes name
   fmt.Fprintf(&builder, "Name: %s\n", e.Attributes.Name)

   // 3 & 4. print edit ID if type is video, otherwise print ID
   if e.Type == "video" {
      fmt.Fprintf(&builder, "Edit ID: %s\n", e.Relationships.Edit.Data.ID)
   } else {
      fmt.Fprintf(&builder, "ID: %s\n", e.ID)
   }

   // 5. print either show type or video type
   if e.Attributes.ShowType != "" {
      fmt.Fprintf(&builder, "Show Type: %s\n", e.Attributes.ShowType)
   } else if e.Attributes.VideoType != "" {
      fmt.Fprintf(&builder, "Video Type: %s\n", e.Attributes.VideoType)
   }

   return strings.TrimSpace(builder.String())
}
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
