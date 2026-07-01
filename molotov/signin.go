// signin.go
package molotov

import (
   "bytes"
   "encoding/json"
   "fmt"
   "net/http"
)

// MolotovAgent represents the JSON structure required for Molotov.tv specific headers
type MolotovAgent struct {
   AppBuild          int      `json:"app_build"`
   AppID             string   `json:"app_id"`
   APIVersion        int      `json:"api_version"`
   Type              string   `json:"type"`
   OS                string   `json:"os"`
   Manufacturer      string   `json:"manufacturer"`
   Model             string   `json:"model"`
   Brand             string   `json:"brand"`
   Serial            string   `json:"serial"`
   FeaturesSupported []string `json:"features_supported"`
}

type SigninRequest struct {
   Username string `json:"username"`
   Password string `json:"password"`
}

type SigninResponse struct {
   AccessToken  string `json:"access_token"`
   RefreshToken string `json:"refresh_token"`
}

// 429 if you call this too many times
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

   // Unwrap the "payload" envelope layer (Specific to api-eu.fubo.tv endpoints)
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

// Refresh uses the Molotov v2 endpoint to obtain a new access and refresh token,
// overwriting the tokens in the receiver.
func (s *SigninResponse) Refresh() error {
   if s.RefreshToken == "" {
      return fmt.Errorf("missing refresh token")
   }

   url := fmt.Sprintf("https://www.molotov.tv/v2/auth/refresh/%s", s.RefreshToken)

   req, err := http.NewRequest("GET", url, nil)
   if err != nil {
      return err
   }

   // Construct the agent data required by molotov.tv endpoints
   agentData := MolotovAgent{
      AppBuild:          1,
      AppID:             "customer_area",
      APIVersion:        8,
      Type:              "desktop",
      OS:                "windows",
      FeaturesSupported: []string{"parental_control_v3", "allow_recurly", "evergreen"},
   }
   agentJSON, _ := json.Marshal(agentData)

   req.Header.Set("Accept", "application/json")
   req.Header.Set("Content-Type", "application/json")
   req.Header.Set("X-Molotov-Agent", string(agentJSON))
   req.Header.Set("X-Molotov-Website", "customer_area")
   req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:140.0) Gecko/20100101 Firefox/140.0")

   resp, err := doRequest(req)
   if err != nil {
      return err
   }
   defer resp.Body.Close()

   if resp.StatusCode != http.StatusOK {
      return fmt.Errorf("refresh failed with status: %d", resp.StatusCode)
   }

   // Decode directly into the receiver, clobbering existing fields
   if err := json.NewDecoder(resp.Body).Decode(s); err != nil {
      return err
   }

   return nil
}
