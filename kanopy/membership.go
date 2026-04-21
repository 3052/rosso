package kanopy

import (
   "encoding/json"
   "net/url"
   "strconv"

   "41.neocities.org/maya"
)

type Membership struct {
   IdentityId         int    `json:"identityId"`
   DomainId           int    `json:"domainId"`
   UserId             int    `json:"userId"`
   Status             string `json:"status"`
   IsDefault          bool   `json:"isDefault"`
   SiteName           string `json:"sitename"`
   Subdomain          string `json:"subdomain"`
   TicketsAvailable   int    `json:"ticketsAvailable"`
   MaxTicketsPerMonth int    `json:"maxTicketsPerMonth"`
}

func GetMemberships(loginResp *LoginResponse) ([]Membership, error) {
   membershipsUrl, err := url.Parse("https://www.kanopy.com/kapi/memberships?userId=" + strconv.Itoa(loginResp.UserId))
   if err != nil {
      return nil, err
   }

   headers := map[string]string{
      "authorization": "Bearer " + loginResp.Jwt,
      "user-agent":    "!",
      "x-version":     "!/!/!/!",
   }

   resp, err := maya.Get(membershipsUrl, headers)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   var result struct {
      List []Membership `json:"list"`
   }
   if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
      return nil, err
   }
   return result.List, nil
}
