package oldflix

import (
   "encoding/json"
   "fmt"
   "net/http"
   "net/url"
   "strings"
)

const BaseUrl = "https://oldflix-api.azurewebsites.net"

type Login struct {
   Token  string `json:"token"`
   Status int    `json:"status"`
}

func FetchLogin(username, password string) (*Login, error) {
   data := url.Values{}
   data.Set("username", username)
   data.Set("password", password)
   req, err := http.NewRequest("POST", BaseUrl+"/api/token", strings.NewReader(data.Encode()))
   if err != nil {
      return nil, err
   }
   req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
   req.Header.Set("User-Agent", "okhttp/4.12.0")
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   var loginResp Login
   if err := json.NewDecoder(resp.Body).Decode(&loginResp); err != nil {
      return nil, fmt.Errorf("failed to decode login response: %w", err)
   }
   if loginResp.Token == "" {
      return nil, fmt.Errorf("authentication failed, no token received")
   }
   return &loginResp, nil
}
