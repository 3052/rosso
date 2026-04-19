// license.go
package kanopy

import (
   "fmt"
   "io"
   "net/url"

   "41.neocities.org/maya"
)

func GetWidevineLicense(drmLicenseID string, payload []byte, authorization string) ([]byte, error) {
   targetUrl, err := url.Parse(fmt.Sprintf("https://www.kanopy.com/kapi/licenses/widevine/%s", url.PathEscape(drmLicenseID)))
   if err != nil {
      return nil, err
   }

   headers := map[string]string{
      "authorization": authorization,
   }

   resp, err := maya.Post(targetUrl, headers, payload)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()

   respBody, err := io.ReadAll(resp.Body)
   if err != nil {
      return nil, err
   }

   return respBody, nil
}
