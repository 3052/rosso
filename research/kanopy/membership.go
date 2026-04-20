// file: memberships.go
package kanopy

import (
   "encoding/json"
   "io"
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

func (l *LoginResponse) GetMemberships() (*MembershipsResponse, error) {
   targetUrl := &url.URL{
      Scheme:   "https",
      Host:     "www.kanopy.com",
      Path:     "/kapi/memberships",
      RawQuery: "userId=" + strconv.Itoa(l.UserID),
   }

   headers := map[string]string{
      "authorization": "Bearer " + l.JWT,
      "user-agent":    "!",
      "x-version":     "!/!/!/!",
   }

   resp, err := maya.Get(targetUrl, headers)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()

   bodyBytes, err := io.ReadAll(resp.Body)
   if err != nil {
      return nil, err
   }

   var membershipsResp MembershipsResponse
   if err := json.Unmarshal(bodyBytes, &membershipsResp); err != nil {
      return nil, err
   }

   return &membershipsResp, nil
}
