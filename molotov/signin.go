// signin.go
package molotov

import (
   "bytes"
   "encoding/json"
   "fmt"
   "net/http"
)

type SigninRequest struct {
   Username string `json:"username"`
   Password string `json:"password"`
}

type SigninResponse struct {
   AccessToken  string `json:"access_token"`
   RefreshToken string `json:"refresh_token"`
}

// too many calls gets 429
func Signin(username, password string) (*SigninResponse, error) {
   url := "https://api-eu.fubo.tv/v2/signin"
   reqBody, err := json.Marshal(SigninRequest{
      Username: username,
      Password: password,
   })
   if err != nil {
      return nil, err
   }
   req, err := http.NewRequest("PUT", url, bytes.NewBuffer(reqBody))
   if err != nil {
      return nil, err
   }
   req.Header.Set("x-device-id", DeviceID)
   resp, err := doRequest(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()

   if resp.StatusCode != http.StatusOK {
      return nil, fmt.Errorf("signin failed with status: %d", resp.StatusCode)
   }

   // Unwrap the "payload" envelope layer
   var envelope struct {
      Payload SigninResponse `json:"payload"`
   }
   if err := json.NewDecoder(resp.Body).Decode(&envelope); err != nil {
      return nil, err
   }

   return &envelope.Payload, nil
}

func (*SigninResponse) CachePath() string {
   return "rosso/molotov/SigninResponse"
}

// Refresh uses the Fubo API endpoint to obtain a new access and refresh token,
// overwriting the tokens in the receiver.
func (s *SigninResponse) Refresh() error {
   if s.RefreshToken == "" {
      return fmt.Errorf("missing refresh token")
   }
   url := "https://api-eu.fubo.tv/refresh"
   // Request has no body (content-length: 0 in the trace)
   req, err := http.NewRequest("POST", url, nil)
   if err != nil {
      return err
   }
   // The refresh token goes in the Authorization header
   req.Header.Set("Authorization", "Bearer "+s.RefreshToken)
   resp, err := doRequest(req)
   if err != nil {
      return err
   }
   defer resp.Body.Close()

   if resp.StatusCode != http.StatusOK {
      return fmt.Errorf("refresh failed with status: %d", resp.StatusCode)
   }

   // Unlike the /signin endpoint, /refresh returns the tokens directly at the root.
   // Decoding directly into `s` clobbers the existing token values.
   if err := json.NewDecoder(resp.Body).Decode(s); err != nil {
      return err
   }

   return nil
}
