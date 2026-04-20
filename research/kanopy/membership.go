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
   IdentityID         int64  `json:"identityId"`
   DomainID           int64  `json:"domainId"`
   UserID             int64  `json:"userId"`
   Status             string `json:"status"`
   IsDefault          bool   `json:"isDefault"`
   Sitename           string `json:"sitename"`
   Subdomain          string `json:"subdomain"`
   TicketsAvailable   int    `json:"ticketsAvailable"`
   MaxTicketsPerMonth int    `json:"maxTicketsPerMonth"`
}

func GetMemberships(loginResponse *LoginResponse) (*MembershipsResponse, error) {
   target := &url.URL{
      Scheme:   "https",
      Host:     "www.kanopy.com",
      Path:     "/kapi/memberships",
      RawQuery: "userId=" + strconv.FormatInt(loginResponse.UserID, 10),
   }

   headers := map[string]string{
      "authorization": "Bearer " + loginResponse.JWT,
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
