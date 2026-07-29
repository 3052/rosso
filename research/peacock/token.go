// token.go
package peacock

import (
   "bytes"
   "crypto/md5"
   "crypto/tls"
   "encoding/hex"
   "encoding/json"
   "fmt"
   "net/http"
   "os"
   "time"
)

const playBase = "https://play.clients.peacocktv.com"

// TokenResponse is the response from POST /auth/throttled/tokens.
type TokenResponse struct {
   UserToken                     string    `json:"userToken"`
   TokenExpiryTime               time.Time `json:"tokenExpiryTime"`
   RecommendedTokenReacquireTime time.Time `json:"recommendedTokenReacquireTime"`
}

type tokenAuth struct {
   AuthScheme        string `json:"authScheme"`
   AuthIssuer        string `json:"authIssuer"`
   Provider          string `json:"provider"`
   ProviderTerritory string `json:"providerTerritory"`
   Proposition       string `json:"proposition"`
   AuthToken         string `json:"authToken"`
}

type tokenDevice struct {
   Type        string `json:"type"`
   Platform    string `json:"platform"`
   ID          string `json:"id"`
   DrmDeviceID string `json:"drmDeviceId"`
}

// tokenRequest is the body sent to POST /auth/throttled/tokens.
type tokenRequest struct {
   Auth   tokenAuth   `json:"auth"`
   Device tokenDevice `json:"device"`
}

// ExchangeToken trades the OAuth2 activation token for a long-lived user token
// using the mTLS certificate configured on the Client.
// The returned user token can be used as a bearer credential for playback.
func (c *Client) ExchangeToken() (*TokenResponse, error) {
   if c.Token == "" {
      return nil, fmt.Errorf("exchange token: no activation token; call Activate first")
   }
   if len(c.CertPEM) == 0 || len(c.KeyPEM) == 0 {
      return nil, fmt.Errorf("exchange token: mTLS certificate not configured; set CertPEM and KeyPEM on Client")
   }

   body := tokenRequest{
      Auth: tokenAuth{
         AuthScheme:        "OAUTH",
         AuthIssuer:        "NOWTV",
         Provider:          "NBCU",
         ProviderTerritory: "US",
         Proposition:       "NBCUOTT",
         AuthToken:         c.Token,
      },
      Device: tokenDevice{
         Type:        "TV",
         Platform:    "ANDROIDTV",
         ID:          c.DeviceID,
         DrmDeviceID: "UNKNOWN",
      },
   }

   raw, err := json.Marshal(body)
   if err != nil {
      return nil, fmt.Errorf("exchange token: marshal: %w", err)
   }

   // Compute Content-MD5 of the request body.
   hash := md5.Sum(raw)
   contentMD5 := hex.EncodeToString(hash[:])

   // Build an mTLS HTTP client using the loaded certificate.
   cert, err := tls.X509KeyPair(c.CertPEM, c.KeyPEM)
   if err != nil {
      return nil, fmt.Errorf("exchange token: load key pair: %w", err)
   }

   mtlsClient := &http.Client{
      Timeout: c.HTTP.Timeout,
      Transport: &http.Transport{
         TLSClientConfig: &tls.Config{
            Certificates: []tls.Certificate{cert},
            MinVersion:   tls.VersionTLS12,
         },
      },
   }

   req, err := c.newRequest(http.MethodPost, playBase+"/auth/throttled/tokens", bytes.NewReader(raw))
   if err != nil {
      return nil, err
   }
   req.Header.Set("Content-Type", "application/vnd.tokens.v1+json")
   req.Header.Set("Accept", "application/vnd.tokens.v1+json")
   req.Header.Set("Content-MD5", contentMD5)
   req.Header.Set("Origin", "https://tv.clients.peacocktv.com")

   resp, err := mtlsClient.Do(req)
   if err != nil {
      return nil, fmt.Errorf("exchange token: %w", err)
   }
   defer resp.Body.Close()

   if resp.StatusCode != http.StatusOK {
      return nil, fmt.Errorf("exchange token: bad status %d", resp.StatusCode)
   }

   var out TokenResponse
   if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
      return nil, fmt.Errorf("exchange token: decode: %w", err)
   }
   return &out, nil
}

// LoadCertFiles reads the mTLS certificate and key from the given file paths
// and stores them on the Client for use by ExchangeToken.
func (c *Client) LoadCertFiles(certPath, keyPath string) error {
   certPEM, err := os.ReadFile(certPath)
   if err != nil {
      return fmt.Errorf("read cert file: %w", err)
   }
   keyPEM, err := os.ReadFile(keyPath)
   if err != nil {
      return fmt.Errorf("read key file: %w", err)
   }
   c.CertPEM = certPEM
   c.KeyPEM = keyPEM
   return nil
}
