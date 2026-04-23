package kanopy

import (
   "encoding/json"
   "net/url"
   "strconv"

   "41.neocities.org/maya"
)

type Membership struct {
   DomainId   int `json:"domainId"`
   IdentityId int `json:"identityId"`
}

type Memberships struct {
   List []Membership `json:"list"`
}

func GetMemberships(loginData *Login) (*Memberships, error) {
   query := url.Values{}
   query.Set("userId", strconv.Itoa(loginData.UserId))

   endpoint := &url.URL{
      Scheme:   "https",
      Host:     "www.kanopy.com",
      Path:     "/kapi/memberships",
      RawQuery: query.Encode(),
   }

   headers := map[string]string{
      "authorization": "Bearer " + loginData.Jwt,
      "x-version":     "!/!/!/!",
   }

   resp, err := maya.Get(endpoint, headers)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()

   var membershipsData Memberships
   if err := json.NewDecoder(resp.Body).Decode(&membershipsData); err != nil {
      return nil, err
   }

   return &membershipsData, nil
}
