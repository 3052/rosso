// licenses.go
package kanopy

import (
   "fmt"
   "io"
   "net/url"

   "41.neocities.org/maya"
)

func GetWidevineLicense(licenseID string, payload []byte, authorization string) ([]byte, error) {
   targetURL, err := url.Parse(fmt.Sprintf("https://www.kanopy.com/kapi/licenses/widevine/%s", licenseID))
   if err != nil {
      return nil, err
   }

   headers := map[string]string{
      "authorization": authorization,
      "user-agent":    "!",
      "x-version":     "!/!/!/!",
   }

   resp, err := maya.Post(targetURL, headers, payload)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()

   if resp.StatusCode != 200 {
      return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
   }

   respBody, err := io.ReadAll(resp.Body)
   if err != nil {
      return nil, err
   }

   return respBody, nil
}
