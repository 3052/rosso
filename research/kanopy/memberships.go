// memberships.go
package kanopy

import (
   "encoding/json"
   "fmt"
   "io"
   "net/url"

   "41.neocities.org/maya"
)

type Membership struct {
   DomainId int `json:"domainId"`
   UserId   int `json:"userId"`
}

type GetMembershipsResponse struct {
   List []Membership `json:"list"`
}

func GetMemberships(userId int, authorization string) (*GetMembershipsResponse, error) {
   targetUrl, err := url.Parse("https://www.kanopy.com/kapi/memberships")
   if err != nil {
      return nil, err
   }

   queryParams := targetUrl.Query()
   queryParams.Set("userId", fmt.Sprintf("%d", userId))
   targetUrl.RawQuery = queryParams.Encode()

   headers := map[string]string{
      "authorization": authorization,
   }

   resp, err := maya.Get(targetUrl, headers)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()

   respBody, err := io.ReadAll(resp.Body)
   if err != nil {
      return nil, err
   }

   var membershipsResp GetMembershipsResponse
   if err := json.Unmarshal(respBody, &membershipsResp); err != nil {
      return nil, err
   }

   return &membershipsResp, nil
}
