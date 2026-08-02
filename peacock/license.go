package peacock

import (
   "bytes"
   "fmt"
   "io"
   "net/http"
)

// AcquireLicense sends a Widevine license acquisition request to the licence acquisition URL
// and returns the raw license bytes. The licenceAcquisitionUrl is the full URL returned in
// the PlayoutVod response.
func (c *Client) AcquireLicense(licenceAcquisitionUrl string, challenge []byte) ([]byte, error) {
   if licenceAcquisitionUrl == "" {
      return nil, fmt.Errorf("acquire license: empty licenceAcquisitionUrl")
   }
   if len(challenge) == 0 {
      return nil, fmt.Errorf("acquire license: empty challenge")
   }
   req, err := http.NewRequest(http.MethodPost, licenceAcquisitionUrl, bytes.NewReader(challenge))
   if err != nil {
      return nil, fmt.Errorf("acquire license: create request: %w", err)
   }
   resp, err := doRequest(http.DefaultClient, req)
   if err != nil {
      return nil, fmt.Errorf("acquire license: %w", err)
   }
   defer resp.Body.Close()

   if resp.StatusCode != http.StatusOK {
      return nil, fmt.Errorf("acquire license: bad status %d", resp.StatusCode)
   }

   license, err := io.ReadAll(resp.Body)
   if err != nil {
      return nil, fmt.Errorf("acquire license: read body: %w", err)
   }

   return license, nil
}

// license.go
