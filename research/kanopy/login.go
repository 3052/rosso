package kanopy

import (
   "bytes"
   "encoding/json"
   "fmt"
   "io"
   "net/http"
)

type LoginRequest struct {
   CredentialType string `json:"credentialType"`
   EmailUser      struct {
      Email    string `json:"email"`
      Password string `json:"password"`
   } `json:"emailUser"`
}

type LoginResponse struct {
   JWT    string `json:"jwt"`
   UserID int    `json:"userId"`
}

// Login authenticates the user and returns the login response containing the JWT and UserID.
func (c *Client) Login(email, password string) (*LoginResponse, error) {
   payload := LoginRequest{
      CredentialType: "email",
   }
   payload.EmailUser.Email = email
   payload.EmailUser.Password = password

   body, err := json.Marshal(payload)
   if err != nil {
      return nil, err
   }

   req, err := http.NewRequest("POST", BaseURL+"/kapi/login", bytes.NewBuffer(body))
   if err != nil {
      return nil, err
   }

   req.Header.Set("Content-Type", "application/json")
   req.Header.Set("User-Agent", c.UserAgent)

   // Explicitly using http.DefaultClient
   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()

   if resp.StatusCode != http.StatusOK {
      return nil, fmt.Errorf("login failed with status: %d", resp.StatusCode)
   }

   respBody, err := io.ReadAll(resp.Body)
   if err != nil {
      return nil, err
   }

   var loginResp LoginResponse
   if err := json.Unmarshal(respBody, &loginResp); err != nil {
      return nil, err
   }

   if loginResp.JWT != "" {
      c.SetToken(loginResp.JWT)
   }

   return &loginResp, nil
}
