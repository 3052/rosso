package plex

import (
   "bytes"
   "fmt"
   "io"
   "net/http"
   "net/url"
)

// GetWidevineLicense submits the Widevine CDM payload using the DASH partID
// (e.g., "62d04a01700e44e5863792a9-dash") and returns the raw binary license bytes.
func GetWidevineLicense(partID, plexToken string, payload []byte) ([]byte, error) {
   baseURL, _ := url.Parse(fmt.Sprintf("https://vod.provider.plex.tv/library/parts/%s/license", partID))

   q := baseURL.Query()
   q.Set("x-plex-drm", "widevine")
   q.Set("x-plex-token", plexToken)
   baseURL.RawQuery = q.Encode()

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
