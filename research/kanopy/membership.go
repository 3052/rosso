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
   SiteName           string `json:"sitename"`
   Subdomain          string `json:"subdomain"`
   TicketsAvailable   int    `json:"ticketsAvailable"`
   MaxTicketsPerMonth int    `json:"maxTicketsPerMonth"`
}

type MembershipsResponse struct {
   List []Membership `json:"list"`
}

func GetMemberships(loginResp *LoginResponse) (*MembershipsResponse, error) {
   membershipsURL, err := url.Parse("https://www.kanopy.com/kapi/memberships?userId=" + strconv.Itoa(loginResp.UserID))
   if err != nil {
      return nil, err
   }

   headers := map[string]string{
      "authorization": "Bearer " + loginResp.JWT,
      "user-agent":    "!",
      "x-version":     "!/!/!/!",
   }

   resp, err := maya.Get(membershipsURL, headers)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()

   var memResp MembershipsResponse
   if err := json.NewDecoder(resp.Body).Decode(&memResp); err != nil {
      return nil, err
   }

   return &memResp, nil
}
