package kanopy

import (
   "fmt"
   "io"
   "net/url"

   "41.neocities.org/maya"
)

func GetWidevineLicense(drmLicenseID string, jwt string, licensePayload []byte) ([]byte, error) {
   targetUrl, err := url.Parse(fmt.Sprintf("https://www.kanopy.com/kapi/licenses/widevine/%s", drmLicenseID))
   if err != nil {
      return nil, err
   }

   requestHeaders := map[string]string{
      "authorization": "Bearer " + jwt,
   }

   response, err := maya.Post(targetUrl, requestHeaders, licensePayload)
   if err != nil {
      return nil, err
   }
   defer response.Body.Close()

   responseBytes, err := io.ReadAll(response.Body)
   if err != nil {
      return nil, err
   }

   return responseBytes, nil
}
