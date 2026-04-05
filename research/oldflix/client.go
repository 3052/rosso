package oldflix

import (
   "net/http"
   "time"
)

const BaseURL = "https://oldflix-api.azurewebsites.net"

// Client holds the HTTP client and the JWT authentication token
type Client struct {
   HTTPClient *http.Client
   Token      string
}

// NewClient creates a new configured Oldflix API client
func NewClient() *Client {
   return &Client{
      HTTPClient: &http.Client{Timeout: 15 * time.Second},
   }
}
