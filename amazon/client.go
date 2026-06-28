package amazon

import (
   "encoding/json"
   "log"
   "net/http"
)

const ( // API Hosts
   HostAmazonAPI = "https://api.amazon.com"
   HostATVPS     = "https://atv-ps.amazon.com"
   HostATVExt    = "https://atv-ext.amazon.com"
)

const DeviceID = "deviceID"

// the wrong DTID will fail the license request. if you change the DTID you
// need to relog
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

// doRequest wraps the http.Client Do method to log every outgoing request.
func doRequest(req *http.Request) (*http.Response, error) {
   log.Println(req.Method, req.URL)
   client := &http.Client{}
   return client.Do(req)
}

func marshal(value any) ([]byte, error) {
   return json.MarshalIndent(value, "", " ")
}

// Device represents the metadata for a supported hardware device.
type Device struct {
   Manufacturer  string
   Model         string
   SecurityLevel int
   DeviceTypeID  string
}
