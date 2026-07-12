// 6_get_license.go
package unext

import (
   "bytes"
   "fmt"
   "io"
   "net/http"
   "net/url"
)

// GetLicense sends the Widevine challenge and retrieves the license binary data
func GetLicense(client *http.Client, licenseUrl, playToken string, challenge []byte) ([]byte, error) {
   u, err := url.Parse(licenseUrl)
   if err != nil {
      return nil, err
   }

   // Append play_token as a query parameter
   q := u.Query()
   q.Set("play_token", playToken)
   u.RawQuery = q.Encode()

   req, err := http.NewRequest("POST", u.String(), bytes.NewReader(challenge))
   if err != nil {
      return nil, err
   }

   req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:140.0) Gecko/20100101 Firefox/140.0")
   req.Header.Set("Accept", "*/*")
   req.Header.Set("Accept-Language", "en-US,en;q=0.5")
   req.Header.Set("Origin", "https://video.unext.jp")
   req.Header.Set("Referer", "https://video.unext.jp/")
   // Note: Content-Type is intentionally omitted to match the HAR, the server accepts the binary body as-is

   resp, err := client.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()

   if resp.StatusCode != http.StatusOK {
      return nil, fmt.Errorf("license request failed with status code: %d", resp.StatusCode)
   }

   licenseData, err := io.ReadAll(resp.Body)
   if err != nil {
      return nil, err
   }

   return licenseData, nil
}
