package amc

import (
   "encoding/json"
   "fmt"
   "net/http"
)

func Refresh(refreshToken string) (*AuthResponse, error) {
   req, err := http.NewRequest(http.MethodPost, "https://gw.cds.amcn.com/auth-orchestration-id/api/v1/refresh", nil)
   if err != nil {
      return nil, err
   }

   req.Header.Set("authorization", "Bearer "+refreshToken)
   req.Header.Set("user-agent", "Go-http-client/2.0")

   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()

   if resp.StatusCode != http.StatusOK {
      return nil, fmt.Errorf("refresh failed with status: %d", resp.StatusCode)
   }

   var authResp AuthResponse
   if err := json.NewDecoder(resp.Body).Decode(&authResp); err != nil {
      return nil, err
   }

   return &authResp, nil
}
