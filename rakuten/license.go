package rakuten

import (
   "io"
   "net/url"

   "41.neocities.org/maya"
)

func (info *StreamInfo) FetchLicense(challenge []byte) ([]byte, error) {
   target, err := url.Parse(info.LicenseUrl)
   if err != nil {
      return nil, err
   }

   resp, err := maya.Post(target, nil, challenge)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()

   return io.ReadAll(resp.Body)
}
