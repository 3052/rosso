package hboMax

import (
   "fmt"
   "net/http"
)

// Client handles communication with the HBO Max API.
type Client struct {
   BaseURL    string
   Token      string
   HTTPClient *http.Client
}

// NewClient creates a new API client with the given st token.
func NewClient(token string) *Client {
   return &Client{
      BaseURL:    "https://default.any-emea.prd.api.hbomax.com",
      Token:      token,
      HTTPClient: &http.Client{},
   }
}

// newRequest sets up the boilerplate headers required by the API.
func (c *Client) newRequest(method, endpoint string) (*http.Request, error) {
   req, err := http.NewRequest(method, c.BaseURL+endpoint, nil)
   if err != nil {
      return nil, err
   }

   req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:140.0) Gecko/20100101 Firefox/140.0")
   req.Header.Set("Accept", "application/json")
   req.Header.Set("Referer", "https://play.hbomax.com/")
   req.Header.Set("X-Device-Info", "hbomax/6.17.1 (desktop/desktop; Windows/NT 10.0; f681564c-1be5-4495-882b-6efc06cd8a9d/da0cdd94-5a39-42ef-aa68-54cbc1b852c3)")
   req.Header.Set("X-Disco-Client", "WEB:NT 10.0:hbomax:6.17.1")
   req.Header.Set("X-Disco-Params", "realm=bolt,bid=beam,features=ar")
   req.Header.Set("Cookie", fmt.Sprintf("st=%s", c.Token))

   return req, nil
}
