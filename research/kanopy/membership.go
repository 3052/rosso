package kanopy

import (
   "fmt"
   "io"
   "net/http"
)

// GetMemberships fetches library memberships/domains associated with the user.
func (c *Client) GetMemberships(userID int) ([]byte, error) {
   url := fmt.Sprintf("%s/kapi/memberships?userId=%d", BaseURL, userID)

   req, err := http.NewRequest("GET", url, nil)
   if err != nil {
      return nil, err
   }

   req.Header.Set("User-Agent", c.UserAgent)
   req.Header.Set("X-Version", c.XVersion)
   req.Header.Set("Authorization", "Bearer "+c.Token)

   // Explicitly using http.DefaultClient
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()

   if resp.StatusCode != http.StatusOK {
      return nil, fmt.Errorf("get memberships failed with status: %d", resp.StatusCode)
   }

   return io.ReadAll(resp.Body)
}
