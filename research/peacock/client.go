// client.go
package peacock

import (
   "crypto/rand"
   "encoding/hex"
   "fmt"
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
