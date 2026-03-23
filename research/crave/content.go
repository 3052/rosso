package crave

import (
   "net/http"
   "net/url"
)

// https://crave.ca/movie/goldeneye-38860
func (t *TokenResponse) content() (*http.Response, error) {
   var req http.Request
   req.Header = http.Header{}
   req.URL = &url.URL{}
   req.URL.Host = "stream.video.9c9media.com"
   req.URL.Path = "/meta/content/938361/contentpackage/8143402/destination/1880/platform/1"
   req.Header.Add("Authorization", "Bearer "+t.AccessToken)
   value := url.Values{}
   value["format"] = []string{"mpd"}
   req.URL.RawQuery = value.Encode()
   req.URL.Scheme = "https"
   return http.DefaultClient.Do(&req)
}
