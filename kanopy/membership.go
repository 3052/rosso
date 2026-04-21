// memberships.go
package kanopy

import (
   "encoding/json"
   "io"
   "net/url"
   "strconv"

   "41.neocities.org/maya"
)

type MembershipsResponse struct {
   List []Membership `json:"list"`
}

type Membership struct {
   IdentityId         int    `json:"identityId"`
   DomainId           int    `json:"domainId"`
   UserId             int    `json:"userId"`
   Status             string `json:"status"`
   IsDefault          bool   `json:"isDefault"`
   Sitename           string `json:"sitename"`
   Subdomain          string `json:"subdomain"`
   TicketsAvailable   int    `json:"ticketsAvailable"`
   MaxTicketsPerMonth int    `json:"maxTicketsPerMonth"`
}

func GetMemberships(userId int, jwt string) (*MembershipsResponse, error) {
   query := url.Values{}
   query.Set("userId", strconv.Itoa(userId))

   target := &url.URL{
      Scheme:   "https",
      Host:     "www.kanopy.com",
      Path:     "/kapi/memberships",
      RawQuery: query.Encode(),
   }

   headers := map[string]string{
      "authorization": "Bearer " + jwt,
      "user-agent":    "!",
      "x-version":     "!/!/!/!",
   }

   resp, err := maya.Get(target, headers)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()

   respBytes, err := io.ReadAll(resp.Body)
   if err != nil {
      return nil, err
   }

   var membershipsResp MembershipsResponse
   if err := json.Unmarshal(respBytes, &membershipsResp); err != nil {
      return nil, err
   }

   return &membershipsResp, nil
}
