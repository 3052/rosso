package amazon

import (
   "bytes"
   "encoding/json"
   "fmt"
   "io"
   "net/http"
)

type RegisterResponse struct {
   Response struct {
      Success struct {
         Tokens struct {
            Bearer struct {
               AccessToken  string `json:"access_token"`
               RefreshToken string `json:"refresh_token"`
            } `json:"bearer"`
            MacDms struct {
               DevicePrivateKey string `json:"device_private_key"`
               AdpToken         string `json:"adp_token"`
            } `json:"mac_dms"`
         } `json:"tokens"`
      } `json:"success"`
   } `json:"response"`
}

// RegisterDevice exchanges the authorization_code and code_verifier for API access tokens and the ADP private key.
func RegisterDevice(authCode, codeVerifier, deviceSerial string) (string, string, string, string, error) {
   url := "https://api.amazon.com/auth/register"

   payload := map[string]interface{}{
      "auth_data": map[string]string{
         "use_global_authentication": "true",
         "authorization_code":        authCode,
         "code_verifier":             codeVerifier,
         "code_algorithm":            "SHA-256",
         "client_domain":             "DeviceLegacy",
         "client_id":                 "61643565316233333062326434653565616338613331646436393462656431372341314d50534c4643374c3541464b",
      },
      "registration_data": map[string]string{
         "domain":           "DeviceLegacy",
         "device_type":      "A1MPSLFC7L5AFK",
         "device_serial":    deviceSerial,
         "app_name":         "com.amazon.avod.thirdpartyclient",
         "app_version":      "458000357",
         "device_model":     "sdk_gphone_x86_64",
         "os_version":       "11",
         "software_version": "130050002",
      },
      "requested_token_type": []string{"bearer", "mac_dms"},
      "device_metadata": map[string]string{
         "device_os_family": "android",
         "device_type":      "A1MPSLFC7L5AFK",
         "device_serial":    deviceSerial,
         "manufacturer":     "Google",
         "model":            "sdk_gphone_x86_64",
         "os_version":       "30",
         "product":          "sdk_gphone_x86_64",
      },
   }

   body, _ := json.Marshal(payload)
   req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(body))
   if err != nil {
      return "", "", "", "", err
   }

   req.Header.Set("Content-Type", "application/json")
   req.Header.Set("User-Agent", "Dalvik/2.1.0 (Linux; U; Android 11; sdk_gphone_x86_64 Build/RSR1.240422.006)")

   client := &http.Client{}
   resp, err := client.Do(req)
   if err != nil {
      return "", "", "", "", err
   }
   defer resp.Body.Close()

   respBody, err := io.ReadAll(resp.Body)
   if err != nil {
      return "", "", "", "", fmt.Errorf("failed to read response body: %v", err)
   }

   if resp.StatusCode != http.StatusOK {
      return "", "", "", "", fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, string(respBody))
   }

   var regResp RegisterResponse
   if err := json.Unmarshal(respBody, &regResp); err != nil {
      return "", "", "", "", fmt.Errorf("failed to parse JSON response: %v\nBody: %s", err, string(respBody))
   }

   accessToken := regResp.Response.Success.Tokens.Bearer.AccessToken
   if accessToken == "" {
      return "", "", "", "", fmt.Errorf("access_token is empty in the response.\nRaw Body: %s", string(respBody))
   }

   refreshToken := regResp.Response.Success.Tokens.Bearer.RefreshToken
   if refreshToken == "" {
      return "", "", "", "", fmt.Errorf("refresh_token is empty in the response.\nRaw Body: %s", string(respBody))
   }

   privateKey := regResp.Response.Success.Tokens.MacDms.DevicePrivateKey
   if privateKey == "" {
      return "", "", "", "", fmt.Errorf("device_private_key (mac_dms) is empty in the response.\nRaw Body: %s", string(respBody))
   }

   adpToken := regResp.Response.Success.Tokens.MacDms.AdpToken
   if adpToken == "" {
      return "", "", "", "", fmt.Errorf("adp_token is empty in the response.\nRaw Body: %s", string(respBody))
   }

   return accessToken, refreshToken, privateKey, adpToken, nil
}
