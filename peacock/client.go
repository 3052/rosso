package peacock

import (
   "crypto/rand"
   "crypto/tls"
   _ "embed"
   "encoding/hex"
   "fmt"
   "net/http"
   "net/http/cookiejar"
   "time"
)

const sasBase = "https://sas.peacocktv.com"

//go:embed cert.pem
var certPEM []byte

//go:embed key.pem
var keyPEM []byte

// mtlsClient returns an *http.Client configured with the embedded mTLS
// certificate, ProxyFromEnvironment, and the given timeout.
func mtlsClient(timeout time.Duration) (*http.Client, error) {
   cert, err := tls.X509KeyPair(certPEM, keyPEM)
   if err != nil {
      return nil, err
   }

   return &http.Client{
      Timeout: timeout,
      Transport: &http.Transport{
         Proxy: http.ProxyFromEnvironment,
         TLSClientConfig: &tls.Config{
            Certificates: []tls.Certificate{cert},
            MinVersion:   tls.VersionTLS12,
         },
      },
   }, nil
}

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
}

// NewClient returns a new Client. If deviceID is empty, a random one is generated.
// The HTTP client is configured with a cookie jar so session cookies persist
// across requests (e.g. from SignIn to OAuthAuthorize).
func NewClient(deviceID string) *Client {
   if deviceID == "" {
      deviceID = randomDeviceID()
   }
   jar, _ := cookiejar.New(nil)
   return &Client{
      HTTP: &http.Client{
         Timeout: 30 * time.Second,
         Jar:     jar,
      },
      DeviceID: deviceID,
   }
}

// client.go
