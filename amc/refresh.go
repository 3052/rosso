package amc

import (
   "bytes"
   "encoding/json"
   "fmt"
   "io"
   "net/http"
)

func SeasonEpisodes(authToken string, seasonID int) (*ContentNode, error) {
   url := fmt.Sprintf("https://gw.cds.amcn.com/content-compiler-cr/api/v1/content/amcn/amcplus/type/season-episodes/id/%d", seasonID)
   
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
      return nil, fmt.Errorf("season episodes failed with status: %d", resp.StatusCode)
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

func License(licenseURL, bcovAuth string, challengePayload []byte) ([]byte, error) {
   req, err := http.NewRequest(http.MethodPost, licenseURL, bytes.NewReader(challengePayload))
   if err != nil {
      return nil, err
   }

   req.Header.Set("bcov-auth", bcovAuth)
   req.Header.Set("user-agent", "Go-http-client/2.0")

   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()

   if resp.StatusCode != http.StatusOK {
      return nil, fmt.Errorf("license request failed with status: %d", resp.StatusCode)
   }

   return io.ReadAll(resp.Body)
}

func Refresh(refreshToken string) (*AuthData, error) {
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
