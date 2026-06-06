package amazon

import (
   "encoding/json"
   "fmt"
   "net/http"
)

type ProfilesResponse struct {
   Resource struct {
      Profiles []struct {
         ProfileId string `json:"profileId"`
      } `json:"profiles"`
   } `json:"resource"`
}

func GetPrimeVideoProfileId(accessToken, deviceId string) (string, error) {
   url := fmt.Sprintf("https://abzq7aq4866p.na.api.amazonvideo.com/cdp/switchblade/android/getDataByJvmTransform/v1/dv-android/profiles/listPrimeVideoProfiles/v1.kt?deviceId=%s&deviceTypeID=A43PXU4ZN2AL1&format=json&version=1", deviceId)

   req, err := http.NewRequest("GET", url, nil)
   if err != nil {
      return "", err
   }

   req.Header.Set("Authorization", "Bearer "+accessToken)
   req.Header.Set("User-Agent", "Dalvik/2.1.0 (Linux; U; Android 11; sdk_gphone_x86_64 Build/RSR1.240422.006)")

   client := &http.Client{}
   resp, err := client.Do(req)
   if err != nil {
      return "", err
   }
   defer resp.Body.Close()

   if resp.StatusCode != http.StatusOK {
      return "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
   }

   var profResp ProfilesResponse
   if err := json.NewDecoder(resp.Body).Decode(&profResp); err != nil {
      return "", err
   }

   if len(profResp.Resource.Profiles) == 0 {
      return "", fmt.Errorf("no profiles found")
   }

   return profResp.Resource.Profiles[0].ProfileId, nil
}
