package kanopy

import (
   "io"
   "net/url"

   "41.neocities.org/maya"
)

func PostLicense(loginData *Login, manifestData *Manifest, challenge []byte) ([]byte, error) {
   endpoint := &url.URL{
      Scheme: "https",
      Host:   "www.kanopy.com",
      Path:   "/kapi/licenses/widevine/" + manifestData.DrmLicenseID,
   }

   headers := map[string]string{
      "authorization": "Bearer " + loginData.Jwt,
      "x-version":     "!/!/!/!",
   }

   resp, err := maya.Post(endpoint, headers, challenge)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()

   return io.ReadAll(resp.Body)
}
