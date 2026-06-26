package amazon

import (
   "net/http"
)

const defaultAPIHost = "atv-ps.primevideo.com"
const defaultDeviceTypeID = "A3NM0WFSU3DLT5"

// Client handles the communication with Amazon APIs.
type Client struct {
   HTTPClient *http.Client
}

// NewClient creates a new Amazon API client.
func NewClient(httpClient *http.Client) *Client {
   if httpClient == nil {
      httpClient = http.DefaultClient
   }
   return &Client{HTTPClient: httpClient}
}

// DeviceProfile holds the specific capabilities and identities to test against the API.
type DeviceProfile struct {
   DeviceID          string
   DRMType           string // "Widevine" or "PlayReady"
   HDCPLevel         string // e.g. "2.1", "2.3"
   MaxResolution     string // e.g. "480p", "720p", "1080p", "1440p", "2160p"
   HDRFormats        string // e.g. "None", "HDR10", "DolbyVision"
   VideoCodec        string // e.g. "H264", "H265"
   BitrateAdaptation string // e.g. "CVBR", "CBR"
   AuthBearer        string // Required for authorization
}
