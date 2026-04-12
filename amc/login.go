package amc

import (
   "bytes"
   "encoding/json"
   "fmt"
   "net/http"
)

// Login authenticates the user. It requires the guest token (access_token) 
// retrieved from calling the Unauth() function.
func Login(guestToken, email, password string) (*AuthData, error) {
   body := map[string]string{
      "email":    email,
      "password": password,
   }
   payload, err := json.Marshal(body)
   if err != nil {
      return nil, err
   }

   req, err := http.NewRequest(http.MethodPost, "https://gw.cds.amcn.com/auth-orchestration-id/api/v1/login", bytes.NewReader(payload))
   if err != nil {
      return nil, err
   }

   req.Header.Set("authorization", "Bearer "+guestToken)
   req.Header.Set("content-type", "application/json")
   req.Header.Set("x-amcn-language", "en")
   req.Header.Set("x-amcn-network", "amcplus")
   req.Header.Set("x-amcn-platform", "web")
   req.Header.Set("x-amcn-service-group-id", "10")
   req.Header.Set("x-amcn-tenant", "amcn")
   req.Header.Set("x-amcn-device-ad-id", "-")
   req.Header.Set("x-amcn-device-id", "-")
   req.Header.Set("x-amcn-service-id", "amcplus")
   req.Header.Set("x-ccpa-do-not-sell", "doNotPassData")
   req.Header.Set("user-agent", "Go-http-client/2.0")

   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()

   if resp.StatusCode != http.StatusOK {
      return nil, fmt.Errorf("login failed with status: %d", resp.StatusCode)
   }

   // Internal envelope to strip the first layer
   var envelope struct {
      Success bool     `json:"success"`
      Status  int      `json:"status"`
      Data    AuthData `json:"data"`
   }

   if err := json.NewDecoder(resp.Body).Decode(&envelope); err != nil {
      return nil, err
   }

   return &envelope.Data, nil
}
