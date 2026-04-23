package kanopy

import (
   "io"
   "net/url"

   "41.neocities.org/maya"
)

func CreateLicense(login *LoginResponse, manifestData *Manifest, challenge []byte) ([]byte, error) {
   endpoint := &url.URL{
      Scheme: "https",
      Host:   "www.kanopy.com",
      Path:   "/kapi/licenses/widevine/" + manifestData.DrmLicenseId,
   }

   headers := map[string]string{
      "authorization": "Bearer " + login.Jwt,
   }

   resp, err := maya.Post(endpoint, headers, challenge)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()

   return io.ReadAll(resp.Body)
}
