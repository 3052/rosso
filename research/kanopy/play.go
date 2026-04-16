package kanopy

import (
   "bytes"
   "encoding/json"
   "fmt"
   "io"
   "net/http"
)

type PlayRequest struct {
   DomainID int `json:"domainId"`
   UserID   int `json:"userId"`
   VideoID  int `json:"videoId"`
}

// CreatePlay registers a play event to retrieve stream manifests and DRM license IDs.
func (c *Client) CreatePlay(domainID, userID, videoID int) ([]byte, error) {
   payload := PlayRequest{
      DomainID: domainID,
      UserID:   userID,
      VideoID:  videoID,
   }

   body, err := json.Marshal(payload)
   if err != nil {
      return nil, err
   }

   req, err := http.NewRequest("POST", BaseURL+"/kapi/plays", bytes.NewBuffer(body))
   if err != nil {
      return nil, err
   }

   req.Header.Set("X-Version", c.XVersion)
   req.Header.Set("Authorization", "Bearer "+c.Token)
   req.Header.Set("Content-Type", "application/json")
   req.Header.Set("User-Agent", c.UserAgent)

   // Explicitly using http.DefaultClient
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()

   if resp.StatusCode != http.StatusOK {
      return nil, fmt.Errorf("create play failed with status: %d", resp.StatusCode)
   }

   return io.ReadAll(resp.Body)
}
