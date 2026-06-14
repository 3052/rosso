package amazon

import (
   "encoding/json"
   "fmt"
   "net/http"
)

// GetPrimaryProfile uses the account access token to fetch available profiles and returns the primary actorId.
func GetPrimaryProfile(accountAccessToken string) (string, error) {
   url := "https://ab8mt4dd97et.na.api.amazonvideo.com/lrcedge/getDataByJavaTransform/v1/lr/profiles/profileSelection"

   req, err := http.NewRequest("GET", url, nil)
   if err != nil {
      return "", err
   }

   q := req.URL.Query()
   q.Add("deviceTypeID", "A2SNKIF736WF4T")
   q.Add("deviceID", "uuidcbb2f9705f13437e9e515622dce02106")
   q.Add("firmware", "google/sdk_gphone_x86/generic_x86_arm:11/RSR1.240422.006/12134477:userdebug/dev-keys")
   q.Add("manufacturer", "Google")
   q.Add("chipset", "goldfish_x86")
   q.Add("model", "sdk_gphone_x86")
   q.Add("operatingSystem", "Android")
   q.Add("clientId", "pv-lrc-rust")
   req.URL.RawQuery = q.Encode()

   req.Header.Set("User-Agent", "Android/google/sdk_gphone_x86/generic_x86_arm:11/RSR1.240422.006/12134477:userdebug/dev-keys, Ignition X/15.5.2026042820-android, Google")
   req.Header.Set("Authorization", "Bearer "+accountAccessToken)
   req.Header.Set("Accept", "application/json")
   req.Header.Set("x-client-app", "avlrc")

   client := &http.Client{}
   resp, err := client.Do(req)
   if err != nil {
      return "", err
   }
   defer resp.Body.Close()

   if resp.StatusCode != http.StatusOK {
      return "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
   }

   var result struct {
      Resource struct {
         Profiles []struct {
            ProfileId        string `json:"profileId"`
            IsDefaultProfile bool   `json:"isDefaultProfile"`
         } `json:"profiles"`
      } `json:"resource"`
   }

   if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
      return "", err
   }

   for _, profile := range result.Resource.Profiles {
      if profile.IsDefaultProfile {
         return profile.ProfileId, nil
      }
   }

   if len(result.Resource.Profiles) > 0 {
      return result.Resource.Profiles[0].ProfileId, nil
   }

   return "", fmt.Errorf("no profiles found")
}
