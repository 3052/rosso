package amazon

const (
   // API Hosts
   HostAmazonAPI = "https://api.amazon.com"
   HostATVPS     = "https://atv-ps.amazon.com"
   HostATVExt    = "https://atv-ext.amazon.com"

   // the wrong DTID will fail the license request. if you change the DTID you
   // need to relog

   //DeviceID = "deviceID"
   DeviceID = "uuidd9b2d63da09b46c1952a8c7b3e1f0a4d"

   // time.is/Unix_time
   DeviceName = "1782591345"
)

// Devices is a slice containing all supported device configurations.
var Devices = []Device{
   {
      Manufacturer:  "Hisense",
      Model:         "HE55A7000EUWTS",
      SecurityLevel: 3000,
      DeviceTypeID:  "A3REWRVYBYPKUM",
   },
   {
      Manufacturer:  "Hisense",
      Model:         "HU50A6100UW",
      SecurityLevel: 3000,
      DeviceTypeID:  "AAJ692ZPT1X85",
   },
   {
      Manufacturer:  "Hisense",
      Model:         "HU32E5600FHWV",
      SecurityLevel: 3000,
      DeviceTypeID:  "A2RGJ95OVLR12U",
   },
   {
      Manufacturer:  "EXPRESS LUCK TECHNOLOGY LIMITED",
      Model:         "LE-*",
      SecurityLevel: 3000,
      DeviceTypeID:  "A3NM0WFSU3DLT5",
   },
}

// Device represents the metadata for a supported hardware device.
type Device struct {
   Manufacturer  string
   Model         string
   SecurityLevel int
   DeviceTypeID  string
}
