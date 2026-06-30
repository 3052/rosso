// asset.go
package molotov

import (
   "encoding/json"
   "fmt"
   "net/http"
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
func GetAsset(assetID string, signinResp *SigninResponse, userResp *UserResponse) (*AssetResponse, error) {
   url := fmt.Sprintf("https://api-eu.fubo.tv/vapi/asset/v1?id=%s&type=vod", assetID)

   req, err := http.NewRequest("GET", url, nil)
   if err != nil {
      return nil, err
   }

   // Accessing the unwrapped field directly
   req.Header.Set("Authorization", "Bearer "+signinResp.AccessToken)

   req.Header.Set("x-user-id", userResp.ID)
   req.Header.Set("x-profile-id", userResp.Profiles[0].ID)

   req.Header.Set("x-device-id", DeviceID)
   req.Header.Set("x-application-id", "molotov")
   req.Header.Set("x-device-group", "desktop")
   req.Header.Set("x-device-type", "desktop")
   req.Header.Set("x-device-app", "web")
   req.Header.Set("x-client-version", "6.12.0")
   req.Header.Set("x-drm-scheme", "widevine")
   req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:140.0) Gecko/20100101 Firefox/140.0")

   resp, err := doRequest(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()

   if resp.StatusCode != http.StatusOK {
      return nil, fmt.Errorf("get asset failed with status: %d", resp.StatusCode)
   }

   var assetResp AssetResponse
   if err := json.NewDecoder(resp.Body).Decode(&assetResp); err != nil {
      return nil, err
   }

   return &assetResp, nil
}

func (*AssetResponse) CachePath() string {
   return "rosso/molotov/AssetResponse"
}
