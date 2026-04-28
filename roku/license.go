package roku

import (
   "io"
   "net/url"

   "41.neocities.org/maya"
)

func GetWidevineLicense(config *PlaybackConfig, challenge []byte) ([]byte, error) {
   target, err := url.Parse(config.Drm.Widevine.LicenseServer)
   if err != nil {
      return nil, err
   }
   headers := map[string]string{
      "content-type": "application/x-protobuf",
      "user-agent":   "Go-http-client/2.0",
   }

   resp, err := maya.Post(target, headers, challenge)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()

   return io.ReadAll(resp.Body)
}
