package amc

import (
   "encoding/json"
   "fmt"
   "net/http"
)

func Unauth() (*AuthResponse, error) {
   req, err := http.NewRequest(http.MethodPost, "https://gw.cds.amcn.com/auth-orchestration-id/api/v1/unauth", nil)
   if err != nil {
      return nil, err
   }

   req.Header.Set("x-amcn-network", "amcplus")
   req.Header.Set("x-amcn-platform", "web")
   req.Header.Set("x-amcn-tenant", "amcn")
   req.Header.Set("x-amcn-device-id", "-")
   req.Header.Set("x-amcn-language", "en")
   req.Header.Set("user-agent", "Go-http-client/2.0")

   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()

   if resp.StatusCode != http.StatusOK {
      return nil, fmt.Errorf("unauth failed with status: %d", resp.StatusCode)
   }

   var authResp AuthResponse
   if err := json.NewDecoder(resp.Body).Decode(&authResp); err != nil {
      return nil, err
   }

   return &authResp, nil
}
