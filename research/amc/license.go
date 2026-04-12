package amc

import (
   "bytes"
   "fmt"
   "io"
   "net/http"
)

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
