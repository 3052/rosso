// client_id.go
package amazon

import (
   "encoding/hex"
   "fmt"
)

// GenerateClientID creates the OAuth 2.0 client_id expected by Amazon's sign-in page.
// The authDeviceType for the Amazon Video Android app is "A1MPSLFC7L5AFK".
func GenerateClientID(deviceID, authDeviceType string) string {
   plainText := fmt.Sprintf("%s#%s", deviceID, authDeviceType)
   hexEncoded := hex.EncodeToString([]byte(plainText))
   return "device:" + hexEncoded
}
