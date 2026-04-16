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

// Login authenticates the user and returns an initialized Session.
func Login(email, password string) (*Session, error) {
   payload := LoginRequest{
      CredentialType: "email",
   }
   payload.EmailUser.Email = email
   payload.EmailUser.Password = password

   body, err := json.Marshal(payload)
   if err != nil {
      return nil, err
   }

   req, err := http.NewRequest("POST", BaseUrl+"/kapi/login", bytes.NewBuffer(body))
   if err != nil {
      return nil, err
   }

   req.Header.Set("Content-Type", "application/json")
   req.Header.Set("User-Agent", UserAgent)

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

   var session Session
   if err := json.Unmarshal(respBody, &session); err != nil {
      return nil, err
   }

   return &session, nil
}
