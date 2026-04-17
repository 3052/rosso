package plex

import (
   "bytes"
   "fmt"
   "io"
   "net/http"
   "net/url"
)

// GetWidevineLicense submits the CDM challenge/payload to Plex's license server
func GetWidevineLicense(partID, plexToken string, payload []byte) ([]byte, error) {
   baseURL, _ := url.Parse(fmt.Sprintf("https://vod.provider.plex.tv/library/parts/%s/license", partID))

   q := baseURL.Query()
   q.Set("x-plex-drm", "widevine")
   q.Set("x-plex-token", plexToken)
   baseURL.RawQuery = q.Encode()

   // Use bytes.NewReader to automatically allow http.NewRequest to set Content-Length
   req, err := http.NewRequest("POST", baseURL.String(), bytes.NewReader(payload))
   if err != nil {
      return nil, err
   }

   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()

   if resp.StatusCode < 200 || resp.StatusCode >= 300 {
      return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
   }

   return io.ReadAll(resp.Body)
}
