// memberships.go
package kanopy

import (
   "encoding/json"
   "fmt"
   "io"
   "net/url"

   "41.neocities.org/maya"
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

func GetMemberships(UserID int, authorization string) (*MembershipsResponse, error) {
   targetURL, err := url.Parse(fmt.Sprintf("https://www.kanopy.com/kapi/memberships?userId=%d", UserID))
   if err != nil {
      return nil, err
   }

   headers := map[string]string{
      "authorization": authorization,
      "user-agent":    "!",
      "x-version":     "!/!/!/!",
   }

   resp, err := maya.Get(targetURL, headers)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()

   if resp.StatusCode != 200 {
      return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
   }

   respBody, err := io.ReadAll(resp.Body)
   if err != nil {
      return nil, err
   }

   var membershipsResp MembershipsResponse
   if err := json.Unmarshal(respBody, &membershipsResp); err != nil {
      return nil, err
   }

   return &membershipsResp, nil
}
