package molotov

import (
   "bytes"
   "encoding/json"
   "fmt"
   "io"
   "log"
   "net/http"
   "net/url"
)

// DeviceID is the centralized value used for the x-device-id header across all requests.
const DeviceID = "x-device-id"

const x_forwarded_for = "178.132.106.134"

// doRequest logs the method and URL, then performs the HTTP request.
func doRequest(req *http.Request) (*http.Response, error) {
   log.Println(req.Method, req.URL)
   client := &http.Client{}
   return client.Do(req)
}

type AssetResponse struct {
   Stream struct {
      URL string `json:"url"` // The MPD URL
   } `json:"stream"`
   DRM struct {
      LicenseURL string `json:"licenseUrl"`
      Token      string `json:"token"`
   } `json:"drm"`
}

// GetAsset retrieves the asset playback details using the auth and user contexts.
func GetAsset(assetID string, signinResp *SigninResponse) (*AssetResponse, error) {
   // Initialize the request with the base URL
   req, err := http.NewRequest("GET", "https://api-eu.fubo.tv/vapi/asset/v1", nil)
   if err != nil {
      return nil, err
   }
   // Properly build and encode the query string
   query := url.Values{}
   query.Set("id", assetID)
   query.Set("type", "vod")
   req.URL.RawQuery = query.Encode()
   // Set Headers
   req.Header.Set("x-forwarded-for", x_forwarded_for)
   // Accessing the unwrapped field directly
   req.Header.Set("x-application-id", "molotov")
   req.Header.Set("x-device-type", "desktop")
   req.Header.Set("x-device-app", "web")
   req.Header.Set("x-drm-scheme", "widevine")
   req.Header.Set("Authorization", "Bearer "+signinResp.AccessToken)
   // Execute request
   resp, err := doRequest(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()

   if resp.StatusCode != http.StatusOK {
      return nil, fmt.Errorf("get asset failed with status: %d", resp.StatusCode)
   }

   // Decode response
   var assetResp AssetResponse
   if err := json.NewDecoder(resp.Body).Decode(&assetResp); err != nil {
      return nil, err
   }

   return &assetResp, nil
}

func (*AssetResponse) CachePath() string {
   return "rosso/molotov/AssetResponse"
}

// GetLicense requests the DRM license. As a method on *AssetResponse,
// it can be used directly as a closure: func([]byte) ([]byte, error).
func (a *AssetResponse) GetLicense(challenge []byte) ([]byte, error) {
   req, err := http.NewRequest("POST", a.DRM.LicenseURL, bytes.NewReader(challenge))
   if err != nil {
      return nil, err
   }
   resp, err := doRequest(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()

   if resp.StatusCode != http.StatusOK {
      return nil, fmt.Errorf("license request failed with status: %d", resp.StatusCode)
   }

   licenseData, err := io.ReadAll(resp.Body)
   if err != nil {
      return nil, err
   }

   return licenseData, nil
}

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

// molotov.go
