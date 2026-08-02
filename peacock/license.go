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

   client, err := mtlsClient(c.HTTP.Timeout)
   if err != nil {
      return nil, fmt.Errorf("acquire license: %w", err)
   }

   req, err := http.NewRequest(http.MethodPost, licenceAcquisitionUrl, bytes.NewReader(challenge))
   if err != nil {
      return nil, fmt.Errorf("acquire license: create request: %w", err)
   }

   req.Header.Set("Content-Type", "application/octet-stream")
   req.Header.Set("User-Agent", "Media3Player/7.6.100 (Linux;Android 12) AndroidXMedia3-Sky-CVSDK/1.8.0 [emulator64_x86_64_arm64, sdk_gphone64_x86_64, Google, 31]")

   resp, err := doRequest(client, req)
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
