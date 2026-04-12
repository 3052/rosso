package amc

import (
   "encoding/json"
   "fmt"
   "net/http"
)

func Unauth() (*AuthData, error) {
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
func SeriesDetail(authToken, seriesID string) (*ContentNode, error) {
   url := fmt.Sprintf("https://gw.cds.amcn.com/content-compiler-cr/api/v1/content/amcn/amcplus/type/series-detail/id/%s", seriesID)

   req, err := http.NewRequest(http.MethodGet, url, nil)
   if err != nil {
      return nil, err
   }

   req.Header.Set("authorization", "Bearer "+authToken)
   req.Header.Set("x-amcn-network", "amcplus")
   req.Header.Set("x-amcn-platform", "android")
   req.Header.Set("x-amcn-tenant", "amcn")
   req.Header.Set("user-agent", "Go-http-client/2.0")

   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()

   if resp.StatusCode != http.StatusOK {
      return nil, fmt.Errorf("series detail failed with status: %d", resp.StatusCode)
   }

   // Internal envelope to strip the first layer
   var envelope struct {
      Success bool        `json:"success"`
      Status  int         `json:"status"`
      Data    ContentNode `json:"data"`
   }

   if err := json.NewDecoder(resp.Body).Decode(&envelope); err != nil {
      return nil, err
   }

   return &envelope.Data, nil
}
