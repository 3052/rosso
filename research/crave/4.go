package crave

import (
   "net/http"
   "net/url"
)

const (
   BaseURL = "https://account.bellmedia.ca"
   // Basic base64("crave-web:default")
   BasicAuth = "Basic Y3JhdmUtd2ViOmRlZmF1bHQ="
   UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:140.0) Gecko/20100101 Firefox/140.0"
)

type TokenResponse struct {
   AccessToken  string `json:"access_token"`
   RefreshToken string `json:"refresh_token"`
   AccountID    string `json:"account_id,omitempty"`
   ExpiresIn    int    `json:"expires_in"`
}

func (t *TokenResponse) four() (*http.Response, error) {
   var req http.Request
   req.Header = http.Header{}
   req.URL = &url.URL{}
   req.URL.Host = "stream.video.9c9media.com"
   req.URL.Path = "/meta/content/938361/contentpackage/8143402/destination/1880/platform/1"
   value := url.Values{}
   value["format"] = []string{"mpd"}
   req.Header.Add("Authorization", "Bearer " + t.AccessToken)
   req.URL.RawQuery = value.Encode()
   req.URL.Scheme = "https"
   return http.DefaultClient.Do(&req)
}
