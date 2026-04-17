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

func (s *Session) GetWidevine(manifest *Manifest, challenge []byte) ([]byte, error) {
   if manifest == nil {
      return nil, fmt.Errorf("a valid stream manifest is required to request a DRM license")
   }
   if manifest.DrmLicenseId == "" {
      return nil, fmt.Errorf("manifest does not contain a DRM license ID")
   }

   url := fmt.Sprintf("%s/kapi/licenses/widevine/%s", BaseUrl, manifest.DrmLicenseId)

   req, err := http.NewRequest("POST", url, bytes.NewBuffer(challenge))
   if err != nil {
      return nil, err
   }

   req.Header.Set("Authorization", "Bearer "+s.Jwt)
   req.Header.Set("User-Agent", UserAgent)
   req.Header.Set("X-Version", Xversion)

   resp, err := http.DefaultClient.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()

   if resp.StatusCode != http.StatusOK {
      return nil, fmt.Errorf("widevine license request failed with status: %d", resp.StatusCode)
   }

   return io.ReadAll(resp.Body)
}

const (
   BaseUrl   = "https://www.kanopy.com"
   UserAgent = "!"
   Xversion  = "!/!/!/!"
)

// Session represents an authenticated user context.
type Session struct {
   Jwt       string `json:"jwt"`
   VisitorId string `json:"visitorId"`
   UserId    int    `json:"userId"`
   UserRole  string `json:"userRole"`
}
