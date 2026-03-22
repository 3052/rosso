package crave

import (
   "net/http"
   "time"
)

const (
   BaseURL = "https://account.bellmedia.ca"
   // Basic base64("crave-web:default")
   BasicAuth = "Basic Y3JhdmUtd2ViOmRlZmF1bHQ="
   UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:140.0) Gecko/20100101 Firefox/140.0"
)

type Client struct {
   HTTPClient *http.Client
}

func NewClient() *Client {
   return &Client{
      HTTPClient: &http.Client{
         Timeout: 15 * time.Second,
      },
   }
}

type TokenResponse struct {
   AccessToken  string `json:"access_token"`
   RefreshToken string `json:"refresh_token"`
   AccountID    string `json:"account_id,omitempty"`
   ExpiresIn    int    `json:"expires_in"`
}
