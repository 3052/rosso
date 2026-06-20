package amazon

import (
   "encoding/json"
   "fmt"
   "net/http"
)

func (*Profile) CachePath() string {
   return "rosso/amazon/Profile"
}

// Profile represents an Amazon actor profile.
type Profile struct {
   ProfileID        string `json:"profileId"`
   IsDefaultProfile bool   `json:"isDefaultProfile"`
}

// GetPrimaryProfile uses the account access token to fetch available profiles and returns the primary profile.
func GetPrimaryProfile(accountAccessToken string) (*Profile, error) {
   url := "https://ab8mt4dd97et.na.api.amazonvideo.com/lrcedge/getDataByJavaTransform/v1/lr/profiles/profileSelection"

   req, err := http.NewRequest("GET", url, nil)
   if err != nil {
      return nil, err
   }

   q := req.URL.Query()
   q.Add("deviceTypeID", DeviceTypeID)
   q.Add("deviceID", DeviceID)
   q.Add("firmware", DeviceFirmware)
   q.Add("manufacturer", DeviceManufacturer)
   q.Add("chipset", DeviceChipset)
   q.Add("model", DeviceModel)
   q.Add("operatingSystem", DeviceOS)
   q.Add("clientId", "pv-lrc-rust")
   req.URL.RawQuery = q.Encode()

   req.Header.Set("User-Agent", UserAgent)
   req.Header.Set("Authorization", "Bearer "+accountAccessToken)
   req.Header.Set("Accept", "application/json")
   req.Header.Set("x-client-app", "avlrc")

   client := &http.Client{}
   resp, err := client.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()

   if resp.StatusCode != http.StatusOK {
      return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
   }

   // Embed our new Profile struct into the anonymous decoder struct
   var result struct {
      Resource struct {
         Profiles []Profile `json:"profiles"`
      } `json:"resource"`
   }

   if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
      return nil, err
   }

   for _, profile := range result.Resource.Profiles {
      if profile.IsDefaultProfile {
         return &profile, nil
      }
   }

   if len(result.Resource.Profiles) > 0 {
      return &result.Resource.Profiles[0], nil
   }

   return nil, fmt.Errorf("no profiles found")
}
