// 3_migrate_tokens.go
package unext

import (
   "bytes"
   "fmt"
   "net/http"
)

// MigrateTokens retrieves the _at JWT access token cookie
func MigrateTokens(client *http.Client) error {
   urlStr := "https://myaccount.unext.jp/api/migrateTokens"

   payload := []byte(`{"forceMigration":true}`)

   req, err := http.NewRequest("POST", urlStr, bytes.NewBuffer(payload))
   if err != nil {
      return err
   }

   req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:140.0) Gecko/20100101 Firefox/140.0")
   req.Header.Set("Accept", "*/*")
   req.Header.Set("Accept-Language", "en-US,en;q=0.5")
   req.Header.Set("Content-Type", "application/json")
   req.Header.Set("Origin", "https://myaccount.unext.jp")
   req.Header.Set("Referer", "https://myaccount.unext.jp/oauth-migration?backurl=https%3A%2F%2Fvideo.unext.jp%2Ftitle%2FSID0020149")

   resp, err := client.Do(req)
   if err != nil {
      return err
   }
   defer resp.Body.Close()

   if resp.StatusCode != http.StatusOK {
      return fmt.Errorf("migrateTokens failed with status code: %d", resp.StatusCode)
   }

   // The _at and _rt cookies are automatically stored in the client's CookieJar
   return nil
}
