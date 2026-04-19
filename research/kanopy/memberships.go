// memberships.go
package kanopy

import (
   "encoding/json"
   "net/url"
   "strconv"

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

func GetMemberships(jwt string, userId int) (*MembershipsResponse, error) {
   targetUrl, err := url.Parse("https://www.kanopy.com/kapi/memberships?userId=" + strconv.Itoa(userId))
   if err != nil {
      return nil, err
   }

   headers := map[string]string{
      "authorization": "Bearer " + jwt,
      "x-version":     "!/!/!/!",
   }

   resp, err := maya.Get(targetUrl, headers)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()

   var membershipsResp MembershipsResponse
   err = json.NewDecoder(resp.Body).Decode(&membershipsResp)
   if err != nil {
      return nil, err
   }

   return &membershipsResp, nil
}
