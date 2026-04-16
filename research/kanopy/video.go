package kanopy

import (
   "fmt"
   "io"
   "net/http"
)

// GetVideo fetches video metadata using the movie alias (e.g., "justwatch-14685304").
func (c *Client) GetVideo(alias string) ([]byte, error) {
   url := fmt.Sprintf("%s/kapi/videos/alias/%s", BaseURL, alias)

   req, err := http.NewRequest("GET", url, nil)
   if err != nil {
      return nil, err
   }

   req.Header.Set("X-Version", c.XVersion)
   req.Header.Set("Authorization", "Bearer "+c.Token)
   req.Header.Set("User-Agent", "Go-http-client/2.0")

   // Explicitly using http.DefaultClient
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()

   if resp.StatusCode != http.StatusOK {
      return nil, fmt.Errorf("get video failed with status: %d", resp.StatusCode)
   }

   return io.ReadAll(resp.Body)
}
