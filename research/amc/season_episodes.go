package amc

import (
   "encoding/json"
   "fmt"
   "net/http"
)

// EpisodesMetadata recursively traverses the Server-Driven UI tree 
// and extracts only the Metadata for playable episodes.
func (c *ContentNode) EpisodesMetadata() []*Metadata {
   var metadata []*Metadata
   var walk func(node *ContentNode)
   walk = func(node *ContentNode) {
      if node.Type == "card" && node.Properties != nil && node.Properties.ContentType == "episode" && node.Properties.Metadata != nil {
         metadata = append(metadata, node.Properties.Metadata)
      }
      for i := range node.Children {
         walk(&node.Children[i])
      }
   }
   walk(c)
   return metadata
}

func SeasonEpisodes(authToken, seasonID string) (*ContentNode, error) {
   url := fmt.Sprintf("https://gw.cds.amcn.com/content-compiler-cr/api/v1/content/amcn/amcplus/type/season-episodes/id/%s", seasonID)

   req, err := http.NewRequest(http.MethodGet, url, nil)
   if err != nil {
      return nil, err
   }

   req.Header.Set("authorization", "Bearer "+authToken)
   req.Header.Set("x-amcn-network", "amcplus")
   req.Header.Set("x-amcn-platform", "android")
   req.Header.Set("x-amcn-tenant", "amcn")
   req.Header.Set("user-agent", "Go-http-client/2.0")

   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()

   if resp.StatusCode != http.StatusOK {
      return nil, fmt.Errorf("season episodes failed with status: %d", resp.StatusCode)
   }

   // Internal envelope to strip the first layer
   var envelope struct {
      Success bool        `json:"success"`
      Status  int         `json:"status"`
      Data    ContentNode `json:"data"`
   }

   if err := json.NewDecoder(resp.Body).Decode(&envelope); err != nil {
      return nil, err
   }

   return &envelope.Data, nil
}
