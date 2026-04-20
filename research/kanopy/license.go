// file: widevine.go
package kanopy

import (
   "io"
   "net/url"

   "41.neocities.org/maya"
)

func GetWidevineLicense(drmLicenseID string, jwt string, body []byte) ([]byte, error) {
   targetUrl := &url.URL{
      Scheme: "https",
      Host:   "www.kanopy.com",
      Path:   "/kapi/licenses/widevine/" + drmLicenseID,
   }

   headers := map[string]string{
      "authorization": "Bearer " + jwt,
      "user-agent":    "!",
      "x-version":     "!/!/!/!",
   }

   resp, err := maya.Post(targetUrl, headers, body)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()

   bodyBytes, err := io.ReadAll(resp.Body)
   if err != nil {
      return nil, err
   }

   return bodyBytes, nil
}
