// get_profiles.go
package amazon

import (
   "encoding/json"
   "fmt"
   "io"
   "net/http"
   "net/url"
)

type ProfilesResponse struct {
   Resource struct {
      Profiles []struct {
         ProfileId       string `json:"profileId"`
         Name            string `json:"name"`
         ProfileAgeGroup string `json:"profileAgeGroup"`
      } `json:"profiles"`
   } `json:"resource"`
}

func GetProfiles(client *http.Client, accessToken, deviceId, deviceTypeId string) (*ProfilesResponse, error) {
   baseURL := "https://abzq7aq4866p.na.api.amazonvideo.com/cdp/switchblade/android/getDataByJvmTransform/v1/dv-android/profiles/listPrimeVideoProfiles/v1.kt"

   params := url.Values{}
   params.Add("ageBracketClassification", "UNKNOWN")
   params.Add("clientName", "ATVAndroidThirdPartyClient")
   params.Add("deviceId", deviceId)
   params.Add("deviceTypeID", deviceTypeId)
   params.Add("firmware", "fmw:30-app:3.0.458.357")
   params.Add("format", "json")
   params.Add("isGeneratedRequest", "false")
   params.Add("osLocale", "en_US")
   params.Add("priorityLevel", "2")
   params.Add("screenDensity", "DEFAULT")
   params.Add("screenWidth", "sw360dp")
   params.Add("softwareVersion", "458")
   params.Add("supportsPKMZ", "false")
   params.Add("swiftPriorityLevel", "critical")
   params.Add("teenProfilesSupported", "true")
   params.Add("uxLocale", "en_US")
   params.Add("version", "1")

   reqURL := fmt.Sprintf("%s?%s", baseURL, params.Encode())

   req, err := http.NewRequest("GET", reqURL, nil)
   if err != nil {
      return nil, fmt.Errorf("failed to create request: %w", err)
   }

   req.Header.Set("Accept", "application/json")
   req.Header.Set("Accept-Language", "en_US")
   req.Header.Set("User-Agent", "Dalvik/2.1.0 (Linux; U; Android 11; sdk_gphone_x86_64 Build/RSR1.240422.006)")
   req.Header.Set("x-atv-page-type", "ATVProfiles")
   req.Header.Set("x-gasc-enabled", "true")
   req.Header.Set("x-request-priority", "CRITICAL")
   req.Header.Set("Authorization", "Bearer "+accessToken)

   resp, err := client.Do(req)
   if err != nil {
      return nil, fmt.Errorf("request failed: %w", err)
   }
   defer resp.Body.Close()

   bodyBytes, err := io.ReadAll(resp.Body)
   if err != nil {
      return nil, fmt.Errorf("failed to read response: %w", err)
   }

   var result ProfilesResponse
   if err := json.Unmarshal(bodyBytes, &result); err != nil {
      return nil, fmt.Errorf("failed to decode response: %w", err)
   }

   return &result, nil
}
