package kanopy

import "net/url"

func (m *Manifest) GetManifest() (*url.URL, error) {
   return url.Parse(m.URL)
}
