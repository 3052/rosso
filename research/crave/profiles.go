package crave

import (
   "encoding/json"
   "fmt"
   "io"
   "net/http"
)

type Profile struct {
   ID        string `json:"id"`
   AccountID string `json:"accountId"`
   Nickname  string `json:"nickname"`
   HasPin    bool   `json:"hasPin"`
   Master    bool   `json:"master"`
   Maturity  string `json:"maturity"`
}

// GetProfiles fetches the list of profiles associated with the account.
func (c *Client) GetProfiles(accountID, accessToken string) ([]*Profile, error) {
   endpoint := fmt.Sprintf("%s/api/profile/v2/account/%s", BaseURL, accountID)

   req, err := http.NewRequest("GET", endpoint, nil)
   if err != nil {
      return nil, err
   }

   req.Header.Set("Authorization", "Bearer "+accessToken)
   req.Header.Set("User-Agent", UserAgent)
   req.Header.Set("Accept", "application/json, text/plain, */*")
   req.Header.Set("Origin", "https://www.crave.ca")

   resp, err := c.HTTPClient.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()

   if resp.StatusCode < 200 || resp.StatusCode >= 300 {
      body, _ := io.ReadAll(resp.Body)
      return nil, fmt.Errorf("failed to fetch profiles with status %d: %s", resp.StatusCode, string(body))
   }

   var profiles []*Profile
   if err := json.NewDecoder(resp.Body).Decode(&profiles); err != nil {
      return nil, err
   }

   return profiles, nil
}
