package amc

import (
   "encoding/json"
   "fmt"
   "net/http"
)

func SeriesDetail(authToken, seriesID string) (*ContentResponse, error) {
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

   var contentResp ContentResponse
   if err := json.NewDecoder(resp.Body).Decode(&contentResp); err != nil {
      return nil, err
   }

   return &contentResp, nil
}
