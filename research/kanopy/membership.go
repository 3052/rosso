package kanopy

import (
   "encoding/json"
   "fmt"
   "io"
   "net/http"
)

type Membership struct {
   IdentityID         int    `json:"identityId"`
   DomainID           int    `json:"domainId"`
   UserID             int    `json:"userId"`
   Status             string `json:"status"`
   IsDefault          bool   `json:"isDefault"`
   Sitename           string `json:"sitename"`
   Subdomain          string `json:"subdomain"`
   TicketsAvailable   int    `json:"ticketsAvailable"`
   MaxTicketsPerMonth int    `json:"maxTicketsPerMonth"`
}

type MembershipsResponse struct {
   List []Membership `json:"list"`
}

// GetMemberships fetches the library memberships associated with the session's UserID.
func (s *Session) GetMemberships() (*MembershipsResponse, error) {
   url := fmt.Sprintf("%s/kapi/memberships?userId=%d", BaseURL, s.UserID)

   req, err := http.NewRequest("GET", url, nil)
   if err != nil {
      return nil, err
   }

   req.Header.Set("User-Agent", UserAgent)
   req.Header.Set("X-Version", XVersion)
   req.Header.Set("Authorization", "Bearer "+s.JWT)

   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()

   if resp.StatusCode != http.StatusOK {
      return nil, fmt.Errorf("get memberships failed with status: %d", resp.StatusCode)
   }

   respBody, err := io.ReadAll(resp.Body)
   if err != nil {
      return nil, err
   }

   var memResp MembershipsResponse
   if err := json.Unmarshal(respBody, &memResp); err != nil {
      return nil, err
   }

   return &memResp, nil
}
