// membership.go
package kanopy

import (
   "encoding/json"
   "net/url"
   "strconv"

   "41.neocities.org/maya"
)

type Membership struct {
   IdentityId int    `json:"identityId"`
   DomainId   int    `json:"domainId"`
   UserId     int    `json:"userId"`
   Status     string `json:"status"`
   IsDefault  bool   `json:"isDefault"`
}

type MembershipResponse struct {
   List []Membership `json:"list"`
}

func GetMemberships(session *Session) (*MembershipResponse, error) {
   targetUrl, parseError := url.Parse("https://www.kanopy.com/kapi/memberships")
   if parseError != nil {
      return nil, parseError
   }

   query := targetUrl.Query()
   query.Set("userId", strconv.Itoa(session.UserId))
   targetUrl.RawQuery = query.Encode()

   headers := map[string]string{
      "authorization": "Bearer " + session.Authorization,
      "x-version":     "!/!/!/!",
      "user-agent":    "!",
   }

   resp, requestError := maya.Get(targetUrl, headers)
   if requestError != nil {
      return nil, requestError
   }
   defer resp.Body.Close()

   var memResp MembershipResponse
   decodeError := json.NewDecoder(resp.Body).Decode(&memResp)
   if decodeError != nil {
      return nil, decodeError
   }
   return &memResp, nil
}
