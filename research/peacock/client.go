// client.go
package peacock

import (
   "crypto/rand"
   "encoding/hex"
   "fmt"
   "io"
   "net/http"
   "time"
)

const sasBase = "https://sas.peacocktv.com"

func randomDeviceID() string {
   b := make([]byte, 8)
   _, _ = rand.Read(b)
   return hex.EncodeToString(b)
}

func randomUUID() string {
   b := make([]byte, 16)
   _, _ = rand.Read(b)
   b[6] = (b[6] & 0x0f) | 0x40
   b[8] = (b[8] & 0x3f) | 0x80
   return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:16])
}

// Client is the Peacock HTTP client.
type Client struct {
   HTTP     *http.Client
   DeviceID string
   Token    string // OAuth2 bearer token; populated after successful login

   // CertPEM and KeyPEM hold the mTLS client certificate and private key
   // required by the play.clients.peacocktv.com token endpoint.
   // Set these via LoadCertFiles or directly before calling ExchangeToken.
   CertPEM []byte
   KeyPEM  []byte
}

// NewClient returns a new Client. If deviceID is empty, a random one is generated.
func NewClient(deviceID string) *Client {
   if deviceID == "" {
      deviceID = randomDeviceID()
   }
   return &Client{
      HTTP:     &http.Client{Timeout: 30 * time.Second},
      DeviceID: deviceID,
   }
}

func (c *Client) newRequest(method, url string, body io.Reader) (*http.Request, error) {
   req, err := http.NewRequest(method, url, body)
   if err != nil {
      return nil, err
   }
   c.setCommonHeaders(req.Header)
   if c.Token != "" {
      req.Header.Set("x-token", "OAuth2 "+c.Token)
   }
   return req, nil
}

func (c *Client) setCommonHeaders(h http.Header) {
   h.Set("User-Agent", "Mozilla/5.0 (Linux; Android 12; sdk_gphone64_x86_64 Build/SE1A.220826.008; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/91.0.4472.114 Mobile Safari/537.36")
   h.Set("x-skyott-platform", "ANDROIDTV")
   h.Set("x-skyott-proposition", "NBCUOTT")
   h.Set("x-skyott-provider", "NBCU")
   h.Set("x-skyott-territory", "US")
   h.Set("x-skyott-activeterritory", "US")
   h.Set("x-skyott-language", "en-US")
   h.Set("x-skyott-device", "TV")
   h.Set("x-skyott-broadcastregions", "INPATTERN_US_CENTRAL")
   h.Set("x-skyint-requestid", randomUUID())
}
