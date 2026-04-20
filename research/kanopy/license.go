package kanopy

import (
   "io"
   "net/url"

   "41.neocities.org/maya"
)

func GetLicense(jwt, drmLicenseID string, payload []byte) ([]byte, error) {
   licenseURL, err := url.Parse("https://www.kanopy.com/kapi/licenses/widevine/" + drmLicenseID)
   if err != nil {
      return nil, err
   }

   headers := map[string]string{
      "authorization": "Bearer " + jwt,
      "user-agent":    "!",
      "x-version":     "!/!/!/!",
   }

   resp, err := maya.Post(licenseURL, headers, payload)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()

   return io.ReadAll(resp.Body)
}
