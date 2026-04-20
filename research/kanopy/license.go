// license.go
package kanopy

import (
   "io"
   "net/url"

   "41.neocities.org/maya"
)

func GetWidevineLicense(Authorization string, DrmLicenseID string, payload []byte) ([]byte, error) {
   targetUrl, parseError := url.Parse("https://www.kanopy.com/kapi/licenses/widevine/" + DrmLicenseID)
   if parseError != nil {
      return nil, parseError
   }

   headers := map[string]string{
      "authorization": "Bearer " + Authorization,
      "x-version":     "!/!/!/!",
      "user-agent":    "!",
   }

   resp, requestError := maya.Post(targetUrl, headers, payload)
   if requestError != nil {
      return nil, requestError
   }
   defer resp.Body.Close()

   return io.ReadAll(resp.Body)
}
