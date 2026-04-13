package itv

import (
   "io"
   "net/http"
   "net/url"
   "strings"
)

func (m *MediaFile) FetchDash() (*Dash, error) {
   resp, err := http.Get(strings.Replace(m.Href, "itvpnpctv", "itvpnpdotcom", 1))
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   body, err := io.ReadAll(resp.Body)
   if err != nil {
      return nil, err
   }
   return &Dash{Body: body, Url: resp.Request.URL}, nil
}

type Dash struct {
   Body []byte
   Url  *url.URL
}
