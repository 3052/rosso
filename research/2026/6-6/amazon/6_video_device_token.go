package amazon

import (
   "bytes"
   "crypto"
   "crypto/rand"
   "crypto/rsa"
   "crypto/sha256"
   "crypto/x509"
   "encoding/base64"
   "encoding/json"
   "encoding/pem"
   "fmt"
   "net/http"
   "time"
)

type VideoTokenResponse struct {
   AccessToken  string `json:"access_token"`
   RefreshToken string `json:"refresh_token"`
}

// GetVideoDeviceToken uses ADP (Amazon Device Provisioning) signature to request a token for the Video app device type (A43PXU4ZN2AL1).
func GetVideoDeviceToken(deviceID, adpToken, privateKeyPEM string) (string, string, error) {
   url := "https://api.amazon.com/auth/token"

   payload := map[string]interface{}{
      "app_name":             "com.amazon.avod.thirdpartyclient",
      "app_version":          "130050002",
      "source_token_type":    "dms_token",
      "source_token":         "source_token",
      "requested_token_type": "refresh_token",
      "device_metadata": map[string]string{
         "device_os_family": "android",
         "device_type":      "A43PXU4ZN2AL1",
         "device_serial":    deviceID,
         "manufacturer":     "Google",
         "model":            "sdk_gphone_x86_64",
         "os_version":       "30",
         "product":          "sdk_gphone_x86_64",
      },
      "map_version": map[string]interface{}{
         "current_version":           "20251126N",
         "package_name":              "com.amazon.avod.thirdpartyclient",
         "platform":                  "Android",
         "client_metrics_integrated": true,
      },
      "age_info": map[string]string{},
   }

   body, _ := json.Marshal(payload)
   req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(body))
   if err != nil {
      return "", "", err
   }

   // Format ISO8601 Timestamp
   timestamp := time.Now().UTC().Format("2006-01-02T15:04:05Z")

   // Generate ADP Signature
   stringToSign := fmt.Sprintf("POST\napi.amazon.com\n/auth/token\n%s\n%s", timestamp, string(body))

   block, _ := pem.Decode([]byte(privateKeyPEM))
   if block == nil {
      return "", "", fmt.Errorf("failed to parse PEM block containing the private key")
   }
   privKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
   if err != nil {
      return "", "", fmt.Errorf("failed to parse RSA private key: %v", err)
   }

   hashed := sha256.Sum256([]byte(stringToSign))
   signatureBytes, err := rsa.SignPKCS1v15(rand.Reader, privKey, crypto.SHA256, hashed[:])
   if err != nil {
      return "", "", fmt.Errorf("failed to sign ADP challenge: %v", err)
   }

   signature := base64.StdEncoding.EncodeToString(signatureBytes)

   req.Header.Set("Content-Type", "application/json")
   req.Header.Set("x-amzn-identity-auth-domain", "api.amazon.com")
   req.Header.Set("x-adp-alg", "SHA256WithRSA:1.0")
   req.Header.Set("x-adp-token", adpToken)
   req.Header.Set("x-adp-signature", fmt.Sprintf("%s:%s", signature, timestamp))
   req.Header.Set("User-Agent", "AmazonWebView/MAPClientLib/130050002/Android/11/sdk_gphone_x86_64")

   client := &http.Client{}
   resp, err := client.Do(req)
   if err != nil {
      return "", "", err
   }
   defer resp.Body.Close()

   if resp.StatusCode != http.StatusOK {
      buf := new(bytes.Buffer)
      buf.ReadFrom(resp.Body)
      return "", "", fmt.Errorf("expected 200 OK, got status code: %d, body: %s", resp.StatusCode, buf.String())
   }

   var tokenResp VideoTokenResponse
   if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
      return "", "", err
   }

   if tokenResp.AccessToken == "" || tokenResp.RefreshToken == "" {
      return "", "", fmt.Errorf("received 200 OK, but access_token or refresh_token was empty")
   }

   return tokenResp.AccessToken, tokenResp.RefreshToken, nil
}
