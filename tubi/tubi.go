package tubi

import (
   "41.neocities.org/maya"
   "io"
   "net/url"
)

func AcquireLicense(server *LicenseServer, body []byte) ([]byte, error) {
   target, err := url.Parse(server.Url)
   if err != nil {
      return nil, err
   }

   resp, err := maya.Post(target, nil, body)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()

   return io.ReadAll(resp.Body)
}

func (v *VideoResource) GetManifest() (*url.URL, error) {
   return url.Parse(v.Manifest.Url)
}
