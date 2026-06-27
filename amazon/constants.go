package amazon

import (
   "crypto/rand"
   "encoding/hex"
   "fmt"
)

const (
   /////////////////////////////////////////////////////////////////////////////////
   DeviceTypeID = "A3NM0WFSU3DLT5"
   /////////////////////////////////////////////////////////////////////////////////
   // API Hosts
   HostAmazonAPI = "https://api.amazon.com"
   HostATVPS     = "https://atv-ps.amazon.com"
   HostATVExt    = "https://atv-ext.amazon.com"
)

// DeviceID represents a unique device identifier for Amazon API requests.
type DeviceID string

// NewDeviceID generates a structurally valid UUID v4 for device registration.
// The caller must save this value and pass it to subsequent API calls.
func NewDeviceID() (DeviceID, error) {
   var uuid [16]byte
   if _, err := rand.Read(uuid[:]); err != nil {
      return "", fmt.Errorf("failed to generate random bytes for device ID: %w", err)
   }

   // Apply standard UUID v4 bitmasks
   uuid[6] = (uuid[6] & 0x0f) | 0x40 // Version 4
   uuid[8] = (uuid[8] & 0x3f) | 0x80 // Variant RFC4122

   return DeviceID("uuid" + hex.EncodeToString(uuid[:])), nil
}

func (DeviceID) CachePath() string {
   return "rosso/amazon/DeviceID"
}
