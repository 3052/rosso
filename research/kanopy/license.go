// File: get_widevine_license.go
package kanopy

import (
   "io"
   "net/url"

   "41.neocities.org/maya"
)

func GetWidevineLicense(drmLicenseID string, payload []byte, token string) ([]byte, error) {
   reqURL, err := url.Parse("https://www.kanopy.com/kapi/licenses/widevine/" + drmLicenseID)
   if err != nil {
      return nil, err
   }

   headers := map[string]string{
      "authorization": "Bearer " + token,
      "user-agent":    "!",
      "x-version":     "!/!/!/!",
   }

   resp, err := maya.Post(reqURL, headers, payload)
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
