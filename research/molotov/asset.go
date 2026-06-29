// asset.go
package molotov

import (
   "encoding/json"
   "fmt"
   "net/http"
)

// GetAsset retrieves the MPD URL, the license URL, and the DRM auth token.
func GetAsset(assetID, accessToken, userID, profileID, deviceID, sessionID string) (string, string, string, error) {
   // Constructing URL using the specific asset ID (VOD_314017 in this case)
   url := fmt.Sprintf("https://api-eu.fubo.tv/vapi/asset/v1?id=%s&type=vod", assetID)

   req, err := http.NewRequest("GET", url, nil)
   if err != nil {
      return "", "", "", err
   }

   req.Header.Set("Authorization", "Bearer "+accessToken)
   req.Header.Set("x-user-id", userID)
   req.Header.Set("x-profile-id", profileID)
   req.Header.Set("x-device-id", deviceID)
   req.Header.Set("x-session-id", sessionID)
   req.Header.Set("x-device-group", "desktop")
   req.Header.Set("x-client-version", "6.12.0")
   req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:140.0) Gecko/20100101 Firefox/140.0")

   client := &http.Client{}
   resp, err := client.Do(req)
   if err != nil {
      return "", "", "", err
   }
   defer resp.Body.Close()

   if resp.StatusCode != http.StatusOK {
      return "", "", "", fmt.Errorf("get asset failed with status: %d", resp.StatusCode)
   }

   var assetResp AssetResponse
   if err := json.NewDecoder(resp.Body).Decode(&assetResp); err != nil {
      return "", "", "", err
   }

   mpdURL := assetResp.Stream.URL
   licenseURL := assetResp.DRM.LicenseURL
   dtAuthToken := assetResp.DRM.LicenseURLHeaders.XDtAuthToken

   return mpdURL, licenseURL, dtAuthToken, nil
}

type AssetResponse struct {
   Stream struct {
      URL string `json:"url"` // The MPD URL
   } `json:"stream"`
   DRM struct {
      LicenseURL        string `json:"licenseUrl"`
      LicenseURLHeaders struct {
         XDtAuthToken string `json:"x-dt-auth-token"`
      } `json:"licenseUrlHeaders"`
   } `json:"drm"`
}
