// season_episodes.go
package hboMax

import (
   "encoding/json"
   "fmt"
   "io"
   "net/http"
   "sort"
)

// Episode represents a parsed episode from a season.
type Episode struct {
   VideoID       string
   EditID        string
   Name          string
   Description   string
   SeasonNumber  int
   EpisodeNumber int
   AirDate       string
}

// Internal structures for parsing the Season API response
type seasonResponse struct {
   Included []seasonIncludedData `json:"included"`
}

type seasonIncludedData struct {
   ID            string              `json:"id"`
   Type          string              `json:"type"`
   Attributes    seasonAttributes    `json:"attributes"`
   Relationships seasonRelationships `json:"relationships"`
}

type seasonAttributes struct {
   MaterialType  string `json:"materialType"`
   Name          string `json:"name"`
   Description   string `json:"description"`
   SeasonNumber  int    `json:"seasonNumber"`
   EpisodeNumber int    `json:"episodeNumber"`
   AirDate       string `json:"airDate"`
}

type seasonRelationships struct {
   Edit seasonRelationshipEdit `json:"edit"`
}

type seasonRelationshipEdit struct {
   Data seasonRelationshipData `json:"data"`
}

type seasonRelationshipData struct {
   ID   string `json:"id"`
   Type string `json:"type"`
}

// GetSeasonEpisodes fetches all episodes for a given show ID and season number.
func (c *Client) GetSeasonEpisodes(showID string, seasonNumber int) ([]Episode, error) {
   // The collection ID '227084608563650952176059252419027445293' represents the "Season Tabbed Content" UI rail.
   endpoint := fmt.Sprintf("/cms/collections/227084608563650952176059252419027445293?include=default&decorators=viewingHistory,isFavorite,contentAction,badges&pf[show.id]=%s&pf[seasonNumber]=%d", showID, seasonNumber)

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

   body, err := io.ReadAll(resp.Body)
   if err != nil {
      return nil, err
   }

   var apiResp seasonResponse
   if err := json.Unmarshal(body, &apiResp); err != nil {
      return nil, err
   }

   var episodes []Episode
   for _, item := range apiResp.Included {
      if item.Type == "video" && item.Attributes.MaterialType == "EPISODE" {
         episodes = append(episodes, Episode{
            VideoID:       item.ID,
            EditID:        item.Relationships.Edit.Data.ID,
            Name:          item.Attributes.Name,
            Description:   item.Attributes.Description,
            SeasonNumber:  item.Attributes.SeasonNumber,
            EpisodeNumber: item.Attributes.EpisodeNumber,
            AirDate:       item.Attributes.AirDate,
         })
      }
   }

   // Sort episodes by EpisodeNumber just in case the API returned them out of order
   sort.Slice(episodes, func(i, j int) bool {
      return episodes[i].EpisodeNumber < episodes[j].EpisodeNumber
   })

   return episodes, nil
}
