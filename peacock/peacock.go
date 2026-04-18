package peacock

import (
   "errors"
   "net/url"
)

func (e *Endpoint) GetManifest() (*url.URL, error) {
   return url.Parse(e.Url)
}

func (p *Playout) GetFastly() (*Endpoint, error) {
   for _, endpoint_data := range p.Asset.Endpoints {
      if endpoint_data.Cdn == "FASTLY" {
         return &endpoint_data, nil
      }
   }
   return nil, errors.New("FASTLY endpoint not found")
}

type Playout struct {
   Asset struct {
      Endpoints []Endpoint
   }
   Description string
   Protection  struct {
      LicenceAcquisitionUrl string
   }
}

type Endpoint struct {
   Cdn string
   Url string
}

const (
   sky_client  = "NBCU-ANDROID-v3"
   sky_key     = "JuLQgyFz9n89D9pxcN6ZWZXKWfgj2PNBUb32zybj"
   sky_version = "1.0"
)

// userToken is good for one day
type Token struct {
   Description string
   UserToken   string
}

var Territory = "US"
