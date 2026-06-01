// link_device.go
package amazon

import (
   "fmt"
   "io"
   "net/http"
   "net/url"
   "strings"
)

func LinkDevice(client *http.Client, endpoint string, referer string, publicCode string, csrfToken string) error {
   data := url.Values{}
   data.Set("ref_", "atv_set_rd_reg")
   data.Set("publicCode", publicCode)
   data.Set("token", csrfToken)

   req, err := http.NewRequest(http.MethodPost, endpoint, strings.NewReader(data.Encode()))
   if err != nil {
      return err
   }

   req.Header.Set("Accept", "*/*")
   req.Header.Set("Accept-Language", "en-US,en;q=0.9,es-US;q=0.8,es;q=0.7")
   req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
   req.Header.Set("Referer", referer)

   resp, err := client.Do(req)
   if err != nil {
      return err
   }
   defer resp.Body.Close()

   if resp.StatusCode != http.StatusOK {
      bodyBytes, _ := io.ReadAll(resp.Body)
      return fmt.Errorf("unexpected response with the codeBasedLinking request: %s [%d]", string(bodyBytes), resp.StatusCode)
   }

   return nil
}
