package amazon

import (
   "encoding/json"
   "fmt"
   "net/http"
   "net/url"
)

type ProfilesResponse struct {
   Resource struct {
      Profiles []struct {
         ProfileId string `json:"profileId"`
         Name      string `json:"name"`
      } `json:"profiles"`
   } `json:"resource"`
}

// GetPrimeVideoProfileId uses the access token to fetch the account's viewing profiles.
// It returns the profileId (actor ID) of the primary profile.
func GetPrimeVideoProfileId(accessToken, deviceId string) (string, error) {
   baseURL, err := url.Parse("https://abzq7aq4866p.na.api.amazonvideo.com/cdp/switchblade/android/getDataByJvmTransform/v1/dv-android/profiles/listPrimeVideoProfiles/v1.kt")
   if err != nil {
      return "", err
   }

   q := baseURL.Query()
   q.Set("ageBracketClassification", "UNKNOWN")
   q.Set("clientName", "ATVAndroidThirdPartyClient")
   q.Set("deviceId", deviceId)
   q.Set("deviceTypeID", "A43PXU4ZN2AL1")
   q.Set("firmware", "fmw:30-app:3.0.458.357")
   q.Set("format", "json")
   q.Set("osLocale", "en_US")
   q.Set("priorityLevel", "2")
   q.Set("softwareVersion", "458")
   q.Set("supportsPKMZ", "false")
   q.Set("swiftPriorityLevel", "critical")
   q.Set("teenProfilesSupported", "true")
   q.Set("uxLocale", "en_US")
   q.Set("version", "1")
   baseURL.RawQuery = q.Encode()

   req, err := http.NewRequest(http.MethodGet, baseURL.String(), nil)
   if err != nil {
      return "", err
   }

   req.Header.Set("Authorization", "Bearer "+accessToken)
   req.Header.Set("User-Agent", "Dalvik/2.1.0 (Linux; U; Android 11; sdk_gphone_x86_64 Build/RSR1.240422.006)")
   req.Header.Set("Accept", "application/json")
   req.Header.Set("x-gasc-enabled", "true")

   client := &http.Client{}
   resp, err := client.Do(req)
   if err != nil {
      return "", err
   }
   defer resp.Body.Close()

   if resp.StatusCode != http.StatusOK {
      return "", fmt.Errorf("expected status 200 OK, got: %d", resp.StatusCode)
   }

   var profResp ProfilesResponse
   if err := json.NewDecoder(resp.Body).Decode(&profResp); err != nil {
      return "", err
   }

   if len(profResp.Resource.Profiles) == 0 {
      return "", fmt.Errorf("no profiles found in the response")
   }

   profileID := profResp.Resource.Profiles[0].ProfileId
   if profileID == "" {
      return "", fmt.Errorf("profile ID was empty in the response")
   }

   return profileID, nil
}
