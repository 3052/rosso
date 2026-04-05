package oldflix

import (
   "encoding/json"
   "fmt"
   "net/http"
   "net/url"
   "strings"
)

type LoginResponse struct {
   Token  string `json:"token"`
   Status int    `json:"status"`
}

// Login authenticates with the API and stores the JWT token in the Client
func (c *Client) Login(username, password string) error {
   data := url.Values{}
   data.Set("username", username)
   data.Set("password", password)

   req, err := http.NewRequest("POST", BaseURL+"/api/token", strings.NewReader(data.Encode()))
   if err != nil {
      return err
   }
   req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
   req.Header.Set("User-Agent", "okhttp/4.12.0")

   resp, err := c.HTTPClient.Do(req)
   if err != nil {
      return err
   }
   defer resp.Body.Close()

   var loginResp LoginResponse
   if err := json.NewDecoder(resp.Body).Decode(&loginResp); err != nil {
      return fmt.Errorf("failed to decode login response: %w", err)
   }

   if loginResp.Token == "" {
      return fmt.Errorf("authentication failed, no token received")
   }

   c.Token = loginResp.Token
   return nil
}
