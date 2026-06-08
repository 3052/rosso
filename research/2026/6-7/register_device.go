// register_device.go
package amazon

import (
   "bytes"
   "encoding/json"
   "fmt"
   "io"
   "net/http"
)

type RegisterDeviceRequest struct {
   AuthData            AuthData          `json:"auth_data"`
   RegistrationData    RegistrationData  `json:"registration_data"`
   RequestedTokenType  []string          `json:"requested_token_type"`
   Cookies             CookiesData       `json:"cookies"`
   AgeInfo             map[string]any    `json:"age_info"`
   UserContextMap      map[string]string `json:"user_context_map,omitempty"`
   DeviceMetadata      DeviceMetadata    `json:"device_metadata"`
   RequestedExtensions []string          `json:"requested_extensions"`
}

type AuthData struct {
   UseGlobalAuthentication string `json:"use_global_authentication"`
   AuthorizationCode       string `json:"authorization_code"`
   CodeVerifier            string `json:"code_verifier"`
   CodeAlgorithm           string `json:"code_algorithm"`
   ClientDomain            string `json:"client_domain"`
   ClientId                string `json:"client_id"`
}

type RegistrationData struct {
   Domain          string `json:"domain"`
   DeviceType      string `json:"device_type"`
   DeviceSerial    string `json:"device_serial"`
   AppName         string `json:"app_name"`
   AppVersion      string `json:"app_version"`
   DeviceModel     string `json:"device_model"`
   OsVersion       string `json:"os_version"`
   SoftwareVersion string `json:"software_version"`
}

type CookiesData struct {
   Domain         string   `json:"domain"`
   WebsiteCookies []string `json:"website_cookies"`
}

type DeviceMetadata struct {
   DeviceOsFamily string `json:"device_os_family"`
   DeviceType     string `json:"device_type"`
   DeviceSerial   string `json:"device_serial"`
   Manufacturer   string `json:"manufacturer"`
   Model          string `json:"model"`
   OsVersion      string `json:"os_version"`
   Product        string `json:"product"`
}

type RegisterDeviceResponse struct {
   Response struct {
      Success struct {
         Tokens struct {
            Bearer struct {
               AccessToken  string `json:"access_token"`
               RefreshToken string `json:"refresh_token"`
               ExpiresIn    string `json:"expires_in"`
            } `json:"bearer"`
            MacDms struct {
               DevicePrivateKey string `json:"device_private_key"`
               AdpToken         string `json:"adp_token"`
            } `json:"mac_dms"`
         } `json:"tokens"`
      } `json:"success"`
   } `json:"response"`
}

func RegisterDevice(client *http.Client, authCode, codeVerifier, deviceSerial string) (*RegisterDeviceResponse, error) {
   reqBody := RegisterDeviceRequest{
      AuthData: AuthData{
         UseGlobalAuthentication: "true",
         AuthorizationCode:       authCode,
         CodeVerifier:            codeVerifier,
         CodeAlgorithm:           "SHA-256",
         ClientDomain:            "DeviceLegacy",
         ClientId:                "61643565316233333062326434653565616338613331646436393462656431372341314d50534c4643374c3541464b",
      },
      RegistrationData: RegistrationData{
         Domain:          "DeviceLegacy",
         DeviceType:      "A1MPSLFC7L5AFK",
         DeviceSerial:    deviceSerial,
         AppName:         "com.amazon.avod.thirdpartyclient",
         AppVersion:      "458000357",
         DeviceModel:     "sdk_gphone_x86_64",
         OsVersion:       "google/sdk_gphone_x86_64/generic_x86_64_arm64:11/RSR1.240422.006/12134477:userdebug/dev-keys",
         SoftwareVersion: "130050002",
      },
      RequestedTokenType: []string{"bearer", "mac_dms", "store_authentication_cookie", "website_cookies"},
      Cookies: CookiesData{
         Domain:         "amazon.com",
         WebsiteCookies: []string{},
      },
      AgeInfo:        make(map[string]any),
      UserContextMap: make(map[string]string),
      DeviceMetadata: DeviceMetadata{
         DeviceOsFamily: "android",
         DeviceType:     "A1MPSLFC7L5AFK",
         DeviceSerial:   deviceSerial,
         Manufacturer:   "Google",
         Model:          "sdk_gphone_x86_64",
         OsVersion:      "30",
         Product:        "sdk_gphone_x86_64",
      },
      RequestedExtensions: []string{"device_info", "customer_info"},
   }

   jsonData, err := json.Marshal(reqBody)
   if err != nil {
      return nil, fmt.Errorf("failed to marshal request: %w", err)
   }

   req, err := http.NewRequest("POST", "https://api.amazon.com/auth/register", bytes.NewBuffer(jsonData))
   if err != nil {
      return nil, fmt.Errorf("failed to create request: %w", err)
   }

   req.Header.Set("Content-Type", "application/json")
   req.Header.Set("Accept-Language", "en-US")
   req.Header.Set("User-Agent", "Dalvik/2.1.0 (Linux; U; Android 11; sdk_gphone_x86_64 Build/RSR1.240422.006)")
   req.Header.Set("x-amzn-identity-auth-domain", "api.amazon.com")

   resp, err := client.Do(req)
   if err != nil {
      return nil, fmt.Errorf("request failed: %w", err)
   }
   defer resp.Body.Close()

   bodyBytes, err := io.ReadAll(resp.Body)
   if err != nil {
      return nil, fmt.Errorf("failed to read response: %w", err)
   }

   var result RegisterDeviceResponse
   if err := json.Unmarshal(bodyBytes, &result); err != nil {
      return nil, fmt.Errorf("failed to decode response: %w", err)
   }

   return &result, nil
}
