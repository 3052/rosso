// licenses.go
package kanopy

import (
   "io"
   "net/url"

   "41.neocities.org/maya"
)

func GetLicenseWidevine(jwt string, drmLicenseID string, payload []byte) ([]byte, error) {
   targetUrl, err := url.Parse("https://www.kanopy.com/kapi/licenses/widevine/" + drmLicenseID)
   if err != nil {
      return nil, err
   }

   headers := map[string]string{
      "authorization": "Bearer " + jwt,
      "x-version":     "!/!/!/!",
   }

   resp, err := maya.Post(targetUrl, headers, payload)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()

   licenseBody, err := io.ReadAll(resp.Body)
   if err != nil {
      return nil, err
   }

   return licenseBody, nil
}
