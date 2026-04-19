package kanopy

import (
   "encoding/json"
   "fmt"
   "io"
   "net/url"

   "41.neocities.org/maya"
)

type Membership struct {
   DomainId int    `json:"domainId"`
   UserId   int    `json:"userId"`
   Status   string `json:"status"`
}

type MembershipsResponse struct {
   List []Membership `json:"list"`
}

func GetMemberships(userId int, jwt string) (*MembershipsResponse, error) {
   targetUrl, err := url.Parse(fmt.Sprintf("https://www.kanopy.com/kapi/memberships?userId=%d", userId))
   if err != nil {
      return nil, err
   }

   requestHeaders := map[string]string{
      "authorization": "Bearer " + jwt,
   }

   response, err := maya.Get(targetUrl, requestHeaders)
   if err != nil {
      return nil, err
   }
   defer response.Body.Close()

   responseBytes, err := io.ReadAll(response.Body)
   if err != nil {
      return nil, err
   }

   var membershipsResponse MembershipsResponse
   err = json.Unmarshal(responseBytes, &membershipsResponse)
   if err != nil {
      return nil, err
   }

   return &membershipsResponse, nil
}
