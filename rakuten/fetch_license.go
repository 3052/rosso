package rakuten

import (
   "io"
   "net/url"

   "41.neocities.org/maya"
)

func (s *StreamInfo) FetchLicense(challenge []byte) ([]byte, error) {
   target, err := url.Parse(s.LicenseUrl)
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
