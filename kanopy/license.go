package kanopy

import (
   "io"
   "net/url"

   "41.neocities.org/maya"
)

func GetLicense(jwt, drmLicenseId string, payload []byte) ([]byte, error) {
   licenseUrl, err := url.Parse("https://www.kanopy.com/kapi/licenses/widevine/" + drmLicenseId)
   if err != nil {
      return nil, err
   }

   headers := map[string]string{
      "authorization": "Bearer " + jwt,
      "user-agent":    "!",
      "x-version":     "!/!/!/!",
   }

   resp, err := maya.Post(licenseUrl, headers, payload)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()

   return io.ReadAll(resp.Body)
}
