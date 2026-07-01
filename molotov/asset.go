// asset.go
package molotov

import (
   "encoding/json"
   "fmt"
   "net/http"
   "net/url"
)

type AssetResponse struct {
   Stream struct {
      URL string `json:"url"` // The MPD URL
   } `json:"stream"`
   DRM struct {
      LicenseURL string `json:"licenseUrl"`
      Token      string `json:"token"`
   } `json:"drm"`
}

// GetAsset retrieves the asset playback details using the auth and user contexts.
func GetAsset(assetID string, signinResp *SigninResponse) (*AssetResponse, error) {
   // Initialize the request with the base URL
   req, err := http.NewRequest("GET", "https://api-eu.fubo.tv/vapi/asset/v1", nil)
   if err != nil {
      return nil, err
   }
   // Properly build and encode the query string
   query := url.Values{}
   query.Add("id", assetID)
   query.Add("type", "vod")
   req.URL.RawQuery = query.Encode()
   // Set Headers
   req.Header.Set("x-forwarded-for", x_forwarded_for)
   // Accessing the unwrapped field directly
   req.Header.Set("x-application-id", "molotov")
   req.Header.Set("x-device-type", "desktop")
   req.Header.Set("x-device-app", "web")
   req.Header.Set("x-drm-scheme", "widevine")
   req.Header.Set("Authorization", "Bearer "+signinResp.AccessToken)
   // Execute request
   resp, err := doRequest(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()

   if resp.StatusCode != http.StatusOK {
      return nil, fmt.Errorf("get asset failed with status: %d", resp.StatusCode)
   }

   // Decode response
   var assetResp AssetResponse
   if err := json.NewDecoder(resp.Body).Decode(&assetResp); err != nil {
      return nil, err
   }

   return &assetResp, nil
}

func (*AssetResponse) CachePath() string {
   return "rosso/molotov/AssetResponse"
}
