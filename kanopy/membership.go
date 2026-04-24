package kanopy

import (
   "encoding/json"
   "io"
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
   Sitename           string `json:"sitename"`
   Subdomain          string `json:"subdomain"`
   TicketsAvailable   int    `json:"ticketsAvailable"`
   MaxTicketsPerMonth int    `json:"maxTicketsPerMonth"`
}

func GetMemberships(login *LoginResponse) ([]Membership, error) {
   endpoint := &url.URL{
      Scheme: "https",
      Host:   "www.kanopy.com",
      Path:   "/kapi/memberships",
   }

   query := url.Values{}
   query.Set("userId", strconv.Itoa(login.UserId))
   endpoint.RawQuery = query.Encode()

   headers := map[string]string{
      "authorization": "Bearer " + login.Jwt,
   }

   resp, err := maya.Get(endpoint, headers)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()

   respBody, err := io.ReadAll(resp.Body)
   if err != nil {
      return nil, err
   }
   var result struct {
      List []Membership `json:"list"`
   }
   if err := json.Unmarshal(respBody, &result); err != nil {
      return nil, err
   }
   return result.List, nil
}
