// refresh_token.go
package peacock

import (
   "bytes"
   "crypto/md5"
   "crypto/tls"
   "encoding/hex"
   "encoding/json"
   "fmt"
   "net/http"
)

// RefreshTokenResponse is the response from PATCH /auth/tokens.
type RefreshTokenResponse struct {
   UserToken                     string `json:"userToken"`
   TokenExpiryTime               string `json:"tokenExpiryTime"`
   RecommendedTokenReacquireTime string `json:"recommendedTokenReacquireTime"`
}

type refreshTokenDevice struct {
   ID string `json:"id"`
}

type refreshTokenRequest struct {
   Device refreshTokenDevice `json:"device"`
}

// RefreshToken refreshes an existing user token using the mTLS certificate
// at the given paths. The current userToken is required.
func (c *Client) RefreshToken(userToken, certPath, keyPath string) (*RefreshTokenResponse, error) {
   if userToken == "" {
      return nil, fmt.Errorf("refresh token: empty userToken")
   }
   if certPath == "" || keyPath == "" {
      return nil, fmt.Errorf("refresh token: certPath and keyPath must be set")
   }

   body := refreshTokenRequest{
      Device: refreshTokenDevice{
         ID: c.DeviceID,
      },
   }

   raw, err := json.Marshal(body)
   if err != nil {
      return nil, fmt.Errorf("refresh token: marshal: %w", err)
   }

   hash := md5.Sum(raw)
   contentMD5 := hex.EncodeToString(hash[:])

   cert, err := tls.LoadX509KeyPair(certPath, keyPath)
   if err != nil {
      return nil, fmt.Errorf("refresh token: load key pair: %w", err)
   }

   mtlsClient := &http.Client{
      Timeout: c.HTTP.Timeout,
      Transport: &http.Transport{
         Proxy: http.ProxyFromEnvironment,
         TLSClientConfig: &tls.Config{
            Certificates: []tls.Certificate{cert},
            MinVersion:   tls.VersionTLS12,
         },
      },
   }

   req, err := http.NewRequest(http.MethodPatch, playBase+"/auth/tokens", bytes.NewReader(raw))
   if err != nil {
      return nil, fmt.Errorf("refresh token: create request: %w", err)
   }

   req.Header.Set("Accept", "application/vnd.tokens.v1+json")
   req.Header.Set("Content-Type", "application/vnd.tokens.v1+json")
   req.Header.Set("Content-MD5", contentMD5)
   req.Header.Set("Origin", "https://tv.clients.peacocktv.com")
   req.Header.Set("User-Agent", "Mozilla/5.0 (Linux; Android 12; sdk_gphone64_x86_64 Build/SE1A.220826.008; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/91.0.4472.114 Mobile Safari/537.36")
   req.Header.Set("x-skyott-activeterritory", "US")
   req.Header.Set("x-skyott-broadcastregions", "INPATTERN_US_CENTRAL")
   req.Header.Set("x-skyott-device", "TV")
   req.Header.Set("x-skyott-language", "en-US")
   req.Header.Set("x-skyott-platform", "ANDROIDTV")
   req.Header.Set("x-skyott-proposition", "NBCUOTT")
   req.Header.Set("x-skyott-provider", "NBCU")
   req.Header.Set("x-skyott-territory", "US")
   req.Header.Set("x-skyott-usertoken", userToken)

   resp, err := mtlsClient.Do(req)
   if err != nil {
      return nil, fmt.Errorf("refresh token: %w", err)
   }
   defer resp.Body.Close()

   if resp.StatusCode != http.StatusOK {
      return nil, fmt.Errorf("refresh token: bad status %d", resp.StatusCode)
   }

   var out RefreshTokenResponse
   if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
      return nil, fmt.Errorf("refresh token: decode: %w", err)
   }
   return &out, nil
}
