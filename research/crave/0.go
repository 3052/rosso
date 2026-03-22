package crave

import (
   "io"
   "net/http"
   "net/url"
   "strings"
)

func zero(username, password string) (*http.Response, error) {
   var data = url.Values{
      "username":[]string{username},
      "password":[]string{password},
      "grant_type":[]string{"password"},
   }.Encode()
   var req http.Request
   req.Header = http.Header{}
   req.Method = "POST"
   req.URL = &url.URL{}
   req.URL.Host = "account.bellmedia.ca"
   req.URL.Path = "/api/login/v2.1"
   req.URL.Scheme = "https"
   req.Body = io.NopCloser(strings.NewReader(data))
   req.Header.Add("Content-Type", "application/x-www-form-urlencoded;charset=utf-8")
   req.Header.Add("Authorization", "Basic Y3JhdmUtd2ViOmRlZmF1bHQ=")
   return http.DefaultClient.Do(&req)
}
