// File: get_memberships.go
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

func GetMemberships(userID int, token string) (*MembershipsResponse, error) {
   reqURL, err := url.Parse("https://www.kanopy.com/kapi/memberships")
   if err != nil {
      return nil, err
   }

   query := reqURL.Query()
   query.Set("userId", strconv.Itoa(userID))
   reqURL.RawQuery = query.Encode()

   headers := map[string]string{
      "authorization": "Bearer " + token,
      "user-agent":    "!",
      "x-version":     "!/!/!/!",
   }

   resp, err := maya.Get(reqURL, headers)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()

   var result MembershipsResponse
   if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
      return nil, err
   }

   return &result, nil
}
