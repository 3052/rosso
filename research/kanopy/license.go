// license.go
package kanopy

import (
   "io"
   "net/url"

   "41.neocities.org/maya"
)

func GetLicense(drmLicenseId string, body []byte, jwt string) ([]byte, error) {
   target := &url.URL{
      Scheme: "https",
      Host:   "www.kanopy.com",
      Path:   "/kapi/licenses/widevine/" + drmLicenseId,
   }

   headers := map[string]string{
      "authorization": "Bearer " + jwt,
      "user-agent":    "!",
      "x-version":     "!/!/!/!",
   }

   resp, err := maya.Post(target, headers, body)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()

   respBytes, err := io.ReadAll(resp.Body)
   if err != nil {
      return nil, err
   }

   return respBytes, nil
}
