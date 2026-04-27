package rakuten

import (
   "io"
   "net/url"

   "41.neocities.org/maya"
)

type LicenseUrl string

func AcquireLicense(licenseUrl LicenseUrl, challengeData []byte) ([]byte, error) {
   endpoint, err := url.Parse(string(licenseUrl))
   if err != nil {
      return nil, err
   }

   resp, err := maya.Post(endpoint, nil, challengeData)
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
