package tubi

import "net/url"

func (v *VideoResource) GetManifest() (*url.URL, error) {
   return url.Parse(v.Manifest.Url)
}
