package crave

import (
   "io"
   "net/http"
   "net/url"
   "strings"
)

func (a *account) magic_link_token() (string, error) {
   var req http.Request
   req.Method = "POST"
   req.URL = &url.URL{
      Scheme: "https",
      Host: "account.bellmedia.ca",
      Path: "/api/magic-link/v2.1/generate",
   }
   req.Header = http.Header{}
   req.Header.Set("authorization", "Bearer " + a.AccessToken)
   resp, err := http.DefaultClient.Do(&req)
   if err != nil {
      return "", err
   }
   defer resp.Body.Close()
   var data strings.Builder
   _, err = io.Copy(&data, resp.Body)
   if err != nil {
      return "", err
   }
   return data.String(), nil
}
